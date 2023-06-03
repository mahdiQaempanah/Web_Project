package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/mahdiQaempanah/Web_Project/Assignment1/authz/server/pb"

	"math/rand"

	"crypto/sha1"
	"encoding/hex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedAuthzServer
	NonceLength int
	redis       *redis.Client
	P           int32
	G           int32
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CalculateSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *server) PGAgreement(ctx context.Context, req *pb.PGRequest) (*pb.PGResponse, error) {
	nonce := req.Nonce
	messageId := req.MessageId
	ServerNonce := RandStringRunes(s.NonceLength)

	hash := CalculateSha1(nonce + ServerNonce)
	err := s.redis.Set(hash, rand.Int31(), 20*time.Minute).Err()
	if err != nil {
		return nil, err
	}
	return &pb.PGResponse{
		Nonce:       nonce,
		ServerNonce: ServerNonce,
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
	GA := req.GA
	hash := CalculateSha1(nonce + serverNonce)
	b, err := s.redis.Get(hash).Result()
	if err != nil {
		return nil, err
	}
	intb, _ := strconv.Atoi(b)
	GB := ModularPower(s.G, int32(intb), s.P)

	GAB := ModularPower(GA, int32(intb), s.P)
	err2 := s.redis.Set(string(GAB), 1, 20*time.Minute).Err()
	if err2 != nil {
		return nil, err
	}
	return &pb.DiffieHellmanResponse{
		Nonce:       nonce,
		ServerNonce: serverNonce,
		MessageId:   messageId + 1,
		GB:          GB,
	}, nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	rand.Seed(382992)

	s := grpc.NewServer()

	pb.RegisterAuthzServer(s, &server{P: 23, G: 5, redis: rdb, NonceLength: 20})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
