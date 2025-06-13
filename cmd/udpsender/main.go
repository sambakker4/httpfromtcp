package main

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"net"
)

func main(){
	address, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		
		input, err := reader.ReadString('\n')
		_, err = connection.Write([]byte(input))
			
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())	
		}
	}
}
