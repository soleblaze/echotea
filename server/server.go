// Package main contains the server for the echo client.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	addr := "localhost:9999"

	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer server.Close()

	log.Println("Server is running on:", addr)

	conn, err := server.Accept()
	if err != nil {
		log.Println("Failed to accept conn.", err)
	}

	for {
		buf := make([]byte, 512)

		_, err = conn.Read(buf)
		if err == io.EOF {
			return
		}

		if err != nil {
			fmt.Println("Error reading:")
			fmt.Println(err)

			continue
		}

		log.Printf("Received: %s", buf)

		i, err := conn.Write([]byte(string(buf) + "\n"))
		if err != nil {
			fmt.Println("Error writing:")
			fmt.Println(err)
		}

		fmt.Printf("Wrote Message: %d\n", i)
	}
}
