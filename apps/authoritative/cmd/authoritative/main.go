package main

// Entry point for author server
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jaximus808/delivery-gdg-platform/main/apps/authoritative/internal/matcher"
	"strconv"

	pb "github.com/jaximus808/delivery-gdg-platform/main/apps/authoritative/proto"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedOrderHandlerServer
	sb *supabase.Client
}

func (s *server) InsertOrder(ctx context.Context, req *pb.InsertOrderRequest) (*pb.InsertOrderResponse, error) {
	order := req.GetOrder()

	//ok time to insert order into database
	orderData := map[string]interface{}{
		"userId":          order.GetUserId(),
		"vendorId":        order.GetVendorId(),
		"status":          order.GetStatus(),
		"dropOffLocation": order.GetDropoffLocId(),
	}

	inserted, _, err := s.sb.
		From("orders").
		Insert(orderData, false, "representation", "", "").
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed inserting order: %v", err)
	}

	var rows []struct {
		OrderID int64 `json:"id"`
	}

	json.Unmarshal(inserted, &rows)
	orderId := rows[0].OrderID

	order.OrderId = orderId

	//ok time to insert items in order into orderitems table
	for _, item := range order.GetItems() {
		itemData := map[string]interface{}{
			"orderId":  orderId,
			"itemName": item.GetItemName(),
			"quantity": item.GetQuantity(),
			"price":    item.GetPrice(),
		}

		_, _, err := s.sb.
			From("orderItems").
			Insert(itemData, false, "", "", "").
			Execute()

		if err != nil {
			return nil, fmt.Errorf("failed inserting order item: %v", err)
		}
	}

	return &pb.InsertOrderResponse{
		Order:     order,
		ReturnMsg: "SUCCESS",
	}, nil
}

func (s *server) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	order := req.GetOrder()
	orderId := order.GetOrderId()

	// Delete order items first due to foreign key constraints
	_, _, err := s.sb.
		From("orderItems").
		Delete("", "").
		Eq("orderId", strconv.Itoa(int(orderId))).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed deleting order items: %v", err)
	}

	// Delete the order
	_, _, err = s.sb.
		From("orders").
		Delete("", "").
		Eq("id", strconv.Itoa(int(orderId))).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed deleting order: %v", err)
	}

	return &pb.DeleteOrderResponse{
		ReturnMsg: "SUCCESS",
	}, nil
}

/*
func (s *server) InsertItem(ctx context.Context, req *pb.InsertItemRequest) (*pb.InsertItemResponse, error) {
	fmt.Print("Received a food item to insert! \n")

	item := req.GetItem()

	return &pb.InsertItemResponse{
		Item: item,
		ReturnMsg: fmt.Sprintf("Received item %s from %d.",
			item.ItemName, item.VendorId),
	}, nil
}
*/

func main() {
	godotenv.Load("../../.env")

	SUPABASE_URL := os.Getenv("SUPABASE_URL")
	SUPABASE_KEY := os.Getenv("SUPABASE_KEY")

	client, err := supabase.NewClient(
		SUPABASE_URL,
		SUPABASE_KEY,
		nil,
	)

	if err != nil {
		log.Fatalf("failed to create supabase client: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	orm := matcher.CreateOrderRobotMatcher()
	grpc_server := grpc.NewServer()
	pb.RegisterOrderHandlerServer(grpc_server, &server{sb: client})

	go orm.StartORM()

	s := grpc.NewServer()
	// proto.RegisterHelloServiceServer(s, &server{})
	log.Println("gRPC server listening on :50051")
	if err := grpc_server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
