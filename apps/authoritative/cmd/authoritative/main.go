package main

// Entry point for author server
import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/jaximus808/delivery-gdg-platform/main/apps/authoritative/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedOrderHandlerServer
}

func (s *server) PlaceOrder(ctx context.Context, order *pb.Order) (*pb.PlaceOrderResponse, error) {
	fmt.Print("Received an order! \n")

	return &pb.PlaceOrderResponse{
		OrderId: order.OrderId,
		Message: fmt.Sprintf("Received order %d from %s! They want %s from %s.",
			order.OrderId, order.CustomerName, order.Item, order.VendorName),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderHandlerServer(s, &server{})

	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
