package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflector/sender"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

// ServerStatus defines the status of each server
type ServerStatus struct {
	PacketsReceived uint64
	LastActivity    time.Time
	IsRunning       bool
}

// Global variables
var (
	servers = map[string]*ServerStatus{
		"DNS":       {},
		"NTP":       {},
		"RDP":       {},
		"SSDP":      {},
		"SNMP":      {},
		"CHARGEN":   {},
		"OPENVPN":   {},
		"CLDAP":     {},
		"MEMCACHED": {},
	}

	ports = map[string]int{
		"DNS":       53,
		"NTP":       123,
		"RDP":       3389,
		"SSDP":      1900,
		"SNMP":      161,
		"CHARGEN":   19,
		"OPENVPN":   1194,
		"CLDAP":     389,
		"MEMCACHED": 11211,
	}

	statColors = color.New(color.FgHiWhite, color.BgBlack)
)

func main() {
	// Clear screen and display banner at start
	clearScreen()
	printBanner()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go signalHandler(cancel)

	// Main loop for menu interaction
	for {
		printMenu()
		handleInput(ctx)
	}
}

func signalHandler(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	cancel()
	os.Exit(0)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printBanner() {
	banner := `
██████╗ ███████╗███████╗██╗     ███████╗ ██████╗████████╗ ██████╗ ██████╗ 
██╔══██╗██╔════╝██╔════╝██║     ██╔════╝██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗
██████╔╝█████╗  █████╗  ██║     █████╗  ██║        ██║   ██║   ██║██████╔╝
██╔══██╗██╔══╝  ██╔══╝  ██║     ██╔══╝  ██║        ██║   ██║   ██║██╔══██╗
██║  ██║███████╗██║     ███████╗███████╗╚██████╗   ██║   ╚██████╔╝██║  ██║
╚═╝  ╚═╝╚══════╝╚═╝     ╚══════╝╚══════╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝
                                                                           
██████╗ ██████╗  ██████╗ ███████╗                                          
██╔══██╗██╔══██╗██╔═══██╗██╔════╝                                          
██║  ██║██║  ██║██║   ██║███████╗                                          
██║  ██║██║  ██║██║   ██║╚════██║                                          
██████╔╝██████╔╝╚██████╔╝███████║                                          
╚═════╝ ╚═════╝  ╚═════╝ ╚══════╝                                          
`
	fmt.Println(color.HiCyanString(banner))
	fmt.Println(color.HiYellowString("DDoS 反射放大攻击模拟器"))
	fmt.Println(color.HiWhiteString("支持协议: DNS, NTP, RDP, SSDP, SNMP, CHARGEN, OPENVPN, CLDAP, MEMCACHED"))
}

func printMenu() {
	fmt.Println("\n" + color.HiYellowString("=== 服务器控制菜单 ==="))
	fmt.Println("1. 启动所有服务器")
	fmt.Println("2. 停止所有服务器")
	fmt.Println("3. 查看服务器状态")
	fmt.Println("4. 单独控制服务器")
	fmt.Println("5. 刷新服务器状态")
	fmt.Println("6. 退出")
	fmt.Print("\n请输入选项: ")
}

func handleInput(ctx context.Context) {
	var choice int
	if _, err := fmt.Scanln(&choice); err != nil {
		fmt.Println(color.RedString("输入错误，请重新输入"))
		return
	}

	switch choice {
	case 1:
		startAllServers(ctx)
		fmt.Println(color.GreenString("所有服务器已启动"))
	case 2:
		stopAllServers()
		fmt.Println(color.RedString("所有服务器已停止"))
	case 3:
		showStatus()
	case 4:
		controlIndividual(ctx)
	case 5:
		resetServerStatus()
		fmt.Println(color.GreenString("服务器状态已刷新"))
	case 6:
		os.Exit(0)
	default:
		fmt.Println(color.RedString("无效选项"))
	}
}

func startAllServers(ctx context.Context) {
	for proto := range servers {
		if !servers[proto].IsRunning {
			go startServer(ctx, proto)
		}
	}
}

func stopAllServers() {
	for proto := range servers {
		servers[proto].IsRunning = false
	}
}

func showStatus() {
	statColors.Println("\n=== 服务器状态 ===")
	fmt.Println("协议名称\t状态\t\t收到包数\t最后活动时间")
	fmt.Println("--------------------------------------------------------")
	for proto, status := range servers {
		stateColor := color.RedString("停止")
		if status.IsRunning {
			stateColor = color.GreenString("运行中")
		}

		lastActivity := "N/A"
		if !status.LastActivity.IsZero() {
			lastActivity = status.LastActivity.Format("2006-01-02 15:04:05")
		}

		fmt.Printf("%-10s\t%s\t%d\t\t%s\n",
			color.HiCyanString(proto),
			stateColor,
			atomic.LoadUint64(&status.PacketsReceived),
			lastActivity)
	}
	fmt.Println()
}

func controlIndividual(ctx context.Context) {
	fmt.Println("\n" + color.HiYellowString("=== 单独控制服务器 ==="))
	fmt.Println("1. DNS")
	fmt.Println("2. NTP")
	fmt.Println("3. RDP")
	fmt.Println("4. SSDP")
	fmt.Println("5. SNMP")
	fmt.Println("6. CHARGEN")
	fmt.Println("7. OPENVPN")
	fmt.Println("8. CLDAP")
	fmt.Println("9. MEMCACHED")
	fmt.Println("0. 返回上级菜单")
	fmt.Print("\n请选择服务器: ")

	var choice int
	if _, err := fmt.Scanln(&choice); err != nil {
		fmt.Println(color.RedString("输入错误"))
		return
	}

	var proto string
	switch choice {
	case 1:
		proto = "DNS"
	case 2:
		proto = "NTP"
	case 3:
		proto = "RDP"
	case 4:
		proto = "SSDP"
	case 5:
		proto = "SNMP"
	case 6:
		proto = "CHARGEN"
	case 7:
		proto = "OPENVPN"
	case 8:
		proto = "CLDAP"
	case 9:
		proto = "MEMCACHED"
	case 0:
		return
	default:
		fmt.Println(color.RedString("无效选项"))
		return
	}

	fmt.Printf("\n当前 [%s] 服务器状态：", proto)
	if servers[proto].IsRunning {
		fmt.Println(color.GreenString("运行中"))
	} else {
		fmt.Println(color.RedString("停止"))
	}
	fmt.Println("1. 启动服务器")
	fmt.Println("2. 停止服务器")
	fmt.Println("3. 返回上级菜单")
	fmt.Print("\n请输入操作选项: ")

	var op int
	if _, err := fmt.Scanln(&op); err != nil {
		fmt.Println(color.RedString("输入错误"))
		return
	}
	switch op {
	case 1:
		if servers[proto].IsRunning {
			fmt.Println(color.YellowString("服务器已经在运行"))
		} else {
			go startServer(ctx, proto)
			fmt.Printf("[%s] 服务器启动中...\n", proto)
		}
	case 2:
		servers[proto].IsRunning = false
		fmt.Printf("[%s] 服务器已停止\n", proto)
	case 3:
		return
	default:
		fmt.Println(color.RedString("无效选项"))
	}
	time.Sleep(1 * time.Second)
}

func startServer(ctx context.Context, proto string) {
	status := servers[proto]
	status.IsRunning = true

	addr := &net.UDPAddr{Port: ports[proto]}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("启动 %s 服务器失败: %v\n", proto, err)
		status.IsRunning = false
		return
	}
	defer conn.Close()

	buf := make([]byte, 1500)
	for status.IsRunning {
		select {
		case <-ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, raddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				// 忽略超时错误
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				fmt.Printf("[%s] 接收数据错误: %v\n", proto, err)
				continue
			}

			atomic.AddUint64(&status.PacketsReceived, 1)
			status.LastActivity = time.Now()

			go handlePacket(proto, conn, raddr, buf[:n])
		}
	}
}

func handlePacket(proto string, conn *net.UDPConn, raddr *net.UDPAddr, data []byte) {
	switch proto {
	case "DNS":
		sender.SendDNSResponse(conn, raddr, data)
	case "NTP":
		sender.SendNTPResponse(conn, raddr)
	case "RDP":
		sender.SendRDPResponse(conn, raddr)
	case "SSDP":
		sender.SendSSDPResponse(conn, raddr)
	case "SNMP":
		sender.SendSNMPResponse(conn, raddr)
	case "CHARGEN":
		sender.SendCHARGENResponse(conn, raddr)
	case "OPENVPN":
		sender.SendOPENVPNResponse(conn, raddr)
	case "CLDAP":
		sender.SendCLDAPResponse(conn, raddr)
	case "MEMCACHED":
		sender.SendMEMCACHEDResponse(conn, raddr)
	}
}

func resetServerStatus() {
	for _, status := range servers {
		atomic.StoreUint64(&status.PacketsReceived, 0)
	}
}
