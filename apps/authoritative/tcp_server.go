package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// TCPClient represents a connected TCP client
type TCPClient struct {
	ID          string
	Conn        net.Conn
	ClientType  string // "robot" or "person"
	ConnectedAt time.Time
}

// TCPServer manages TCP connections
type TCPServer struct {
	clients map[string]*TCPClient
	mu      sync.RWMutex
	port    string
}

// NewTCPServer creates a new TCP server instance
func NewTCPServer(port string) *TCPServer {
	return &TCPServer{
		clients: make(map[string]*TCPClient),
		port:    port,
	}
}

// Start begins listening for TCP connections
func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("TCP Server listening on %s", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection manages individual client connections
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("New TCP connection from %s", clientAddr)

	// Read client type and ID
	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading from client %s: %v", clientAddr, err)
		return
	}

	// Parse client info (format: "TYPE:ID\n")
	var clientType, clientID string
	fmt.Sscanf(data, "%s:%s", &clientType, &clientID)

	// Register client
	client := &TCPClient{
		ID:          clientID,
		Conn:        conn,
		ClientType:  clientType,
		ConnectedAt: time.Now(),
	}

	s.addClient(client)
	defer s.removeClient(clientID)

	// Send welcome message
	welcome := fmt.Sprintf("Welcome %s (Type: %s)! You are connected.\n", clientID, clientType)
	conn.Write([]byte(welcome))

	// Handle messages from client
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected: %v", clientID, err)
			break
		}

		log.Printf("[TCP] Received from %s (%s): %s", clientID, clientType, message)

		// Echo message back
		response := fmt.Sprintf("Server received: %s", message)
		conn.Write([]byte(response))

		// Broadcast to all other clients
		s.broadcast(clientID, fmt.Sprintf("[%s]: %s", clientID, message))
	}
}

// addClient adds a client to the server's client map
func (s *TCPServer) addClient(client *TCPClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client.ID] = client
	log.Printf("TCP Client registered: %s (Type: %s) - Total clients: %d",
		client.ID, client.ClientType, len(s.clients))
}

// removeClient removes a client from the server's client map
func (s *TCPServer) removeClient(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, clientID)
	log.Printf("TCP Client removed: %s - Total clients: %d", clientID, len(s.clients))
}

// broadcast sends a message to all connected clients except the sender
func (s *TCPServer) broadcast(senderID, message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for id, client := range s.clients {
		if id != senderID {
			client.Conn.Write([]byte(message))
		}
	}
}

// GetClients returns a snapshot of all connected clients
func (s *TCPServer) GetClients() map[string]*TCPClient {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make(map[string]*TCPClient)
	for k, v := range s.clients {
		clients[k] = v
	}
	return clients
}

// GetClientCount returns the number of connected clients
func (s *TCPServer) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}
