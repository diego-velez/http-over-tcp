package main

import (
	"fmt"
	"log"
	"net"

	"github.com/diego-velez/http-from-scratch-course/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close() // nolint: errcheck

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("New connection accepted!")

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range r.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
	}
}
