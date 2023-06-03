package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/mahdiQaempanah/Web_Project/Assignment1/biz/grpc/biz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	biz.UnimplementedBizServer
	db     *sql.DB
	redis  *redis.Client
	logger *log.Logger
}

func (s *server) isAuthenticated(auth_key int32) error {
	val, err := s.redis.Get(strconv.Itoa(int(auth_key))).Result()
	if err != nil {
		return err
	}
	valint, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	if valint == 1 {
		return nil
	}
	return errors.New("authentication failed")
}

func checkUserValidity(user string) bool {
	res, err := regexp.MatchString("^[0-9]*$", user)
	if err != nil {
		return false
	}
	return res
}

func (s *server) postgresSelectUser(userId string) ([]*biz.User, error) {
	s.logger.Println("User is valid")

	var is_empty int
	statment := fmt.Sprintf(`select count(*) from USERS where id='%s';`, userId)
	s.logger.Println("statement is : " + statment)
	if err := s.db.QueryRow(statment).Scan(&is_empty); err != nil {
		return nil, err
	}
	var rows *sql.Rows

	s.logger.Println("Is empty: " + fmt.Sprint(is_empty))

	if is_empty == 0 {
		var err error
		query := "select * from USERS limit 100;"
		rows, err = s.db.Query(query)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		query := fmt.Sprintf(`select * from USERS where id='%s';`, userId)
		rows, err = s.db.Query(query)
		if err != nil {
			return nil, err
		}
	}

	result := []*biz.User{}
	for rows.Next() {
		var user biz.User
		if err := rows.Scan(&user.Name, &user.Surname, &user.Id, &user.Age, &user.Sex); err != nil {
			return nil, err
		}
		result = append(result, &user)
	}

	return result, nil
}

func (s *server) GetUsers(ctx context.Context, req *biz.GetUserRequest) (*biz.GetUserResponse, error) {
	userId := req.UserID
	auth_key := req.AuthKey
	messageId := req.MessageId

	err := s.isAuthenticated(auth_key)
	if err != nil {
		return nil, err
	}

	if !checkUserValidity(userId) {
		return nil, errors.New("user_id is not numeric")
	}

	users, err := s.postgresSelectUser(userId)
	if err != nil {
		return nil, err
	}
	return &biz.GetUserResponse{
		Users:     users,
		MessageId: messageId + 1}, nil
}

func (s *server) GetUsersWithSQLInject(ctx context.Context, req *biz.GetUserRequest) (*biz.GetUserResponse, error) {
	userId := req.UserID
	auth_key := req.AuthKey
	messageId := req.MessageId

	err := s.isAuthenticated(auth_key)
	if err != nil {
		return nil, err
	}

	users, err := s.postgresSelectUser(userId)
	if err != nil {
		return nil, err
	}
	return &biz.GetUserResponse{
		Users:     users,
		MessageId: messageId + 1}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":5062")
	if err != nil {
		panic(err)
	}

	// TODO: use viper
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error: " + err.Error())
	} else {
		fmt.Println("Successfully connected to Fucking DB")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	s := grpc.NewServer()
	biz.RegisterBizServer(s, &server{db: db, redis: rdb, logger: log.New(os.Stdout, "logger: ", log.Ldate|log.Ltime)})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
