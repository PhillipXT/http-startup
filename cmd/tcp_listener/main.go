package main

import (
	"fmt"
	"log"
	"net"

	"github.com/PhillipXT/http-startup/internal/request"
)

const filePath = "messages.txt"
const port = ":42069"

// Run the application, then print to the port from the command line:
// go run . | tee /tmp/tcp.txt
// printf "Can you hear me now?\r\n" | nc -w 1 127.0.0.1 42069
// nc -v localhost 42069

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err)
	}

	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}

		fmt.Println("Accepted connection from:", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("error reading from connection: %s\n", err)
		}

		line := req.RequestLine

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", line.Method)
		fmt.Printf("- Target: %s\n", line.RequestTarget)
		fmt.Printf("- Version: %s\n", line.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Println("Body:")
		fmt.Println(string(req.Body))

		fmt.Printf("Connection to %s closed\n", conn.RemoteAddr())
	}
}
