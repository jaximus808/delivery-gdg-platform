package main

import (
	"log"
	"sync"
	"time"
)

// RunTCPTestClients runs fake robot and person clients for TCP testing
func RunTCPTestClients() {
	var wg sync.WaitGroup

	// Give server time to start
	time.Sleep(2 * time.Second)

	log.Println("Starting TCP test clients...")

	// Start 2 robot clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		SimulateRobotTCP("localhost:8080", "Robot-001", 30*time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		SimulateRobotTCP("localhost:8080", "Robot-002", 30*time.Second)
	}()

	// Start 2 person clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		SimulatePersonTCP("localhost:8080", "Person-Alice", 30*time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		SimulatePersonTCP("localhost:8080", "Person-Bob", 30*time.Second)
	}()

	wg.Wait()
	log.Println("All TCP test clients completed")
}

// RunUDPTestClients runs fake robot and person clients for UDP testing
func RunUDPTestClients() {
	var wg sync.WaitGroup

	// Give server time to start
	time.Sleep(2 * time.Second)

	log.Println("Starting UDP test clients...")

	// Start 2 robot clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		SimulateRobotUDP("localhost:8081", "Robot-UDP-001", 30*time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		SimulateRobotUDP("localhost:8081", "Robot-UDP-002", 30*time.Second)
	}()

	// Start 2 person clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		SimulatePersonUDP("localhost:8081", "Person-UDP-Charlie", 30*time.Second)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		SimulatePersonUDP("localhost:8081", "Person-UDP-Diana", 30*time.Second)
	}()

	wg.Wait()
	log.Println("All UDP test clients completed")
}

// RunAllTestClients runs both TCP and UDP test clients
func RunAllTestClients() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunTCPTestClients()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunUDPTestClients()
	}()

	wg.Wait()
	log.Println("All test clients completed")
}
