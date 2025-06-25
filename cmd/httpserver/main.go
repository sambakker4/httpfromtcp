package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sambakker4/httpfromtcp/internal/request"
	"github.com/sambakker4/httpfromtcp/internal/response"
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

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		handlerHTTPBin(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/video" && req.RequestLine.Method == "GET" {
		handlerGetVideo(w, req)
		return
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handlerYourProblem(w, req)
	case "/myproblem":
		handlerMyProblem(w, req)
	default:
		html := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`
		headers := response.GetDefaultHeaders(len([]byte(html)))
		headers.Set("Content-type", "text/html")
		w.WriteStatusLine(200)
		w.WriteHeaders(headers)
		w.WriteBody([]byte(html))
	}
}
