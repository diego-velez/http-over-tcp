package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/diego-velez/http-from-scratch-course/internal/request"
	"github.com/diego-velez/http-from-scratch-course/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

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

	var buf bytes.Buffer
	buf.Write([]byte("\r\n"))
	handlerErr := s.handler(&buf, r)
	if handlerErr != nil {
		_ = response.WriteStatusLine(conn, handlerErr.Code)
		_, _ = buf.Write([]byte(handlerErr.Msg))
	} else {
		_ = response.WriteStatusLine(conn, response.StatusOK)
	}

	headers := response.GetDefaultHeaders(len(buf.Bytes()) - 2)
	_ = response.WriteHeaders(conn, headers)

	_, _ = conn.Write(buf.Bytes())
}
