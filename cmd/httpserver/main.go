package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diego-velez/http-from-scratch-course/internal/request"
	"github.com/diego-velez/http-from-scratch-course/internal/response"
	"github.com/diego-velez/http-from-scratch-course/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handleConn)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handleConn(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{Code: response.StatusBadRequest, Msg: "Your problem is not my problem\n"}
	case "/myproblem":
		return &server.HandlerError{Code: response.StatusInternalServerError, Msg: "Woopsie, my bad\n"}
	default:
		_, _ = w.Write([]byte("All good, frfr\n"))
	}

	return nil
}
