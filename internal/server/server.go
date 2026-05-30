package server

import (
	"fmt"
	"log"
	"net"

	"github.com/diego-velez/http-from-scratch-course/internal/request"
)

type Server struct {
	listener net.Listener
	closed   bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{listener: listener, closed: false}
	s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed = true
	_ = s.listener.Close()
	return nil
}

func (s *Server) listen() {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			if s.closed {
				return
			}

			go s.handle(conn)
		}
	}()
}

func (s *Server) handle(conn net.Conn) {
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
	fmt.Println("Body:")
	fmt.Printf("%s\n", string(r.Body))

	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n")
	_, _ = conn.Write(out)
	_ = conn.Close()
}
