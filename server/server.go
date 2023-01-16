// Package main contains the server for the echo client.
package main

import (
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
		server.Close()
		log.Fatal(err)
	}

	for {
		buf := make([]byte, 512)

		_, err = conn.Read(buf)
		if err != nil {
			log.Printf("Error reading: %s\n", err)
		}

		log.Printf("Received: %s\n", buf)

		_, err := conn.Write([]byte(string(buf) + "\n"))
		if err != nil {
			log.Printf("Error writing: %s\n", err)
		}
	}
}
