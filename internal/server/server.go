package server

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/cmingou/tcp-hb/internal/protocol"
)

type HeartbeatServer struct {
	Addr string
	Port int
}

func NewServer(addr string, port int) *HeartbeatServer {
	return &HeartbeatServer{
		Addr: addr,
		Port: port,
	}
}

func (s *HeartbeatServer) Start() {
	addr := fmt.Sprintf("%s:%d", s.Addr, s.Port)
	// Start a TCP server with addr and port, and listen for incoming connections. This server will print the TCP payload to stdout.
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening on %s: %v\n", addr, err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("Server listening on %s\n", addr)

	// Accept connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}

		// Handle the connection in a new goroutine.
		go handleConnection(conn)
	}
}

// handleConnection reads data from the connection and prints it to stdout
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Extract and print the source IP and port
	if tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		fmt.Printf("Accepted connection. Source IP: %s, Source Port: %d\n", tcpAddr.IP.String(), tcpAddr.Port)
	} else {
		fmt.Fprintf(os.Stderr, "Failed to get remote address information\n")
	}

	buf := make([]byte, 13) // Type (1) + Heartbeat ID (4) + Timestamp (8)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		heartbeatID := binary.BigEndian.Uint32(buf[1:5])
		timestamp := binary.BigEndian.Uint64(buf[5:13])
		rtt := time.Since(time.Unix(0, int64(timestamp)))

		fmt.Printf("Received heartbeat request: ID=%d, RTT=%v\n", heartbeatID, rtt)

		if heartbeatID%5 == 0 {
			// Simulate a dropped packet
			fmt.Println("Dropping packet")
			continue
		}

		// Send heartbeat response
		buf[0] = protocol.TypeHeartbeatResponse
		_, err = conn.Write(buf)
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}
