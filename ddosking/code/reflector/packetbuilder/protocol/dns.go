package protocol

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// 设置随机种子
func init() {
    rand.Seed(time.Now().UnixNano())
}

// DNS配置函数
func DNSQuery() layers.DNS {
	return layers.DNS{
		ID:     uint16(rand.Intn(65536)), // 随机 ID
		QR:     false,                   // 查询请求
		OpCode: layers.DNSOpCodeQuery,
		RD:     true,                    // 递归查询
		Questions: []layers.DNSQuestion{
			{
				Name:  []byte("example.com"), 
				Type:  layers.DNSTypeNS,      
				Class: layers.DNSClassIN,
			},
		},
		Additionals: []layers.DNSResourceRecord{
			{
				Type:  layers.DNSTypeOPT, // OPT 记录，启用 EDNS0
				Class: 4096,              // 客户端支持的最大 UDP 响应大小
			},
		},
	}
}

func DNSResponse() layers.DNS{
	return layers.DNS{
		ID:           uint16(rand.Intn(65536)), // 随机 ID，与请求匹配时应保持一致
		QR:           true,                    // 表示这是响应报文
		OpCode:       layers.DNSOpCodeQuery,   // 查询操作码
		AA:           false,                    // 权威答案标志
		TC:           false,                   // 没有被截断
		RD:           false,                    // 递归查询
		RA:           true,                    // 服务器支持递归查询
		ResponseCode: layers.DNSResponseCodeNoErr, // 无错误
		
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
				Class: layers.DNSClassIN,      // IN类   
				TTL:   36000,                    // TTL，单位秒
				IP:    net.ParseIP("192.168.5.107"), // 响应的 IP 地址
			},
		},


		// 	{
		// 		Name:  []byte("ns2.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.108"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("ns3.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.109"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("ns4.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.110"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("ns5.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.107"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("ns6.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.108"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("ns7.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.109"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("ns8.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.110"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("ns9.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.110"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("dns1.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.107"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("dns2.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.108"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("dns3.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.109"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("dns4.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.110"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("dns5.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.107"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("dns6.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.108"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("dns7.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.109"), // 响应的 IP 地址
		// 	},
		// 	{
		// 		Name:  []byte("dns8.example.com"),         
		// 		Type:  layers.DNSTypeA,     
		// 		Class: layers.DNSClassIN,      // IN类   
		// 		TTL:   36000,                    // TTL，单位秒
		// 		IP:    net.ParseIP("192.168.5.110"), // 响应的 IP 地址
		// 	},
		// },
	}
}

func DNSPacket(srcIP, dstIP string, srcPort, dstPort int,dnsPayload layers.DNS) ([]byte, error) {
    // 序列化 DNS 层并生成字节切片
    dnsBuffer := gopacket.NewSerializeBuffer()
    options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
    if err := dnsPayload.SerializeTo(dnsBuffer, options); err != nil {
        return nil, fmt.Errorf("序列化 DNS 层失败: %v", err)
    }

    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, dnsBuffer.Bytes())
}

func DNSParse(payload []byte) (dnsPacket layers.DNS, err error) {
	// 创建一个 DNS 层解析器
	packet := gopacket.NewPacket(payload, layers.LayerTypeDNS, gopacket.Default)

	// 获取 DNS 层
	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		// 将解析结果转换为 DNS 类型
		dnsPacket, _ := dnsLayer.(*layers.DNS)
		return *dnsPacket, nil
	}

	// 如果没有解析出 DNS 层，返回错误
	return layers.DNS{}, fmt.Errorf("无法解析为有效的 DNS 报文")
}