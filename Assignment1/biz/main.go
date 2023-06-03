package main

import (
	"context"
	"log"
	"net"

	"github.com/mahdiQaempanah/Web_Project/Assignment1/biz/server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedBizServer
}

// We should implement this
func (s *server) GetUsers(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// nonce := req.UserID
	// auth_key = req.AuthKey
	// messageId := req.MessageId
	return nil, nil
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

	s := grpc.NewServer()

	pb.RegisterBizServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
