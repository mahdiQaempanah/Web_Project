package main

import (
	"context"
	"log"
	"net"

	"github.com/mahdiQaempanah/Web_Project/Assignment1/authz/server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedAuthzServer
}

// We should implement this
func (s *server) PGAgreement(ctx context.Context, req *pb.PGRequest) (*pb.PGResponse, error) {
	nonce := req.Nonce
	messageId := req.MessageId

	return &pb.PGResponse{
		Nonce:       nonce,
		ServerNonce: nonce,
		MessageId:   messageId + 1,
		P:           23,
		G:           5,
	}, nil
}

func (s *server) DiffieHellman(ctx context.Context, req *pb.DiffieHellmanRequest) (*pb.DiffieHellmanResponse, error) {
	nonce := req.Nonce
	serverNonce := req.ServerNonce
	messageId := req.MessageId
	GA := req.GA

	return &pb.DiffieHellmanResponse{
		Nonce:       nonce,
		ServerNonce: serverNonce,
		MessageId:   messageId + 1,
		GB:          GA,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthzServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
