package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// UDPClient represents a UDP client
type UDPClient struct {
	ID         string
	Addr       *net.UDPAddr
	ClientType string // "robot" or "person"
	LastSeen   time.Time
}

// UDPServer manages UDP connections
type UDPServer struct {
	clients map[string]*UDPClient
	mu      sync.RWMutex
	port    string
	conn    *net.UDPConn
}

// NewUDPServer creates a new UDP server instance
func NewUDPServer(port string) *UDPServer {
	return &UDPServer{
		clients: make(map[string]*UDPClient),
		port:    port,
	}
}

// Start begins listening for UDP packets
func (s *UDPServer) Start() error {
	addr, err := net.ResolveUDPAddr("udp", s.port)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to start UDP server: %v", err)
	}
	defer conn.Close()

	s.conn = conn
	log.Printf("UDP Server listening on %s", s.port)

	// Start cleanup routine for inactive clients
	go s.cleanupInactiveClients()

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP packet: %v", err)
			continue
		}

		go s.handlePacket(buffer[:n], clientAddr)
	}
}

// handlePacket processes incoming UDP packets
func (s *UDPServer) handlePacket(data []byte, addr *net.UDPAddr) {
	message := string(data)

	// Parse message (format: "TYPE:ID:MESSAGE")
	var clientType, clientID, msg string
	n, _ := fmt.Sscanf(message, "%s:%s:%s", &clientType, &clientID, &msg)

	if n < 2 {
		log.Printf("Invalid UDP packet format from %s", addr.String())
		return
	}

	// Update or add client
	s.updateClient(clientID, clientType, addr)

	log.Printf("[UDP] Received from %s (%s): %s", clientID, clientType, msg)

	// Send acknowledgment
	response := fmt.Sprintf("ACK:%s", clientID)
	s.conn.WriteToUDP([]byte(response), addr)

	// Broadcast to all other clients
	s.broadcast(clientID, fmt.Sprintf("[%s]: %s", clientID, msg))
}

// updateClient updates or adds a client to the map
func (s *UDPServer) updateClient(clientID, clientType string, addr *net.UDPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, exists := s.clients[clientID]; exists {
		client.LastSeen = time.Now()
	} else {
		s.clients[clientID] = &UDPClient{
			ID:         clientID,
			Addr:       addr,
			ClientType: clientType,
			LastSeen:   time.Now(),
		}
		log.Printf("UDP Client registered: %s (Type: %s) - Total clients: %d",
			clientID, clientType, len(s.clients))
	}
}

// broadcast sends a message to all connected clients except the sender
func (s *UDPServer) broadcast(senderID, message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for id, client := range s.clients {
		if id != senderID {
			s.conn.WriteToUDP([]byte(message), client.Addr)
		}
	}
}

// cleanupInactiveClients removes clients that haven't been seen in 30 seconds
func (s *UDPServer) cleanupInactiveClients() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, client := range s.clients {
			if now.Sub(client.LastSeen) > 30*time.Second {
				delete(s.clients, id)
				log.Printf("UDP Client timeout: %s - Total clients: %d", id, len(s.clients))
			}
		}
		s.mu.Unlock()
	}
}

// GetClients returns a snapshot of all connected clients
func (s *UDPServer) GetClients() map[string]*UDPClient {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make(map[string]*UDPClient)
	for k, v := range s.clients {
		clients[k] = v
	}
	return clients
}

// GetClientCount returns the number of connected clients
func (s *UDPServer) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}
