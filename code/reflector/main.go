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
	fmt.Println(color.HiYellowString("DDoS Reflection Amplification Attack Simulator"))
	fmt.Println(color.HiWhiteString("Supported Protocols: DNS, NTP, RDP, SSDP, SNMP, CHARGEN, OPENVPN, CLDAP, MEMCACHED"))
}

func printMenu() {
	fmt.Println("\n" + color.HiYellowString("=== Server Control Menu ==="))
	fmt.Println("1. Start all servers")
	fmt.Println("2. Stop all servers")
	fmt.Println("3. View server status")
	fmt.Println("4. Control individual server")
	fmt.Println("5. Refresh server status")
	fmt.Println("6. Exit")
	fmt.Print("\nEnter your choice: ")
}

func handleInput(ctx context.Context) {
	var choice int
	if _, err := fmt.Scanln(&choice); err != nil {
		fmt.Println(color.RedString("Invalid input, please try again"))
		return
	}

	switch choice {
	case 1:
		startAllServers(ctx)
		fmt.Println(color.GreenString("All servers started"))
	case 2:
		stopAllServers()
		fmt.Println(color.RedString("All servers stopped"))
	case 3:
		showStatus()
	case 4:
		controlIndividual(ctx)
	case 5:
		resetServerStatus()
		fmt.Println(color.GreenString("Server status refreshed"))
	case 6:
		os.Exit(0)
	default:
		fmt.Println(color.RedString("Invalid option"))
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
	statColors.Println("\n=== Server Status ===")
	fmt.Println("Protocol\tStatus\t\tPackets Received\tLast Activity Time")
	fmt.Println("--------------------------------------------------------")
	for proto, status := range servers {
		stateColor := color.RedString("Stopped")
		if status.IsRunning {
			stateColor = color.GreenString("Running")
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
	fmt.Println("\n" + color.HiYellowString("=== Individual Server Control ==="))
	fmt.Println("1. DNS")
	fmt.Println("2. NTP")
	fmt.Println("3. RDP")
	fmt.Println("4. SSDP")
	fmt.Println("5. SNMP")
	fmt.Println("6. CHARGEN")
	fmt.Println("7. OPENVPN")
	fmt.Println("8. CLDAP")
	fmt.Println("9. MEMCACHED")
	fmt.Println("0. Return to previous menu")
	fmt.Print("\nSelect a server: ")

	var choice int
	if _, err := fmt.Scanln(&choice); err != nil {
		fmt.Println(color.RedString("Invalid input"))
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
		fmt.Println(color.RedString("Invalid option"))
		return
	}

	fmt.Printf("\nCurrent [%s] server status: ", proto)
	if servers[proto].IsRunning {
		fmt.Println(color.GreenString("Running"))
	} else {
		fmt.Println(color.RedString("Stopped"))
	}
	fmt.Println("1. Start server")
	fmt.Println("2. Stop server")
	fmt.Println("3. Return to previous menu")
	fmt.Print("\nEnter your choice: ")

	var op int
	if _, err := fmt.Scanln(&op); err != nil {
		fmt.Println(color.RedString("Invalid input"))
		return
	}
	switch op {
	case 1:
		if servers[proto].IsRunning {
			fmt.Println(color.YellowString("Server is already running"))
		} else {
			go startServer(ctx, proto)
			fmt.Printf("[%s] Server starting...\n", proto)
		}
	case 2:
		servers[proto].IsRunning = false
		fmt.Printf("[%s] Server stopped\n", proto)
	case 3:
		return
	default:
		fmt.Println(color.RedString("Invalid option"))
	}
	time.Sleep(1 * time.Second)
}

func startServer(ctx context.Context, proto string) {
	status := servers[proto]
	status.IsRunning = true

	addr := &net.UDPAddr{Port: ports[proto]}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Failed to start %s server: %v\n", proto, err)
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
				// Ignore timeout errors
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				fmt.Printf("[%s] Error receiving data: %v\n", proto, err)
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
