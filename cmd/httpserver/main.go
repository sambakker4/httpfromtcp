package main

import (
	"log"
	"os"
	"os/signal"
	"io"
	"syscall"
	"github.com/sambakker4/httpfromtcp/internal/request"
	"github.com/sambakker4/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
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

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget{
	case "/yourproblem":
		return &server.HandlerError{
			Message: "Your problem is not my problem\n",
			StatusCode: 400,
		}
	case "/myproblem":
		return &server.HandlerError{
			Message: "Woopsie, my bad\n",
			StatusCode: 500,
		}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
}
