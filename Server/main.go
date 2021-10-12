package main

import (
	"context"
	"log"
	"net"

	//pb "MiniProject2/src/Chitty_Chat_Server"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedCoursesServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) BroadcastMessage(ctx context.Context, in *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
	log.Printf("Received: %v", in.GetMessage())
	return &pb.CoursesReply{Id: "4567"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCoursesServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}