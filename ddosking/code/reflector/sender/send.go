package sender

import (
	"fmt"
	"net"
	"reflector/packetbuilder/protocol"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// SendDNSResponse 构造并发送放大的DNS响应
func SendDNSResponse(conn *net.UDPConn, raddr *net.UDPAddr, req []byte) {
	query, err := protocol.DNSParse(req)
	if err != nil {
		fmt.Println("解析DNS报文失败:", err)
		return
	}
	domain := query.Questions[0].Name

	response := protocol.DNSResponse()
	response.ID = query.ID
	response.RA = false
	response.Questions[0] = query.Questions[0]
	response.Questions[0].Name = []byte(domain)

	// 为了绕过域名压缩技术，尽可能使得NS记录长且每个域名完全不同
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
		fmt.Printf("序列化 DNS 层失败: %v\n", err)
		return
	}

	conn.WriteToUDP(dnsBuffer.Bytes(), raddr)
}

// SendNTPResponse 构造并发送放大的NTP响应
func SendNTPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.NTPResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendRDPResponse 构造并发送放大的RDP响应
func SendRDPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.RDPResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendSSDPResponse 构造并发送放大的SSDP响应
// 典型放大倍数：1:30
func SendSSDPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	// 获取多个响应包
	responses := protocol.SSDPResponseBuffer()

	// 循环发送所有响应包以实现更高的放大倍数
	for _, resp := range responses {
		conn.WriteToUDP(resp, raddr)
	}
}

// SendSNMPResponse 构造并发送放大的SNMP响应
// 放大倍数：约 6-10 倍
func SendSNMPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.SNMPResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendCHARGENResponse 构造并发送放大的CHARGEN响应
// 放大倍数：高达 358 倍
func SendCHARGENResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.CHARGENResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendOPENVPNResponse 构造并发送放大的OpenVPN响应
// 放大倍数：约 2-3 倍
func SendOPENVPNResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	buf := protocol.OPENVPNResponseBuffer()
	conn.WriteToUDP(buf, raddr)
}

// SendCLDAPResponse 构造并发送放大的CLDAP响应
// 放大倍数：约 56-70 倍
func SendCLDAPResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	// 获取多个响应包
	responses := protocol.CLDAPResponseBuffer()

	// 循环发送所有响应包以实现更高的放大倍数
	for _, resp := range responses {
		conn.WriteToUDP(resp, raddr)
	}
}

// SendMEMCACHEDResponse 构造并发送放大的Memcached响应
// 放大倍数：高达 51,000 倍
func SendMEMCACHEDResponse(conn *net.UDPConn, raddr *net.UDPAddr) {
	// 获取多个响应包
	responses := protocol.MEMCACHEDResponseBuffer()

	// 循环发送所有响应包以实现极高的放大倍数
	for _, resp := range responses {
		conn.WriteToUDP(resp, raddr)
	}
}
