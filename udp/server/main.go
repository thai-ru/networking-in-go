package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	maxPacketSize = 1024
	serverAddr    = ":8081"
)

// Message represents our protocol
type Message struct {
	SequenceNumber uint32
	Payload        string
	TimeStamp      time.Time
}

func main() {
	//	Create UDP address
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to resolve address: %v", err)
	}

	//	Create UDP listener
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Printf("UDP server listening on %s\n", serverAddr)

	//	Create a map to track client sequences
	clientSequences := make(map[string]uint32)

	buffer := make([]byte, maxPacketSize)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		//	Handle packet in a goroutine
		go handlePacket(conn, remoteAddr, buffer[:n], clientSequences)

	}
}

func handlePacket(conn *net.UDPConn, addr *net.UDPAddr, data []byte, clientSequences map[string]uint32) {
	clientKey := addr.String()
	expectedSeq := clientSequences[clientKey]

	//	Simple packet parser (not for real-world use case, use proper serialization)
	sequenceStr := string(data[:4]) // first 4 bytes as a sequence
	payload := string(data[4:])

	var sequence uint32
	fmt.Sscanf(sequenceStr, "%04d", &sequence)

	//	Check packet loss or reordering
	if sequence != expectedSeq {
		//	Send NACK(Negative Acknowledgement)
		response := fmt.Sprintf("NACK:%d", expectedSeq)
		conn.WriteToUDP([]byte(response), addr)
		return
	}

	//	Send ACK (Acknowledgement)
	response := fmt.Sprintf("ACK: %d", sequence)
	conn.WriteToUDP([]byte(response), addr)

	//	 update expected sequence
	clientSequences[clientKey] = sequence + 1

	fmt.Printf("Received from %v: seq=%d, payload=%s\n", addr, sequence, payload)
}
