package main

import (
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

func handleConn(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusBadRequest)
		w.WriteBody([]byte("Your problem is not my problem\n"))
	case "/myproblem":
		w.WriteStatusLine(response.StatusInternalServerError)
		w.WriteBody([]byte("Woopsie, my bad\n"))
	default:
		w.WriteBody([]byte("All good, frfr\n"))
	}
}
