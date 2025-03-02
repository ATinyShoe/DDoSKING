package protocol

func TFTPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
	// 构造tftp读请求（尝试获取大文件）
	payload := []byte{
		0x00, 0x01, // OP=Read Request
		0x66, 0x69, 0x6C, 0x65, 0x6E, 0x61, 0x6D, 0x65, 0x00, // "filename"
		0x6F, 0x63, 0x74, 0x65, 0x74, 0x00, // "octet" mode
	}
	return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}
