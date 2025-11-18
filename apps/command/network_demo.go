package main

import (
	"flag"
	"log"
	"sync"
)

// RunNetworkDemo is the main entry point for the network demo
func RunNetworkDemo() {
	mode := flag.String("mode", "all", "Mode: 'server', 'client', or 'all'")
	flag.Parse()

	var wg sync.WaitGroup

	switch *mode {
	case "server":
		log.Println("Running in SERVER mode")
		RunServers()

	case "client":
		log.Println("Running in CLIENT mode")
		RunAllTestClients()

	case "all":
		log.Println("Running in ALL mode (servers + test clients)")

		// Start servers
		wg.Add(1)
		go func() {
			defer wg.Done()
			RunServers()
		}()

		// Start test clients after a delay
		wg.Add(1)
		go func() {
			defer wg.Done()
			RunAllTestClients()
		}()

		wg.Wait()

	default:
		log.Fatalf("Invalid mode: %s. Use 'server', 'client', or 'all'", *mode)
	}
}
