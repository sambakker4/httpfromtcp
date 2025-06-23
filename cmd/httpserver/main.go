package main

import (
	"errors"
	"io"
	"log"
	"net/http"
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
		endpoint := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

		resp, err := http.Get("https://httpbin.org/" + endpoint)
		if err != nil {
			log.Printf("error: %v\n", err)
		}

		err = w.WriteStatusLine(200)
		if err != nil {
			log.Printf("error: %v\n", err)
		}

		headers := response.GetDefaultHeaders(0)
		headers.Set("Transfer-Encoding", "chunked")
		delete(headers, "content-length")

		err = w.WriteHeaders(headers)

		if err != nil {
			log.Printf("error: %v\n", err)
		}

		buf := make([]byte, 1024)
		n := -1
		for n != 0 {
			_, err = resp.Body.Read(buf)
			if err != nil {
				if errors.Is(io.EOF, err) {
					break
				}
				log.Printf("error: %v\n", err)
			}

			n, err = w.WriteChunkedBody(buf)
			if err != nil {
				log.Printf("error: %v\n", err)
			}
		}

		_, err = w.WriteChunkedBodyDone()
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		return
	}
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		html := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
		headers := response.GetDefaultHeaders(len([]byte(html)))
		headers.Set("Content-type", "text/html")
		w.WriteStatusLine(400)
		w.WriteHeaders(headers)
		w.WriteBody([]byte(html))
	case "/myproblem":
		html := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`
		headers := response.GetDefaultHeaders(len([]byte(html)))
		headers.Set("Content-type", "text/html")
		w.WriteStatusLine(500)
		w.WriteHeaders(headers)
		w.WriteBody([]byte(html))
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
