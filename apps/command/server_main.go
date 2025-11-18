package main

import (
	"log"
	"sync"
	"time"
)

// RunServers starts both TCP and UDP servers
func RunServers() {
	var wg sync.WaitGroup

	// Start TCP Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		tcpServer := NewTCPServer(":8080")
		log.Println("Starting TCP Server on :8080")
		if err := tcpServer.Start(); err != nil {
			log.Fatalf("TCP Server error: %v", err)
		}
	}()

	// Start UDP Server
	wg.Add(1)
	go func() {
		defer wg.Done()
		udpServer := NewUDPServer(":8081")
		log.Println("Starting UDP Server on :8081")
		if err := udpServer.Start(); err != nil {
			log.Fatalf("UDP Server error: %v", err)
		}
	}()

	// Monitor server status
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("=== Server Status ===")
			log.Printf("Servers running...")
		}
	}()

	wg.Wait()
}
