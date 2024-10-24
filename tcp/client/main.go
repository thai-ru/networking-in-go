package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	//	 Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	//	Read user input and send to server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		//	Send user input to server
		fmt.Fprintf(conn, "%s\n", scanner.Text())

		//	Read servers response
		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Server response: ", response)
	}
}
