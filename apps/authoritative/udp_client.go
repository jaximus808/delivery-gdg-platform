package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

// UDPClientConnection represents a UDP client connection
type UDPClientConnection struct {
	conn       *net.UDPConn
	serverAddr *net.UDPAddr
	clientType string
	clientID   string
}

// NewUDPClient creates a new UDP client
func NewUDPClient(serverAddr, clientType, clientID string) (*UDPClientConnection, error) {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve server address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %v", err)
	}

	client := &UDPClientConnection{
		conn:       conn,
		serverAddr: addr,
		clientType: clientType,
		clientID:   clientID,
	}

	log.Printf("UDP Client created as %s (Type: %s)", clientID, clientType)

	return client, nil
}

// Listen starts listening for messages from the server
func (c *UDPClientConnection) Listen() {
	buffer := make([]byte, 1024)
	for {
		n, _, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from server: %v", err)
			return
		}
		fmt.Printf("[Server]: %s\n", string(buffer[:n]))
	}
}

// Send sends a message to the server
func (c *UDPClientConnection) Send(message string) error {
	// Format: TYPE:ID:MESSAGE
	packet := fmt.Sprintf("%s:%s:%s", c.clientType, c.clientID, message)
	_, err := c.conn.Write([]byte(packet))
	if err != nil {
		return fmt.Errorf("failed to send packet: %v", err)
	}
	return nil
}

// Close closes the connection
func (c *UDPClientConnection) Close() {
	c.conn.Close()
}

// SimulateRobotUDP simulates a robot client sending periodic updates
func SimulateRobotUDP(serverAddr, robotID string, duration time.Duration) {
	client, err := NewUDPClient(serverAddr, "robot", robotID)
	if err != nil {
		log.Fatalf("Failed to create robot client: %v", err)
	}
	defer client.Close()

	// Start listening in background
	go client.Listen()

	// Simulate robot sending status updates
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(duration)
	counter := 0

	for {
		select {
		case <-ticker.C:
			counter++
			status := fmt.Sprintf("Position[%d,%d]_Battery:%d%%",
				counter*10, counter*5, 100-counter*2)
			if err := client.Send(status); err != nil {
				log.Printf("Error sending packet: %v", err)
				return
			}
		case <-timeout:
			log.Printf("Robot %s simulation complete", robotID)
			return
		}
	}
}

// SimulatePersonUDP simulates a person client sending messages
func SimulatePersonUDP(serverAddr, personID string, duration time.Duration) {
	client, err := NewUDPClient(serverAddr, "person", personID)
	if err != nil {
		log.Fatalf("Failed to create person client: %v", err)
	}
	defer client.Close()

	// Start listening in background
	go client.Listen()

	// Simulate person sending messages
	messages := []string{
		"Hello_from_" + personID,
		"Status_check",
		"Test_message",
		"Delivery_inquiry",
		"Thanks!",
	}

	timeout := time.After(duration)
	messageIdx := 0

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if messageIdx < len(messages) {
				if err := client.Send(messages[messageIdx]); err != nil {
					log.Printf("Error sending packet: %v", err)
					return
				}
				messageIdx++
			}
		case <-timeout:
			log.Printf("Person %s simulation complete", personID)
			return
		}
	}
}
