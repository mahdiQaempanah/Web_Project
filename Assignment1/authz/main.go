package main

import (
	"context"
	"log"
	"net"

	"github.com/mahdiQaempanah/Web_Project/Assignment1/authz/server/pb"

	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedAuthzServer
	P           int32
	G           int32
	b           int32
	NonceLength int
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// We should implement this
func (s *server) PGAgreement(ctx context.Context, req *pb.PGRequest) (*pb.PGResponse, error) {
	nonce := req.Nonce
	messageId := req.MessageId

	return &pb.PGResponse{
		Nonce:       nonce,
		ServerNonce: RandStringRunes(s.NonceLength),
		MessageId:   messageId + 1,
		P:           s.P,
		G:           s.G,
	}, nil
}

func ModularPower(a int32, b int32, p int32) int32 {
	if b == 0 {
		return 1
	}
	var result int32 = ModularPower(a, b/2, p)
	result = int32((((int64)(result)) * ((int64)(result))) % (int64(p)))
	if b%2 == 1 {
		result = int32(((int64(result)) * (int64(a))) % (int64(p)))
	}
	return result
}
func (s *server) DiffieHellman(ctx context.Context, req *pb.DiffieHellmanRequest) (*pb.DiffieHellmanResponse, error) {
	nonce := req.Nonce
	serverNonce := req.ServerNonce
	messageId := req.MessageId
	// GA := req.GA
	GB := ModularPower(s.G, s.b, s.P)

	return &pb.DiffieHellmanResponse{
		Nonce:       nonce,
		ServerNonce: serverNonce,
		MessageId:   messageId + 1,
		GB:          GB,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	rand.Seed(382992)

	s := grpc.NewServer()

	pb.RegisterAuthzServer(s, &server{P: 23, G: 5, b: 15, NonceLength: 20})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
