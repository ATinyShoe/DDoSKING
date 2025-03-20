package sender

import (
	"fmt"
	"net"
	"reflector/packetbuilder/protocol"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// SendDNSResponse constructs and sends an amplified DNS response
func SendDNSResponse(conn *net.UDPConn, raddr *net.UDPAddr, req []byte) {
	query, err := protocol.DNSParse(req)
	if err != nil {
		fmt.Println("Failed to parse DNS packet:", err)
		return
	}
	domain := query.Questions[0].Name

	response := protocol.DNSResponse()
	response.ID = query.ID
	response.RA = false
	response.Questions[0] = query.Questions[0]
	response.Questions[0].Name = []byte(domain)

	// To bypass domain name compression, make the NS records as long as possible and ensure each domain name is unique
	response.Answers = []layers.DNSResourceRecord{
		{
			Name:  []byte(domain),
			Type:  layers.DNSTypeNS,
			Class: layers.DNSClassIN,
			TTL:   36000,
			NS:    []byte("fd92yovrr5t3oaa9jahtdjauh4rtguv1da3ovzge399f2mvagbvxg4r000cjllr.v0fg52suitru09b8sarrxk7mz3u7kzkwnapfb0b5vfvhl6tf6xh0pff1nk87qhg.ojrjh5vnqm2uot8kti5r9w5cwzx1oeki6t9jfq2dlj6izlbfy4ha5na64do.ssonk0151omineqtnjdsvh6ubag4o53n3civj9q6mecxldunywpteuxtn8878.com"),
		},
		{
			Name:  []byte(domain),
			Type:  layers.DNSTypeNS,
			Class: layers.DNSClassIN,
			TTL:   36000,
			NS:    []byte("80yzwm6r7zjkg834ixo2lqmv5ddozfnniedcb3c3uwmjn03qtup0zb5ewpqwuhc.nijsi94mxd12esjonaq9c48dz3bj8svicu7kor5n9ls4ykp4igu1r3yf9lssb80.3vnw5xqxtvvc06zw5agg2fn0cz7z1p6my45vgb4nnxekab9pm4jff8vf6ep.973z6kw2ytcg94xzxoszxeiy5pub3cakv5y3bcfno894l33e3hdqfkapkbb0w.com"),
		},
		{
			Name:  []byte(domain),
			Type:  layers.DNSTypeNS,
			Class: layers.DNSClassIN,
			TTL:   36000,
			NS:    []byte("1nquklci8g555cpze82c7e02uqfkyxoy24yi18cpja88kqy33g5smqg3yjecz0r.bm85yd0c063zk7hsd9yykykxhj4p6lub8tmqob8zshhac8by6sn7puj9ya2i7ci.epr5viravte0dtqqus7l46djvfijk11ffc947s70f98u456trimyi882bqo.1niplcvnt6agrvbyyi6qij8bpdaishs4zmeqx0pceiwetuoxwlurso4l4n9ln.com"),
		},
		{
			Name:  []byte(domain),
			Type:  layers.DNSTypeNS,
			Class: layers.DNSClassIN,
			TTL:   36000,
			NS:    []byte("2rafjtv0y398mrttsx6ezohz4ult8c336eoi43k1suwmcbn9gg26sfy449v75.com"),
		},
	}

	dnsBuffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	if err := response.SerializeTo(dnsBuffer, options); err != nil {
		fmt.Printf("Failed to serialize DNS layer: %v\n", err)
		return
	}

	conn.WriteToUDP(dnsBuffer.Bytes(), raddr)
}

// SendNTPResponse constructs and sends an amplified NTP response
func SendNTPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.NTPResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendRDPResponse constructs and sends an amplified RDP response
func SendRDPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.RDPResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendSSDPResponse constructs and sends an amplified SSDP response
// Typical amplification factor: 1:30
func SendSSDPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	// Get multiple response packets
	responses := protocol.SSDPResponseBuffer()

	// Loop through and send all response packets to achieve higher amplification
	for _, resp := range responses {
		conn.WriteToUDP(resp, raddr)
	}
}

// SendSNMPResponse constructs and sends an amplified SNMP response
// Amplification factor: approximately 6-10x
func SendSNMPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.SNMPResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendCHARGENResponse constructs and sends an amplified CHARGEN response
// Amplification factor: up to 358x
func SendCHARGENResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.CHARGENResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendOPENVPNResponse constructs and sends an amplified OpenVPN response
// Amplification factor: approximately 2-3x
func SendOPENVPNResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.OPENVPNResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendCLDAPResponse constructs and sends an amplified CLDAP response
// Amplification factor: approximately 56-70x
func SendCLDAPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	// Get multiple response packets
	responses := protocol.CLDAPResponseBuffer()

	// Loop through and send all response packets to achieve higher amplification
	for _, resp := range responses {
		conn.WriteToUDP(resp, raddr)
	}
}

// SendMEMCACHEDResponse constructs and sends an amplified Memcached response
// Amplification factor: up to 51,000x
func SendMEMCACHEDResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	// Get multiple response packets
	responses := protocol.MEMCACHEDResponseBuffer()

	// Loop through and send all response packets to achieve extremely high amplification
	for _, resp := range responses {
		conn.WriteToUDP(resp, raddr)
	}
}
