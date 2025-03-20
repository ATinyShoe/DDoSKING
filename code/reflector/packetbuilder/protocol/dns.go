package protocol

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Set random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// DNS configuration function
func DNSQuery() layers.DNS {
	return layers.DNS{
		ID:     uint16(rand.Intn(65536)), // Random ID
		QR:     false,                   // Query request
		OpCode: layers.DNSOpCodeQuery,
		RD:     true,                    // Recursive query
		Questions: []layers.DNSQuestion{
			{
				Name:  []byte("example.com"),
				Type:  layers.DNSTypeNS,
				Class: layers.DNSClassIN,
			},
		},
		Additionals: []layers.DNSResourceRecord{
			{
				Type:  layers.DNSTypeOPT, // OPT record, enabling EDNS0
				Class: 4096,              // Maximum UDP response size supported by the client
			},
		},
	}
}

func DNSResponse() layers.DNS {
	return layers.DNS{
		ID:           uint16(rand.Intn(65536)), // Random ID, should match the request
		QR:           true,                    // Indicates this is a response packet
		OpCode:       layers.DNSOpCodeQuery,   // Query opcode
		AA:           false,                   // Authoritative answer flag
		TC:           false,                   // Not truncated
		RD:           false,                   // Recursive query
		RA:           true,                    // Server supports recursive query
		ResponseCode: layers.DNSResponseCodeNoErr, // No error

		Questions: []layers.DNSQuestion{
			{
				Name:  []byte("example.com"),
				Type:  layers.DNSTypeNS,
				Class: layers.DNSClassIN,
			},
		},

		Answers: []layers.DNSResourceRecord{
			{
				Name:  []byte("example.com"),
				Type:  layers.DNSTypeNS,
				Class: layers.DNSClassIN,
				TTL:   300,
				IP:    net.ParseIP("127.0.0.1"),
			},
		},

		Additionals: []layers.DNSResourceRecord{
			{
				Name:  []byte("ns1.example.com"),
				Type:  layers.DNSTypeA,
				Class: layers.DNSClassIN,      // IN class
				TTL:   36000,                 // TTL in seconds
				IP:    net.ParseIP("192.168.5.107"), // Response IP address
			},
		},
	}
}

func DNSPacket(srcIP, dstIP string, srcPort, dstPort int, dnsPayload layers.DNS) ([]byte, error) {
	// Serialize the DNS layer and generate a byte slice
	dnsBuffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	if err := dnsPayload.SerializeTo(dnsBuffer, options); err != nil {
		return nil, fmt.Errorf("Failed to serialize DNS layer: %v", err)
	}

	return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, dnsBuffer.Bytes())
}

func DNSParse(payload []byte) (dnsPacket layers.DNS, err error) {
	// Create a DNS layer parser
	packet := gopacket.NewPacket(payload, layers.LayerTypeDNS, gopacket.Default)

	// Get the DNS layer
	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		// Convert the parsed result to DNS type
		dnsPacket, _ := dnsLayer.(*layers.DNS)
		return *dnsPacket, nil
	}

	// If no DNS layer is parsed, return an error
	return layers.DNS{}, fmt.Errorf("Failed to parse as a valid DNS packet")
}
