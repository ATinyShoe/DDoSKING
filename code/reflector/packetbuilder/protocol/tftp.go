package protocol

func TFTPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
	// Construct a TFTP read request (attempting to fetch a large file)
	payload := []byte{
		0x00, 0x01, // OP=Read Request
		0x66, 0x69, 0x6C, 0x65, 0x6E, 0x61, 0x6D, 0x65, 0x00, // "filename"
		0x6F, 0x63, 0x74, 0x65, 0x74, 0x00, // "octet" mode
	}
	return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}
