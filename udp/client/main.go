package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	maxRetries    = 3
	retryInterval = 500 * time.Millisecond
)

func main() {
	// Resolve server address
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	if err != nil {
		log.Fatalf("Failed to resolve server address: %v", err)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	sequence := uint32(0)
	buffer := make([]byte, 1024)

	for {
		// Create message with sequence number
		message := fmt.Sprintf("%04d:Hello, UDP server! Msg #%d", sequence, sequence)

		// Implement reliable delivery with retries
		success := false
		for retry := 0; retry < maxRetries; retry++ {
			// Send message
			_, err := conn.Write([]byte(message))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				continue
			}

			// Set read deadline
			conn.SetReadDeadline(time.Now().Add(retryInterval))

			// Wait for ACK
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					fmt.Printf("Timeout, retrying... (%d)\n", retry+1)
					continue
				}
				log.Printf("Error reading response: %v", err)
				continue
			}

			response := string(buffer[:n])
			if strings.HasPrefix(response, "ACK") {
				success = true
				fmt.Printf("Message %d delivered successfully\n", sequence)
				break
			} else if strings.HasPrefix(response, "NACK") {
				fmt.Printf("Received NACK for message %d\n", sequence)
				continue
			}
		}

		if !success {
			fmt.Printf("Failed to deliver message %d after %d attempts\n",
				sequence, maxRetries)
		}

		sequence++
		time.Sleep(time.Second)
	}
}
