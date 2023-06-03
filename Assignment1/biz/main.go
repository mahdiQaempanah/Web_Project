package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/mahdiQaempanah/Web_Project/Assignment1/biz/grpc/biz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	biz.UnimplementedBizServer
	sql_db *sql.DB
	redis  *redis.Client
}

func (s *server) ValidateUser(auth_key int32) (bool, error) {
	val, err := s.redis.Get(strconv.Itoa(int(auth_key))).Result()
	if err != nil {
		return false, err
	}
	valint, err2 := strconv.Atoi(val)
	if err2 != nil {
		return false, err2
	}
	if valint == 1 {
		return true, nil
	}
	return false, nil
}

func checkUserValidity(user string) bool {
	res, err := regexp.MatchString("^[0-9]*$", user)
	if err != nil {
		return false
	}
	return res
}

func (s *server) GetUsers(ctx context.Context, req *biz.GetUserRequest) (*biz.GetUserResponse, error) {
	userId := req.UserID
	auth_key := req.AuthKey
	messageId := req.MessageId

	validate_result, err := s.ValidateUser(auth_key)
	if err != nil {
		return nil, err
	}
	if !validate_result {
		return nil, errors.New("authentication failed")
	}

	if !checkUserValidity(userId) {
		return nil, errors.New("userId is not numeric")
	}

	var is_empty bool
	if err = s.sql_db.QueryRow("select * from USERS where id = ? IS EMPTY", userId).Scan(&is_empty); err != nil {
		return nil, err
	}
	var rows *sql.Rows

	if is_empty {
		var err error
		rows, err = s.sql_db.Query("select top 100 * from USERS")
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		rows, err = s.sql_db.Query("select * from USERS where Id = ?", userId)
		if err != nil {
			return nil, err
		}
	}

	result := []*biz.User{}
	for rows.Next() {
		var user biz.User
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		result = append(result, &user)
	}

	return &biz.GetUserResponse{
		Users:     result,
		MessageId: messageId + 1}, nil
}

func (s *server) GetUsersWithSQLInject(ctx context.Context, req *biz.GetUserRequest) (*biz.GetUserResponse, error) {
	userId := req.UserID
	auth_key := req.AuthKey
	messageId := req.MessageId

	validate_result, err := s.ValidateUser(auth_key)
	if err != nil {
		return nil, err
	}
	if !validate_result {
		return nil, errors.New("authentication failed")
	}

	var is_empty bool
	if err = s.sql_db.QueryRow("select * from USERS where id = ? IS EMPTY", userId).Scan(&is_empty); err != nil {
		return nil, err
	}
	var rows *sql.Rows

	if is_empty {
		var err error
		rows, err = s.sql_db.Query("select top 100 * from USERS")
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		rows, err = s.sql_db.Query("select * from USERS where Id = ?", userId)
		if err != nil {
			return nil, err
		}
	}

	result := []*biz.User{}
	for rows.Next() {
		var user biz.User
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		result = append(result, &user)
	}

	return &biz.GetUserResponse{
		Users:     result,
		MessageId: messageId + 1}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":5062")
	if err != nil {
		panic(err)
	}

	cfg := mysql.Config{
		User:   os.Getenv("POSTGRES_USER"),
		Passwd: os.Getenv("POSTGRES_PASSWORD"),
		Net:    "tcp",
		Addr:   "localhost:5432",
		DBName: "POSTGRES_DB",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	s := grpc.NewServer()
	biz.RegisterBizServer(s, &server{sql_db: db, redis: rdb})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
