package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {

	//  Create a TCP listener on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to start  server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server listening on: 8080")

	for {
		//  Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		//    Handle each connection in a go routine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	//  create a buffer for incoming data
	reader := bufio.NewReader(conn)

	for {
		// Read incoming message
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			return
		}

		fmt.Printf("Received: %s", message)

		//  Echo message back
		conn.Write([]byte("Server received: " + message))
	}
}
