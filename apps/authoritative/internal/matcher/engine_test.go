package matcher

import (
	"testing"
	"time"
)

func TestStartEngine(t *testing.T) {
	tests := []struct {
		name            string
		orders          []*OrderItem
		robots          []*RobotUpdate
		expectedMatches int
		waitTime        time.Duration
	}{
		{
			name: "basic matching - 3 orders and 3 robots",
			orders: []*OrderItem{
				CreateOrder(1, 101, 0),
				CreateOrder(2, 102, 0),
				CreateOrder(3, 103, 0),
			},
			robots: []*RobotUpdate{
				{robotID: 201, status: "online"},
				{robotID: 202, status: "online"},
				{robotID: 203, status: "online"},
			},
			expectedMatches: 3,
			waitTime:        4 * time.Second,
		},
		{
			name: "more orders than robots",
			orders: []*OrderItem{
				CreateOrder(1, 101, 0),
				CreateOrder(2, 102, 0),
				CreateOrder(3, 103, 0),
				CreateOrder(4, 104, 0),
				CreateOrder(5, 105, 0),
			},
			robots: []*RobotUpdate{
				{robotID: 201, status: "online"},
				{robotID: 202, status: "online"},
			},
			expectedMatches: 2,
			waitTime:        3 * time.Second,
		},
		{
			name: "more robots than orders",
			orders: []*OrderItem{
				CreateOrder(1, 101, 0),
				CreateOrder(2, 102, 0),
			},
			robots: []*RobotUpdate{
				{robotID: 201, status: "online"},
				{robotID: 202, status: "online"},
				{robotID: 203, status: "online"},
				{robotID: 204, status: "online"},
			},
			expectedMatches: 2,
			waitTime:        3 * time.Second,
		},
		{
			name: "robot goes offline",
			orders: []*OrderItem{
				CreateOrder(1, 101, 0),
				CreateOrder(2, 102, 0),
			},
			robots: []*RobotUpdate{
				{robotID: 201, status: "online"},
				{robotID: 202, status: "online"},
				{robotID: 202, status: "offline"}, // same robot goes offline
			},
			expectedMatches: 1,
			waitTime:        3 * time.Second,
		},
		{
			name:            "no orders or robots",
			orders:          []*OrderItem{},
			robots:          []*RobotUpdate{},
			expectedMatches: 0,
			waitTime:        2 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create matcher
			orm := CreateOrderRobotMatcher()
			matchesChan := make(chan *OrderRobotMatch, 100)

			// Start the engine in a goroutine
			go orm.startEngine(matchesChan)

			// Send orders
			for _, order := range tt.orders {
				orm.orderIntake <- order
			}

			// Send robot updates
			for _, robot := range tt.robots {
				orm.robotIntake <- robot
			}

			// Collect matches
			matches := make([]*OrderRobotMatch, 0)
			done := make(chan bool)

			go func() {
				timeout := time.After(tt.waitTime)
				for {
					select {
					case match := <-matchesChan:
						matches = append(matches, match)
						t.Logf("Match received: OrderID=%d, RobotID=%d", match.orderID, match.robotID)
					case <-timeout:
						done <- true
						return
					}
				}
			}()

			<-done

			// Verify match count
			if len(matches) != tt.expectedMatches {
				t.Errorf("Expected %d matches, got %d", tt.expectedMatches, len(matches))
			}

			// Verify no duplicate order IDs in matches
			orderIDsSeen := make(map[int]bool)
			for _, match := range matches {
				if orderIDsSeen[match.orderID] {
					t.Errorf("Duplicate order ID %d in matches", match.orderID)
				}
				orderIDsSeen[match.orderID] = true
			}

			// Verify no duplicate robot IDs in matches
			robotIDsSeen := make(map[int]bool)
			for _, match := range matches {
				if robotIDsSeen[match.robotID] {
					t.Errorf("Duplicate robot ID %d in matches", match.robotID)
				}
				robotIDsSeen[match.robotID] = true
			}

			// Log final state
			t.Logf("Final state - Orders in queue: %d, Robots in queue: %d",
				orm.orderQueue.Len(), orm.robotQueue.Len())
		})
	}
}

func TestStartEngineMatchingOrder(t *testing.T) {
	// Test that orders are matched in FIFO order (based on orderNum)
	orm := CreateOrderRobotMatcher()
	matchesChan := make(chan *OrderRobotMatch, 100)

	go orm.startEngine(matchesChan)

	// Send orders
	orders := []*OrderItem{
		CreateOrder(1, 101, 0),
		CreateOrder(2, 102, 0),
		CreateOrder(3, 103, 0),
	}

	for _, order := range orders {
		orm.orderIntake <- order
		time.Sleep(50 * time.Millisecond) // Small delay to ensure order
	}

	// Send robots
	robots := []*RobotUpdate{
		{robotID: 201, status: "online"},
		{robotID: 202, status: "online"},
		{robotID: 203, status: "online"},
	}

	for _, robot := range robots {
		orm.robotIntake <- robot
	}

	// Collect matches
	matches := make([]*OrderRobotMatch, 0)
	timeout := time.After(4 * time.Second)

	collecting := true
	for collecting && len(matches) < 3 {
		select {
		case match := <-matchesChan:
			matches = append(matches, match)
		case <-timeout:
			collecting = false
		}
	}

	// Verify matches are in order
	if len(matches) >= 3 {
		if matches[0].orderID != 101 {
			t.Errorf("Expected first match to be order 101, got %d", matches[0].orderID)
		}
		if matches[1].orderID != 102 {
			t.Errorf("Expected second match to be order 102, got %d", matches[1].orderID)
		}
		if matches[2].orderID != 103 {
			t.Errorf("Expected third match to be order 103, got %d", matches[2].orderID)
		}
	}
}
