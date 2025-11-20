package main

// Entry point for author server
import (
	"log"
	"net"

	"github.com/jaximus808/delivery-gdg-platform/main/apps/authoritative/internal/matcher"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	orm := matcher.CreateOrderRobotMatcher()

	go orm.StartORM()

	s := grpc.NewServer()
	// proto.RegisterHelloServiceServer(s, &server{})
	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
