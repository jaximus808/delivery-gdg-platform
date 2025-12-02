package matcher

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateOrderRobotMatcher(t *testing.T) {
	orm := CreateOrderRobotMatcher()

	if orm == nil {
		t.Fatal("CreateOrderRobotMatcher returned nil")
	}

	if orm.orderIntake == nil {
		t.Error("orderIntake channel not initialized")
	}

	if orm.robotIntake == nil {
		t.Error("robotIntake channel not initialized")
	}

	if orm.orderQueue == nil {
		t.Error("orderQueue not initialized")
	}

	if orm.robotQueue == nil {
		t.Error("robotQueue not initialized")
	}

	if orm.orderCount != 0 {
		t.Errorf("expected orderCount to be 0, got %d", orm.orderCount)
	}
}

func TestSubmitOrder(t *testing.T) {
	orm := CreateOrderRobotMatcher()

	order := &OrderItem{
		orderId: 1,
	}

	// Submit order in a goroutine to avoid blocking
	done := make(chan bool)
	go func() {
		orm.SubmitOrder(order)
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(time.Second):
		t.Fatal("SubmitOrder blocked unexpectedly")
	}

	// Verify order is in channel
	select {
	case receivedOrder := <-orm.orderIntake:
		if receivedOrder.orderId != order.orderId {
			t.Errorf("expected orderId %d, got %d", order.orderId, receivedOrder.orderId)
		}
	case <-time.After(time.Second):
		t.Fatal("Order not received in channel")
	}
}

func TestSubmitRobot(t *testing.T) {
	orm := CreateOrderRobotMatcher()

	robot := &RobotUpdate{
		robotID: "robot-1",
		status:  "online",
	}

	// Submit robot in a goroutine to avoid blocking
	done := make(chan bool)
	go func() {
		orm.SubmitRobot(robot)
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(time.Second):
		t.Fatal("SubmitRobot blocked unexpectedly")
	}

	// Verify robot is in channel
	select {
	case receivedRobot := <-orm.robotIntake:
		if receivedRobot.robotID != robot.robotID {
			t.Errorf("expected robotID %s, got %s", robot.robotID, receivedRobot.robotID)
		}
	case <-time.After(time.Second):
		t.Fatal("Robot not received in channel")
	}
}

func TestAttemptMatchNoOrdersOrRobots(t *testing.T) {
	orm := CreateOrderRobotMatcher()
	matchesChan := make(chan *OrderRobotMatch, 10)

	orm.attemptMatch(matchesChan)

	// Should not produce a match
	select {
	case <-matchesChan:
		t.Error("Expected no match when queues are empty")
	case <-time.After(100 * time.Millisecond):
		// Expected behavior
	}
}

func TestAttemptMatchWithOrderAndRobot(t *testing.T) {
	orm := CreateOrderRobotMatcher()
	matchesChan := make(chan *OrderRobotMatch, 10)

	// Add order to queue
	order := &OrderItem{
		orderId: 123,
	}
	orm.orderQueue.Insert(order)

	// Add robot to queue
	robot := RobotItem{
		robotID: "robot-456",
	}
	orm.robotQueue.Enqueue(robot)

	orm.attemptMatch(matchesChan)

	// Should produce a match
	select {
	case match := <-matchesChan:
		if match.OrderID != 123 {
			t.Errorf("expected OrderID 123, got %d", match.OrderID)
		}
		if match.RobotID != "robot-456" {
			t.Errorf("expected RobotID robot-456, got %s", match.RobotID)
		}
	case <-time.After(time.Second):
		t.Fatal("Expected a match but none was produced")
	}
}

func TestStartORMReturnsChannel(t *testing.T) {
	orm := CreateOrderRobotMatcher()

	matchesChan := orm.StartORM()

	if matchesChan == nil {
		t.Fatal("StartORM returned nil channel")
	}
}

func TestEngineProcessesOrders(t *testing.T) {
	orm := CreateOrderRobotMatcher()
	matchesChan := orm.StartORM()

	// Submit an order
	order := &OrderItem{
		orderId: 999,
	}
	orm.SubmitOrder(order)

	// Give engine time to process
	time.Sleep(100 * time.Millisecond)

	// Verify order was added to queue and orderCount incremented
	if orm.orderCount != 1 {
		t.Errorf("expected orderCount 1, got %d", orm.orderCount)
	}

	if orm.orderQueue.Len() != 1 {
		t.Errorf("expected orderQueue length 1, got %d", orm.orderQueue.Len())
	}

	// Clean up
	close(matchesChan)
}

