package main

import (
	"context"
	"log"
	"net"

	"main/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedReqPqServiceServer
}

func (s *server) RequestPq(ctx context.Context, in *pb.RequestPqRequest) (*pb.RequestPqResponse, error) {
	nonce := in.Nonce
	messageId := in.MessageId
	return &pb.RequestPqResponse{
		Nonce:       nonce,
		ServerNonce: nonce,
		MessageId:   messageId + 1,
		P:           23,
		G:           5,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterReqPqServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
