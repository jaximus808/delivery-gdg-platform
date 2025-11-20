package db

import (
	"context"
	"os"

	"github.com/supabase-community/postgrest-go"

	"github.com/joho/godotenv"
)

type Database struct {
	client *postgrest.Client
}

func New() *Database {
	godotenv.Load()
	url := os.Getenv("SUPABASE_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	client := postgrest.NewClient(url, apiKey, nil)
	return &Database{client: client}
}

// Coordinate Type Enum
// 1 = Vendor
// 2 = Dropoff
// 3 = Waypoint

const (
	CoordinateTypeVendor   = 1
	CoordinateTypeDropoff  = 2
	CoordinateTypeWaypoint = 3
)

type Coordinate struct {
	ID   string      `json:"id"`
	X    int         `json:"x"`
	Y    int         `json:"y"`
	Meta interface{} `json:"meta"`
	Type int16       `json:"type"`
}

type OrderItem struct {
	ID       string  `json:"id"`
	OrderID  int64   `json:"orderId"`
	ItemName string  `json:"itemName"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	ID              int64  `json:"id"`
	UserID          string `json:"userId"`
	VendorID        string `json:"vendorId"`
	Status          int    `json:"status"`
	CreatedAt       string `json:"createdAt"`
	RobotID         string `json:"robotId"`
	DropOffLocation string `json:"dropOffLocation"`
}

type Robot struct {
	ID         string `json:"id"`
	Status     int    `json:"status"`
	LastUpdate string `json:"lastUpdate"`
	CurrentLoc string `json:"currentLoc"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhoneNum string `json:"phoneNum"`
}

type Vendor struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Hours       interface{} `json:"hours"`
	Coordinates string      `json:"coordinates"`
}

func (db *Database) InsertCoordinate(ctx context.Context, c Coordinate) error { return nil }
func (db *Database) GetCoordinate(ctx context.Context, id string) (Coordinate, error) {
	return Coordinate{}, nil
}
func (db *Database) ListCoordinates(ctx context.Context) ([]Coordinate, error) { return nil, nil }
func (db *Database) DeleteCoordinate(ctx context.Context, id string) error     { return nil }

func (db *Database) CreateOrder(ctx context.Context, o Order) error        { return nil }
func (db *Database) GetOrder(ctx context.Context, id int64) (Order, error) { return Order{}, nil }
func (db *Database) ListOrdersByUser(ctx context.Context, userID string) ([]Order, error) {
	return nil, nil
}

func (db *Database) ListOrdersByVendor(ctx context.Context, vendorID string) ([]Order, error) {
	return nil, nil
}
func (db *Database) UpdateOrderStatus(ctx context.Context, id int64, status int) error { return nil }
func (db *Database) AssignOrderToRobot(ctx context.Context, orderID int64, robotID string) error {
	return nil
}
func (db *Database) DeleteOrder(ctx context.Context, id int64) error { return nil }

func (db *Database) CreateOrderWithItems(ctx context.Context, order Order, items []OrderItem) error {
	return nil
}

func (db *Database) AddOrderItem(ctx context.Context, item OrderItem) error { return nil }
func (db *Database) GetOrderItems(ctx context.Context, orderID int64) ([]OrderItem, error) {
	return nil, nil
}
func (db *Database) DeleteOrderItem(ctx context.Context, id string) error { return nil }

func (db *Database) GetRobot(ctx context.Context, id string) (Robot, error)          { return Robot{}, nil }
func (db *Database) SetRobotStatus(ctx context.Context, id string, status int) error { return nil }
func (db *Database) UpdateRobotLocation(ctx context.Context, id string, coordinateID string) error {
	return nil
}
func (db *Database) ListRobots(ctx context.Context) ([]Robot, error)  { return nil, nil }
func (db *Database) DeleteRobot(ctx context.Context, id string) error { return nil }

func (db *Database) InsertUser(ctx context.Context, u User) error         { return nil }
func (db *Database) GetUser(ctx context.Context, id string) (User, error) { return User{}, nil }
func (db *Database) ListUsers(ctx context.Context) ([]User, error)        { return nil, nil }
func (db *Database) DeleteUser(ctx context.Context, id string) error      { return nil }

func (db *Database) InsertVendor(ctx context.Context, v Vendor) error         { return nil }
func (db *Database) GetVendor(ctx context.Context, id string) (Vendor, error) { return Vendor{}, nil }
func (db *Database) ListVendors(ctx context.Context) ([]Vendor, error)        { return nil, nil }
func (db *Database) DeleteVendor(ctx context.Context, id string) error        { return nil }
