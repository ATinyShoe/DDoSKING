package cli

import (
	"fmt"
	"os"
	"strings"

	"c2/attack"
	"c2/bot"
	"c2/config"

	"github.com/chzyer/readline"
)

// HandleUserInput handles user command input
func HandleUserInput() {
	// Create a readline instance with command history support
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "C2> ",
		HistoryFile:     config.HistoryFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Println("Failed to initialize command line:", err)
		os.Exit(1)
	}
	defer rl.Close()

	// Set up auto-completion
	rl.Config.AutoComplete = setupAutoComplete()

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
			attack.HandleAttack(args[1:])
		case "list":
			ShowBotList()
		case "info":
			ShowBotInfo(args[1:])
		case "clear":
			ClearScreen()
			rl.SetPrompt("C2> ") // Reset prompt after clearing screen
		case "help":
			ShowHelp()
		case "stop":
			HandleStopCommand(args[1:])
		case "exit":
			CleanExit()
		default:
			fmt.Println("[!] Unknown command. Type 'help' to see available commands")
		}
	}
}

// setupAutoComplete sets up command auto-completion
func setupAutoComplete() *readline.PrefixCompleter {
	// Create completion items for "attack" subcommands
	attackItems := make([]readline.PrefixCompleterInterface, 0)
	for _, method := range config.GetAllMethods() {
		attackItems = append(attackItems, readline.PcItem(method))
	}

	// Create the main completer
	return readline.NewPrefixCompleter(
		readline.PcItem("attack", attackItems...),
		readline.PcItem("list"),
		readline.PcItem("info"),
		readline.PcItem("clear"),
		readline.PcItem("help"),
		readline.PcItem("stop"),
		readline.PcItem("exit"),
	)
}

// HandleStopCommand handles the stop command
func HandleStopCommand(args []string) {
	if len(args) == 0 {
		bot.SendStopToAllBots()
	} else {
		ip := args[0]
		bot.SendStopToSpecificBot(ip)
	}
}
