// Simple standalone demo - run this to see the servers and clients in action
// To run this file standalone:
// 1. Copy this file to a new directory
// 2. Run: go run demo_standalone.go
//
// Or call RunSimpleDemo() from your main.go

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// RunSimpleDemo runs a simple demonstration of TCP and UDP servers with clients
func RunSimpleDemo() {
	log.Println("=== TCP/UDP Network Demo ===")
	log.Println("Starting servers and test clients...")

	var wg sync.WaitGroup

	// Start TCP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		startDemoTCPServer()
	}()

	// Start UDP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		startDemoUDPServer()
	}()

	// Wait for servers to start
	time.Sleep(2 * time.Second)

	// Start demo clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		runDemoClients()
	}()

	wg.Wait()
}

func startDemoTCPServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("TCP server failed: %v", err)
	}
	defer listener.Close()
	log.Println("✅ TCP Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleDemoTCPConnection(conn)
	}
}

func handleDemoTCPConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Read client info
	data, _ := reader.ReadString('\n')
	log.Printf("📥 TCP: New client connected: %s", data)

	conn.Write([]byte("Welcome to TCP server!\n"))

	// Echo messages
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		log.Printf("📨 TCP received: %s", message)
		conn.Write([]byte(fmt.Sprintf("Echo: %s", message)))
	}
}

func startDemoUDPServer() {
	addr, _ := net.ResolveUDPAddr("udp", ":8081")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("UDP server failed: %v", err)
	}
	defer conn.Close()
	log.Println("✅ UDP Server listening on :8081")

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		message := string(buffer[:n])
		log.Printf("📨 UDP received: %s", message)
		conn.WriteToUDP([]byte("ACK"), clientAddr)
	}
}

func runDemoClients() {
	var wg sync.WaitGroup

	// TCP Robot Client
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Printf("Failed to connect TCP client: %v", err)
			return
		}
		defer conn.Close()

		conn.Write([]byte("robot:DemoBot-001\n"))
		go func() {
			reader := bufio.NewReader(conn)
			for {
				msg, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				log.Printf("🤖 Robot received: %s", msg)
			}
		}()

		for i := 0; i < 5; i++ {
			time.Sleep(2 * time.Second)
			msg := fmt.Sprintf("Position[%d,%d] Battery:%d%%\n", i*10, i*5, 100-i*5)
			conn.Write([]byte(msg))
			log.Printf("🤖 Robot sent: %s", msg)
		}
	}()

	// UDP Robot Client
	wg.Add(1)
	go func() {
		defer wg.Done()
		addr, _ := net.ResolveUDPAddr("udp", "localhost:8081")
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			log.Printf("Failed to connect UDP client: %v", err)
			return
		}
		defer conn.Close()

		for i := 0; i < 5; i++ {
			time.Sleep(2 * time.Second)
			msg := fmt.Sprintf("robot:UDPBot-001:Status_%d", i)
			conn.Write([]byte(msg))
			log.Printf("🤖 UDP Robot sent: %s", msg)

			buffer := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, _ := conn.Read(buffer)
			if n > 0 {
				log.Printf("🤖 UDP Robot received: %s", string(buffer[:n]))
			}
		}
	}()

	// TCP Person Client
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Printf("Failed to connect TCP client: %v", err)
			return
		}
		defer conn.Close()

		conn.Write([]byte("person:Alice\n"))
		go func() {
			reader := bufio.NewReader(conn)
			for {
				msg, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				log.Printf("👤 Alice received: %s", msg)
			}
		}()

		messages := []string{"Hello!", "How are you?", "Testing the server"}
		for _, msg := range messages {
			time.Sleep(3 * time.Second)
			conn.Write([]byte(msg + "\n"))
			log.Printf("👤 Alice sent: %s", msg)
		}
	}()

	wg.Wait()
	log.Println("\n✅ Demo completed! All clients finished.")
}
