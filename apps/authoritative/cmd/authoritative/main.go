package main

// Entry point for author server
import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jaximus808/delivery-gdg-platform/main/apps/authoritative/example/proto"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func add(x int, y int) int {
	return x + y
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterHelloServiceServer(s, &server{})

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
