package server

import (
	"context"
	"flag"
	"log"

	pb "auth/auth"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.AuthRequest) (*pb.AuthReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.AuthReply{Message: "Hello " + in.GetName()}, nil
}