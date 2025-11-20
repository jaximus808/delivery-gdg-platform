#!/bin/bash

echo "=== TCP/UDP Network Demo Test Script ==="
echo ""
echo "Choose an option:"
echo "1. Run simple standalone demo"
echo "2. Run full demo with all features"
echo "3. Run TCP server only"
echo "4. Run UDP server only"
echo "5. Run test clients (requires servers running)"
echo ""
read -p "Enter option (1-5): " option

case $option in
  1)
    echo "Running simple demo..."
    go run demo_standalone.go
    ;;
  2)
    echo "Running full demo..."
    go run tcp_server.go udp_server.go tcp_client.go udp_client.go server_main.go test_clients.go examples.go network_demo.go -mode all
    ;;
  3)
    echo "Starting TCP server on :8080..."
    go run tcp_server.go server_main.go
    ;;
  4)
    echo "Starting UDP server on :8081..."
    go run udp_server.go server_main.go
    ;;
  5)
    echo "Running test clients..."
    go run tcp_client.go udp_client.go test_clients.go
    ;;
  *)
    echo "Invalid option"
    exit 1
    ;;
esac
