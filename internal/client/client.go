package client

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cmingou/tcp-hb/internal/protocol"
)

type HeartbeatClient struct {
	Addr     string
	Port     int
	Internal int
	Timeout  int
}

func NewHeartbeatClient(addr string, port, interval, timeout int) *HeartbeatClient {
	return &HeartbeatClient{
		Addr:     addr,
		Port:     port,
		Internal: interval,
		Timeout:  timeout,
	}
}

func (c *HeartbeatClient) Connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Addr, c.Port))
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	ticker := time.NewTicker(time.Duration(c.Internal) * time.Second)
	defer ticker.Stop()

	var (
		heartbeatID uint32 = 1
		totalSent   uint64
		totalLost   uint64
	)
	acks := sync.Map{}
	done := make(chan struct{})

	go c.receiveAcks(conn, &acks, &totalLost, done)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		c.printStats(totalSent, totalLost)
		close(done)
		os.Exit(0)
	}()

	for range ticker.C {
		select {
		case <-done:
			return
		default:

			buf := make([]byte, 13) // Type (1) + Heartbeat ID (4) + Timestamp (8)
			buf[0] = protocol.TypeHeartbeatRequest
			binary.BigEndian.PutUint32(buf[1:5], heartbeatID)
			binary.BigEndian.PutUint64(buf[5:13], uint64(time.Now().UnixNano()))

			_, err := conn.Write(buf)
			if err != nil {
				fmt.Println("Failed to send heartbeat:", err)
				return
			}

			acks.Store(heartbeatID, time.Now())
			atomic.AddUint64(&totalSent, 1)

			// Set up a timeout for this heartbeat
			go func(id uint32) {
				time.Sleep(time.Duration(c.Timeout) * time.Second)
				if _, ok := acks.Load(id); ok {
					fmt.Printf("Packet loss or timeout for ID=%v\n", id)
					atomic.AddUint64(&totalLost, 1)
					acks.Delete(id)
				}
			}(heartbeatID)

			heartbeatID++
		}
	}

	<-done
}

func (c *HeartbeatClient) receiveAcks(conn net.Conn, acks *sync.Map, totalLost *uint64, done chan struct{}) {
	buf := make([]byte, 13) // Type (1) + Heartbeat ID (4) + Timestamp (8)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, err := conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			} else {
				fmt.Println("Server stopped responding:", err)
				close(done)
				return
			}
		}

		if buf[0] == protocol.TypeHeartbeatResponse {
			responseID := binary.BigEndian.Uint32(buf[1:5])
			if value, ok := acks.Load(responseID); ok {
				sendTime := value.(time.Time)
				rtt := time.Since(sendTime)
				fmt.Printf("Received valid heartbeat response: ID=%d, RTT=%v\n", responseID, rtt/2)
				acks.Delete(responseID)
			} else {
				fmt.Printf("Received unexpected heartbeat response: ID=%d\n", responseID)
			}
		} else {
			fmt.Println("Received invalid message type")
		}
	}
}

func (c *HeartbeatClient) printStats(totalSent, totalLost uint64) {
	fmt.Printf("\nTotal heartbeats sent: %d\n", totalSent)
	fmt.Printf("Total packet loss: %d\n", totalLost)
	if totalSent > 0 {
		fmt.Printf("Packet loss rate: %.2f%%\n", float64(totalLost*100)/float64(totalSent))
	}
}
