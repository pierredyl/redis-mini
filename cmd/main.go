package main

import (
	"log"
	"net"
	"redis-mini/internal/data"
	"redis-mini/internal/handlers"
)

func main() {
	// Create concurrent TCP servers

	// Start a TCP server
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Failed to start the TCP server.")
	}
	log.Println("listening on :8080")

	// Make the database
	data := data.NewStore()

	// Load data from AOF file
	/*
		err = handlers.HandleAOFRead()
		if err != nil {
			fmt.Println("Couldn't read AOF file")
		}
	*/

	// Accepting connections and sending to its each goroutine
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Failed to accept new connection")
			return
		}
		log.Println("accepted connection from", conn.RemoteAddr())
		go handlers.HandleConnection(conn, data)

	}
}
