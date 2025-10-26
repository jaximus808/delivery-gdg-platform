package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

// TCPClientConnection represents a TCP client connection
type TCPClientConnection struct {
	conn       net.Conn
	clientType string
	clientID   string
}

// NewTCPClient creates a new TCP client
func NewTCPClient(serverAddr, clientType, clientID string) (*TCPClientConnection, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	client := &TCPClientConnection{
		conn:       conn,
		clientType: clientType,
		clientID:   clientID,
	}

	// Send client info to server
	_, err = conn.Write([]byte(fmt.Sprintf("%s:%s\n", clientType, clientID)))
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send client info: %v", err)
	}

	log.Printf("TCP Client connected as %s (Type: %s)", clientID, clientType)

	return client, nil
}

// Listen starts listening for messages from the server
func (c *TCPClientConnection) Listen() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection closed: %v", err)
			return
		}
		fmt.Printf("[Server]: %s", message)
	}
}

// Send sends a message to the server
func (c *TCPClientConnection) Send(message string) error {
	_, err := c.conn.Write([]byte(message + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

// Close closes the connection
func (c *TCPClientConnection) Close() {
	c.conn.Close()
}

// SimulateRobotTCP simulates a robot client sending periodic updates
func SimulateRobotTCP(serverAddr, robotID string, duration time.Duration) {
	client, err := NewTCPClient(serverAddr, "robot", robotID)
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
			status := fmt.Sprintf("Robot status: Position[%d,%d] Battery:%d%%",
				counter*10, counter*5, 100-counter*2)
			if err := client.Send(status); err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
		case <-timeout:
			log.Printf("Robot %s simulation complete", robotID)
			return
		}
	}
}

// SimulatePersonTCP simulates a person client sending messages
func SimulatePersonTCP(serverAddr, personID string, duration time.Duration) {
	client, err := NewTCPClient(serverAddr, "person", personID)
	if err != nil {
		log.Fatalf("Failed to create person client: %v", err)
	}
	defer client.Close()

	// Start listening in background
	go client.Listen()

	// Simulate person sending messages
	messages := []string{
		"Hello from " + personID,
		"How is everyone doing?",
		"This is a test message",
		"Checking the delivery status",
		"Goodbye!",
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
					log.Printf("Error sending message: %v", err)
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
