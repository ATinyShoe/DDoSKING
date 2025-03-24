package protocol

func ARDPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
    payload := []byte{0x00, 0x14, 0x00, 0x00}
    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}