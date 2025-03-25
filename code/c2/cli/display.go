package cli

import (
	"fmt"
	"os"

	"c2/bot"
	"c2/config"
)

// ShowBanner displays the program startup banner
func ShowBanner() {
	fmt.Println(config.GetBanner())
}

// ShowBotList displays the list of bots
func ShowBotList() {
	bot.BotMu.Lock()
	defer bot.BotMu.Unlock()

	fmt.Println("\nActive Bot List:")
	fmt.Printf("%-18s %-22s %-22s %-8s %s\n",
		"IP Address", "Connect Time", "Last Active", "Status", "Current Task")

	for ip, bot := range bot.Bots {
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
		} else if len(task) > 50 {
			task = task[:47] + "..."
		}

		fmt.Printf("%-18s %-22s %-22s %-8s %s\n",
			ip,
			bot.ConnectTime.Format("2006-01-02 15:04:05"),
			bot.LastActive.Format("2006-01-02 15:04:05"),
			status,
			task)
	}
}

// ShowHelp displays help information
func ShowHelp() {
	fmt.Println(config.GetCommandHelp())
}

// ShowBotInfo displays detailed information about a specific bot
func ShowBotInfo(args []string) {
	if len(args) == 0 {
		fmt.Println("[!] Please specify bot IP")
		return
	}

	botIP := args[0]
	bot.BotMu.Lock()
	defer bot.BotMu.Unlock()

	if b, exists := bot.Bots[botIP]; exists {
		fmt.Printf("\nBot Details - %s\n", botIP)
		fmt.Println("--------------------------------------")
		fmt.Printf("Status:       %s\n", b.Status)
		fmt.Printf("Connect Time: %s\n", b.ConnectTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("Last Active:  %s\n", b.LastActive.Format("2006-01-02 15:04:05"))
		if b.Status == "Offline" {
			fmt.Printf("Disconnect:   %s\n", b.DisconnectTime.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("Current Task: %s\n", b.CurrentTask)
		fmt.Printf("Connection:   %v\n", b.Conn)
		fmt.Println("--------------------------------------")
	} else {
		fmt.Println("[!] Specified bot not found")
	}
}

// ClearScreen clears the screen
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// CleanExit exits the program cleanly
func CleanExit() {
	fmt.Println("\n[!] Disconnecting all bots...")
	bot.BotMu.Lock()
	for _, b := range bot.Bots {
		b.Conn.Close()
	}
	bot.BotMu.Unlock()
	os.Exit(0)
}
