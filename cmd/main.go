package main

import (
	"log"
	"github.com/bryanwahyu/test-jt/internal/server"
)

func main() {
	// Start the server
	if err := server.StartServer(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
