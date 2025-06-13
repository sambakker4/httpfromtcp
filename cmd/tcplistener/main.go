package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		linesChannel := getLinesChannel(connection)
		
		for item := range linesChannel {
			fmt.Printf("read: %s\n", item)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)

	go func(){
		defer f.Close()
		defer close(channel)

		currentLine := ""
		for {
			data := make([]byte, 8, 8)
			n, err := f.Read(data)
			if err != nil {
				if currentLine != "" {
					channel <- currentLine
				}

				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Println("error:", err.Error())
			}

			parts := strings.Split(string(data[:n]), "\n")

			for i, part := range parts {
				if i == len(parts) - 1 {
					currentLine +=  part
					break
				}	

				channel <- (currentLine + part)
				currentLine = ""
			}
		}
	}()

	return channel
}