func TestEngineProcessesRobotOnline(t *testing.T) {
	orm := CreateOrderRobotMatcher()
	matchesChan := orm.StartORM()

	// Submit a robot with online status
	robot := &RobotUpdate{
		robotID: "robot-online",
		status:  "online",
	}
	orm.SubmitRobot(robot)

	// Give engine time to process
	time.Sleep(100 * time.Millisecond)

	// Verify robot was added to queue
	if orm.robotQueue.Len() != 1 {
		t.Errorf("expected robotQueue length 1, got %d", orm.robotQueue.Len())
	}

	// Clean up
	close(matchesChan)
}

func TestEngineProcessesRobotOffline(t *testing.T) {
	orm := CreateOrderRobotMatcher()
	matchesChan := orm.StartORM()

	// First add a robot
	robot := &RobotUpdate{
		robotID: "robot-test",
		status:  "online",
	}
	orm.SubmitRobot(robot)
	time.Sleep(100 * time.Millisecond)

	// Now send offline status
	robotOffline := &RobotUpdate{
		robotID: "robot-test",
		status:  "offline",
	}
	orm.SubmitRobot(robotOffline)
	time.Sleep(100 * time.Millisecond)

	// Verify robot was removed from queue
	if orm.robotQueue.Len() != 0 {
		t.Errorf("expected robotQueue length 0, got %d", orm.robotQueue.Len())
	}

	// Clean up
	close(matchesChan)
}

func TestEngineCreatesMatchesOnTicker(t *testing.T) {
	orm := CreateOrderRobotMatcher()
	matchesChan := orm.StartORM()

	// Submit order and robot
	order := &OrderItem{
		orderId: 111,
	}
	orm.SubmitOrder(order)

	robot := &RobotUpdate{
		robotID: "robot-222",
		status:  "online",
	}
	orm.SubmitRobot(robot)

	// Wait for ticker to fire (slightly over 1 second)
	select {
	case match := <-matchesChan:
		if match.OrderID != 111 {
			t.Errorf("expected OrderID 111, got %d", match.OrderID)
		}
		if match.RobotID != "robot-222" {
			t.Errorf("expected RobotID robot-222, got %s", match.RobotID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Expected match to be created within 2 seconds")
	}

	// Clean up
	close(matchesChan)
}

func TestEngineMultipleMatches(t *testing.T) {
	orm := CreateOrderRobotMatcher()
	matchesChan := orm.StartORM()

	// Submit multiple orders and robots
	for i := 1; i <= 3; i++ {
		order := &OrderItem{
			orderId: i * 100,
		}
		orm.SubmitOrder(order)

		robot := &RobotUpdate{
			robotID: fmt.Sprintf("robot-%d", i),
			status:  "online",
		}
		orm.SubmitRobot(robot)
	}

	// Collect matches (should get 3 matches over ~3 seconds)
	matchCount := 0
	timeout := time.After(5 * time.Second)

	for matchCount < 3 {
		select {
		case match := <-matchesChan:
			if match == nil {
				t.Error("Received nil match")
			}
			matchCount++
		case <-timeout:
			t.Fatalf("Expected 3 matches, got %d", matchCount)
		}
	}

	if matchCount != 3 {
		t.Errorf("expected 3 matches, got %d", matchCount)
	}

	// Clean up
	close(matchesChan)
}

func TestEngineOrderCountIncrement(t *testing.T) {
	orm := CreateOrderRobotMatcher()
	matchesChan := orm.StartORM()

	// Submit multiple orders
	for i := 0; i < 5; i++ {
		order := &OrderItem{
			orderId: i,
		}
		orm.SubmitOrder(order)
	}

	// Give engine time to process all orders
	time.Sleep(200 * time.Millisecond)

	if orm.orderCount != 5 {
		t.Errorf("expected orderCount 5, got %d", orm.orderCount)
	}

	// Clean up
	close(matchesChan)
}

func TestEngineChannelBuffering(t *testing.T) {
	orm := CreateOrderRobotMatcher()

	// Test order intake buffer (100)
	for i := 0; i < 100; i++ {
		order := &OrderItem{
			orderId: i,
		}
		// Should not block
		orm.SubmitOrder(order)
	}

	// Test robot intake buffer (100)
	for i := 0; i < 100; i++ {
		robot := &RobotUpdate{
			robotID: fmt.Sprintf("robot-%d", i),
			status:  "online",
		}
		// Should not block
		orm.SubmitRobot(robot)
	}

	// If we got here without blocking, test passes
}
