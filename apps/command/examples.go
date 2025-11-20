// Example: How to use the TCP/UDP servers and clients
//
// This file demonstrates various ways to run the network demo.
// Uncomment the desired section in your main.go to use it.

package main

import (
	"sync"
	"time"
)

// Example1: Run both servers with test clients
func Example1_ServersWithTestClients() {
	var wg sync.WaitGroup

	// Start servers
	wg.Add(1)
	go func() {
		defer wg.Done()
		RunServers()
	}()

	// Start test clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		RunAllTestClients()
	}()

	wg.Wait()
}

// Example2: Run only TCP server with TCP test clients
func Example2_TCPServerWithClients() {
	var wg sync.WaitGroup

	// Start TCP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		tcpServer := NewTCPServer(":8080")
		tcpServer.Start()
	}()

	// Give server time to start
	time.Sleep(1 * time.Second)

	// Start TCP test clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		RunTCPTestClients()
	}()

	wg.Wait()
}

// Example3: Run only UDP server with UDP test clients
func Example3_UDPServerWithClients() {
	var wg sync.WaitGroup

	// Start UDP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		udpServer := NewUDPServer(":8081")
		udpServer.Start()
	}()

	// Give server time to start
	time.Sleep(1 * time.Second)

	// Start UDP test clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		RunUDPTestClients()
	}()

	wg.Wait()
}

// Example4: Run single robot client (connect to already running server)
func Example4_SingleRobotClient() {
	SimulateRobotTCP("localhost:8080", "MyRobot-001", 60*time.Second)
}

// Example5: Run single person client (connect to already running server)
func Example5_SinglePersonClient() {
	SimulatePersonTCP("localhost:8080", "Alice", 60*time.Second)
}

// Example6: Custom client with manual control
func Example6_CustomClient() {
	// Create a custom TCP client
	client, err := NewTCPClient("localhost:8080", "robot", "CustomBot-123")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Start listening in background
	go client.Listen()

	// Send custom messages
	client.Send("Initializing systems...")
	time.Sleep(2 * time.Second)

	client.Send("Position: [0,0,0]")
	time.Sleep(2 * time.Second)

	client.Send("Task completed!")
}
