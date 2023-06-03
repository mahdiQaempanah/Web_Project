package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"os"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
	"github.com/mahdiQaempanah/Web_Project/Assignment1/biz/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedBizServer
	sql_db *sql.DB
	redis  *redis.Client
}

func (s *server) ValidateUser(auth_key int32) (bool, error) {
	val, err := s.redis.Get(auth_key).Result()
	if err != nil {
		return false, err
	}
	if val == 1 {
		return true, nil
	}
	return false, nil
}

// We should implement this
func (s *server) GetUsers(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userId := req.UserID
	auth_key := req.AuthKey
	messageId := req.MessageId

	validate_result, err := s.ValidateUser(auth_key)
	if err != nil {
		return nil, err
	}
	if validate_result == false {
		return nil, errors.New("Authentication Failed.")
	}

	var is_empty bool
	if err = s.sql_db.QueryRow("select * from USERS where id = ? IS EMPTY", userId).Scan(&is_empty); err != nil {
		return nil, err
	}
	var rows *sql.Rows

	if is_empty == true {
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

	result := []*pb.User{}
	for rows.Next() {
		var user pb.User
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		result = append(result, &user)
	}

	return &pb.GetUserResponse{
		Users:     result,
		MessageId: messageId + 1}, nil
}

func (s *server) GetUsersWithSQLInject(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// nonce := req.UserID
	// auth_key = req.AuthKey
	// messageId := req.MessageId
	return nil, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
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
	pb.RegisterBizServer(s, &server{sql_db: db, redis: rdb})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
