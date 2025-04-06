package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"auth/config"
	"auth/handler"
	"auth/server"
	"auth/ui"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Create DNS handler
	dnsHandler := handler.NewDNSHandler(cfg)

	// Create TCP and UDP servers
	tcpServer := server.NewTCPServer(cfg, dnsHandler)
	udpServer := server.NewUDPServer(cfg, dnsHandler)

	// Create user interface
	userInterface := ui.NewUI(dnsHandler)

	// Set up signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Start servers
	go func() {
		if err := tcpServer.Start(); err != nil {
			log.Fatalf("Failed to start TCP server: %v", err)
		}
	}()

	go func() {
		if err := udpServer.Start(); err != nil {
			log.Fatalf("Failed to start UDP server: %v", err)
		}
	}()

	// Start UI in a goroutine
	go userInterface.Run()

	// Wait for signal
	<-signalCh

	// Shutdown
	tcpServer.Shutdown()
	udpServer.Shutdown()
	userInterface.Stop()
}