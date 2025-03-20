package protocol

import (
    "encoding/binary"
    "bytes"
)

// RDPPacket constructs an RDP request packet
func RDPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
    payload := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}

func RDPResponseBuffer() []byte {
    // A more realistic RDP initial response size is around 1200-1500 bytes
    buf := new(bytes.Buffer)
    
    // Create content first, then backfill the TPKT header length
    // Reserve space for the TPKT header
    buf.Write([]byte{0x03, 0x00, 0x00, 0x00})
    
    // X.224 Connection Confirm (RFC 2126)
    buf.Write([]byte{
        0x0D,       // Length field (TPDU code length)
        0xD0,       // CC-TPDU type (Connection Confirm)
        0x00, 0x00, // DST-REF Destination reference
        0x12, 0x34, // SRC-REF Source reference (random value)
        0x00,       // Class options
    })
    
    // RDP Negotiation Response (MS-RDPBCGR 2.2.1.2)
    buf.Write([]byte{
        0x01,                   // TYPE - TYPE_RDP_NEG_RSP
        0x00,                   // Flags - Reserved
        0x08, 0x00,             // Length field (little-endian)
        0x03, 0x00, 0x00, 0x00, // Supported protocols: TLS 1.0 + CredSSP
    })
    
    // MCS Connect Response
    mcsHeader := []byte{
        0x02, 0xF0, 0x80,       // TPKT header + BER: CONNECT-RESPONSE
        0x7F, 0x66,             // BER: APPLICATION 101, length
    }
    buf.Write(mcsHeader)
    
    // Result = rt-successful
    buf.Write([]byte{0x0A, 0x01, 0x00})
    
    // Connect ID
    buf.Write([]byte{0x02, 0x01, 0x22})
    
    // Domain parameters - BER encoded
    buf.Write([]byte{
        0x04, 0x82, 0x01, 0x4B, // OCTET STRING, length = 331 bytes
        // Below is the actual structure of domain parameters, but we use random data for simplicity
    })
    for i := 0; i < 320; i++ {
        buf.WriteByte(byte(i % 256))
    }
    
    // GCC Conference Create Response
    gccHeader := []byte{
        0x00, 0x05, 0x00, 0x14, // PER encoded header
        0x7C, 0x00, 0x01,       // Conference create response
        0x2A, 0x14, 0x76, 0x0A, // User data length
    }
    buf.Write(gccHeader)
    
    // Server Core Data (MS-RDPBCGR 2.2.1.4.2)
    buf.Write([]byte{
        0x01, 0x0C,             // TS_UD_HEADER.type = SC_CORE
        0x6C, 0x00,             // TS_UD_HEADER.length = 108 (actual length)
        0x04, 0x00, 0x08, 0x00, // Version and flags
        0x01, 0x00,             // Connection type = RDP
        0x20, 0x00,             // 32 bits per pixel
        0x03, 0x00,             // Color depth = 24 bpp
        0x60, 0x00,             // SASSequence
        0x00, 0x04, 0x00, 0x00, // Keyboard layout = US
        0x2C, 0x01, 0x00, 0x00, // Client build = 300
    })
    
    // Add server name (fake a reasonable name)
    serverName := "RD-SERVER01"
    nameBytes := []byte(serverName)
    for i := 0; i < 32; i++ {
        if i < len(nameBytes) {
            buf.WriteByte(nameBytes[i])
        } else {
            buf.WriteByte(0x00) // Fill the remaining part with 0
        }
    }
    
    // Add more realistic fields
    buf.Write([]byte{
        0x30, 0x00, 0x00, 0x00, // Key length
        0x01, 0x00, 0x00, 0x00, // IO channel count
        0x01, 0x00, 0x00, 0x00, // High color depth
        0x01, 0x00,             // Flag
        0x01, 0x00,             // Compression
        0x00, 0x00,             // Padding
    })
    
    // Server Security Data (MS-RDPBCGR 2.2.1.4.3)
    buf.Write([]byte{
        0x02, 0x0C,             // TS_UD_HEADER.type = SC_SECURITY
        0x2C, 0x00,             // TS_UD_HEADER.length = 44
        0x03, 0x00, 0x00, 0x00, // Encryption methods
        0x00, 0x00, 0x00, 0x00, // Encryption level
        0x02, 0x00, 0x00, 0x00, // Server Random length
        0x00, 0x00, 0x00, 0x00, // Server certificate length
    })
    
    // Add random numbers (some RDP versions have this)
    for i := 0; i < 32; i++ {
        buf.WriteByte(byte(i * 3 % 256))
    }
    
    // Actual capability set section - Server Capability Set (MS-RDPBCGR 2.2.7)
    capSetHeader := []byte{
        0x0A, 0x0C,             // TS_UD_HEADER.type = SC_CAPS
        0x58, 0x02,             // TS_UD_HEADER.length = 600
    }
    buf.Write(capSetHeader)
    
    // Add 8 more realistic capability sets
    // 1. General Capability Set
    buf.Write([]byte{
        0x01, 0x00,             // Capability set type
        0x14, 0x00,             // Length: 20 bytes
        0x01, 0x00,             // Protocol version
        0x03, 0x00,             // General compression flags
        0x00, 0x00, 0x00, 0x00, // Extra flags
        0x00, 0x00, 0x00, 0x00, // Compression type
        0x01, 0x00,             // Update capability
        0x00, 0x00, 0x00, 0x00, // Remote frame buffer size
    })
    
    // 2. Bitmap Capability Set
    buf.Write([]byte{
        0x02, 0x00,             // Capability set type
        0x1C, 0x00,             // Length: 28 bytes
        0x20, 0x00, 0x00, 0x00, // Color depth
        0x01, 0x00, 0x00, 0x00, // Support flags
        0x01, 0x00,             // bpp
        0x00, 0x07,             // Width resolution
        0x00, 0x05,             // Height resolution
        0x00, 0x01,             // Desktop width
        0x00, 0x00,             // Desktop height
        0x00, 0x01,             // Desktop size
        0x01, 0x00,             // Palette flag
        0x00, 0x01,             // Palette cache size
    })
    
    // 3. Order Capability Set
    buf.Write([]byte{
        0x03, 0x00,             // Capability set type
        0x18, 0x00,             // Length: 24 bytes
    })
    for i := 0; i < 20; i++ {
        buf.WriteByte(byte(i % 256))
    }
    
    // 4. Audio Capability Set 
    buf.Write([]byte{
        0x0C, 0x00,             // Capability set type
        0x0C, 0x00,             // Length: 12 bytes
        0x01, 0x00, 0x00, 0x00, // Sound flags
        0x40, 0x6F, 0x00, 0x00, // Frequency (28,480 Hz)
    })
    
    // 5-8. Other four capability sets
    // For simplicity, only add reasonable-sized padding
    for i := 0; i < 4; i++ {
        buf.Write([]byte{
            byte(4 + i), 0x00,          // Capability set type
            0x18, 0x00,                 // Length: 24 bytes
        })
        for j := 0; j < 20; j++ {
            buf.WriteByte(byte((i*j) % 256))
        }
    }
    
    // After completion, backfill the TPKT header length
    completed := buf.Bytes()
    binary.BigEndian.PutUint16(completed[2:4], uint16(len(completed)))
    
    return completed
}
