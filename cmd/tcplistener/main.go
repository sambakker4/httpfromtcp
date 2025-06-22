package main

import (
	"fmt"
	"github.com/sambakker4/httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection accepted")

		req, err := request.RequestFromReader(connection)
		if err != nil {
			log.Printf("error: %s", err.Error())
		}

		fmt.Println("Request Line:")
		fmt.Println(" - Method:", req.RequestLine.Method)
		fmt.Println(" - Target:", req.RequestLine.RequestTarget)
		fmt.Println(" - Version:", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf(" - %s: %s\n", key, value)
		}

		fmt.Println("Body:")
		fmt.Println(string(req.Body))
		fmt.Println()
	}
}
