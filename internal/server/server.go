package server

import (
	"fmt"
	"log"
	"net"

	"github.com/diego-velez/http-from-scratch-course/internal/request"
	"github.com/diego-velez/http-from-scratch-course/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	Code response.StatusCode
	Msg  string
}

type Server struct {
	listener net.Listener
	handler  Handler
	closed   bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{listener: listener, handler: handler, closed: false}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed = true
	_ = s.listener.Close()
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		if s.closed {
			return
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close() // nolint: errcheck

	fmt.Println("New connection accepted!")

	r, err := request.RequestFromReader(conn)
	if err != nil {
		_ = response.WriteStatusLine(conn, response.StatusBadRequest)
		headers := response.GetDefaultHeaders(0)
		_ = response.WriteHeaders(conn, headers)
		return
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

	var w response.Writer
	s.handler(&w, r)

	_, _ = conn.Write(w.Buf.Bytes())
}
