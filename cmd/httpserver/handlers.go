package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sambakker4/httpfromtcp/internal/headers"
	"github.com/sambakker4/httpfromtcp/internal/request"
	"github.com/sambakker4/httpfromtcp/internal/response"
)

func handlerHTTPBin(w *response.Writer, req *request.Request) {
	endpoint := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

	resp, err := http.Get("https://httpbin.org/" + endpoint)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	err = w.WriteStatusLine(200)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	hdrs := response.GetDefaultHeaders(0)
	hdrs.Set("Transfer-Encoding", "chunked")
	hdrs.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	delete(hdrs, "content-length")

	err = w.WriteHeaders(hdrs)

	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fullResponse := ""

	buf := make([]byte, 1024)
	n := -1
	for n != 0 {
		n, err = resp.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("error: %v\n", err)
		}
		fullResponse += string(buf[:n])

		n, err = w.WriteChunkedBody(buf[:n])
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	hash := sha256.Sum256([]byte(fullResponse))
	length := len(fullResponse)

	headers := headers.NewHeaders()
	headers.Set("X-Content-SHA256", fmt.Sprintf("%x", hash))
	headers.Set("X-Content-Length", strconv.Itoa(length))

	err = w.WriteTrailers(headers)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

func handlerYourProblem(w *response.Writer, req *request.Request) {
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
}

func handlerMyProblem(w *response.Writer, req *request.Request) {
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
}

func handlerGetVideo(w *response.Writer, req *request.Request) {
	data, err := os.ReadFile("/home/sambakker/workspace/github.com/sambakker4/httpfromtcp/assets/vim.mp4")
	if err != nil {
		log.Printf("error: %s\n", err.Error())
	}
	headers := response.GetDefaultHeaders(len(data))
	headers.Set("Content-Type", "video/mp4")

	w.WriteStatusLine(200)
	w.WriteHeaders(headers)
	w.WriteBody(data)
}
