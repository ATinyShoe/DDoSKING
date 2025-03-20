package protocol

import (
    "bytes"
    "math/rand"
    "strconv"
    "strings"
)

func SNMPResponseBuffer() []byte {
    // Create a response buffer
    var buf bytes.Buffer
    
    // First create the inner content, then add the header and length
    var contentBuf bytes.Buffer
    
    // Version: 1 (SNMP v2c)
    contentBuf.Write([]byte{0x02, 0x01, 0x01})
    
    // Community string: "public"
    community := "public"
    contentBuf.Write([]byte{0x04, byte(len(community))})
    contentBuf.Write([]byte(community))
    
    // Create GetResponse PDU content
    var pduBuf bytes.Buffer
    
    // Request ID (random)
    reqID := rand.Uint32()
    pduBuf.Write([]byte{0x02, 0x04})
    pduBuf.Write([]byte{
        byte(reqID >> 24),
        byte(reqID >> 16),
        byte(reqID >> 8),
        byte(reqID),
    })
    
    // Error status: 0 (noError)
    pduBuf.Write([]byte{0x02, 0x01, 0x00})
    
    // Error index: 0
    pduBuf.Write([]byte{0x02, 0x01, 0x00})
    
    // Create variable bindings content
    var varbindsBuf bytes.Buffer
    
    // Use a more reasonable number of OIDs (15-25)
    numOIDs := 18
    
    // Some common SNMP OIDs
    commonOIDs := []string{
        "1.3.6.1.2.1.1.1.0",    // sysDescr
        "1.3.6.1.2.1.1.3.0",    // sysUpTime
        "1.3.6.1.2.1.1.4.0",    // sysContact
        "1.3.6.1.2.1.1.5.0",    // sysName
        "1.3.6.1.2.1.1.6.0",    // sysLocation
        "1.3.6.1.2.1.1.7.0",    // sysServices
        "1.3.6.1.2.1.2.1.0",    // ifNumber
        "1.3.6.1.2.1.2.2.1.1.1", // ifIndex.1
        "1.3.6.1.2.1.2.2.1.2.1", // ifDescr.1
        "1.3.6.1.2.1.2.2.1.3.1", // ifType.1
        "1.3.6.1.2.1.2.2.1.5.1", // ifSpeed.1
        "1.3.6.1.2.1.4.1.0",    // ipForwarding
        "1.3.6.1.2.1.4.20.1.1.192.168.1.1", // ipAdEntAddr
        "1.3.6.1.2.1.25.1.1.0",  // hrSystemUptime
        "1.3.6.1.2.1.25.2.2.0",  // hrMemorySize
    }
    
    // Add reasonable values for each OID
    for i := 0; i < numOIDs; i++ {
        var oidBuf bytes.Buffer
        
        // Use predefined OIDs or generate reasonable OIDs
        var oid string
        if i < len(commonOIDs) {
            oid = commonOIDs[i]
        } else {
            // Generate an extended enterprise OID (1.3.6.1.4.1....)
            oid = "1.3.6.1.4.1." + strconv.Itoa(9999) + "." + strconv.Itoa(i)
        }
        
        // Convert OID string to byte sequence
        oidBytes := encodeOID(oid)
        oidBuf.Write([]byte{0x06, byte(len(oidBytes))})
        oidBuf.Write(oidBytes)
        
        // Add reasonable values based on OID type
        var valueBytes []byte
        switch {
        case strings.HasPrefix(oid, "1.3.6.1.2.1.1.1.0"): // sysDescr
            value := "Hardware: Dell PowerEdge R640, Intel Xeon Gold 6230 CPU\nSoftware: Linux 5.10.0-15-amd64 #1 SMP Debian 5.10.120-1"
            valueBytes = []byte{0x04, byte(len(value))}
            valueBytes = append(valueBytes, []byte(value)...)
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.1.3.0"): // sysUpTime
            // Assume uptime is 10 days in Timeticks format (1/100 seconds)
            uptime := uint32(10 * 24 * 60 * 60 * 100)
            valueBytes = []byte{0x43, 0x04}
            valueBytes = append(valueBytes, []byte{
                byte(uptime >> 24),
                byte(uptime >> 16),
                byte(uptime >> 8),
                byte(uptime),
            }...)
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.1.4.0"): // sysContact
            value := "admin@example.com"
            valueBytes = []byte{0x04, byte(len(value))}
            valueBytes = append(valueBytes, []byte(value)...)
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.1.5.0"): // sysName
            value := "core-router-01"
            valueBytes = []byte{0x04, byte(len(value))}
            valueBytes = append(valueBytes, []byte(value)...)
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.1.6.0"): // sysLocation
            value := "Server Room 3, Rack 12, Unit 5-6"
            valueBytes = []byte{0x04, byte(len(value))}
            valueBytes = append(valueBytes, []byte(value)...)
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.1.7.0"): // sysServices
            valueBytes = []byte{0x02, 0x01, 0x7F} // 127 (typical router value)
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.2.1.0"): // ifNumber
            valueBytes = []byte{0x02, 0x01, 0x08} // 8 interfaces
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.2.2.1.1"): // ifIndex
            valueBytes = []byte{0x02, 0x01, byte(i % 8 + 1)}
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.2.2.1.2"): // ifDescr
            value := "GigabitEthernet0/" + strconv.Itoa(i % 8)
            valueBytes = []byte{0x04, byte(len(value))}
            valueBytes = append(valueBytes, []byte(value)...)
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.2.2.1.3"): // ifType
            valueBytes = []byte{0x02, 0x01, 0x06} // ethernetCsmacd(6)
            
        case strings.HasPrefix(oid, "1.3.6.1.2.1.2.2.1.5"): // ifSpeed
            // 1 Gbps = 1000000000
            speed := uint32(1000000000)
            valueBytes = []byte{0x42, 0x04}
            valueBytes = append(valueBytes, []byte{
                byte(speed >> 24),
                byte(speed >> 16),
                byte(speed >> 8),
                byte(speed),
            }...)
            
        default:
            // For other OIDs, generate a reasonable string value (20-30 bytes)
            valueLen := 20 + rand.Intn(10)
            value := make([]byte, valueLen)
            // Use readable characters
            for j := range value {
                // Generate letters and numbers
                if rand.Intn(2) == 0 {
                    value[j] = byte(rand.Intn(26) + 'a')
                } else {
                    value[j] = byte(rand.Intn(10) + '0')
                }
            }
            valueBytes = []byte{0x04, byte(len(value))}
            valueBytes = append(valueBytes, value...)
        }
        
        // Add OID and value to varbind
        var varbindBuf bytes.Buffer
        varbindBuf.Write(oidBuf.Bytes())
        varbindBuf.Write(valueBytes)
        
        // Wrap varbind in a sequence
        varbindBytes := varbindBuf.Bytes()
        varbindsBuf.Write([]byte{0x30, byte(len(varbindBytes))})
        varbindsBuf.Write(varbindBytes)
    }
    
    // Wrap varbinds as a varbinds sequence
    varbindsBytes := varbindsBuf.Bytes()
    pduBuf.Write([]byte{0x30, byte(len(varbindsBytes))})
    pduBuf.Write(varbindsBytes)
    
    // Wrap PDU content as GetResponse PDU
    pduBytes := pduBuf.Bytes()
    contentBuf.Write([]byte{0xA2, byte(len(pduBytes))})
    contentBuf.Write(pduBytes)
    
    // Finally, wrap the entire content as an SNMP message
    contentBytes := contentBuf.Bytes()
    buf.Write([]byte{0x30, byte(len(contentBytes))})
    buf.Write(contentBytes)
    
    return buf.Bytes()
}

// Helper function: Encode OID string to BER format
func encodeOID(oidStr string) []byte {
    parts := strings.Split(oidStr, ".")
    if len(parts) < 2 {
        return []byte{0}
    }
    
    // The first byte is 1*40 + 2
    var result []byte
    firstVal := 40*parseInt(parts[0]) + parseInt(parts[1])
    result = append(result, byte(firstVal))
    
    // Process the remaining parts
    for i := 2; i < len(parts); i++ {
        value := parseInt(parts[i])
        
        // Simple handling for values < 127
        if value < 128 {
            result = append(result, byte(value))
        } else {
            // Handle large values (actual implementation should support arbitrary-sized integers)
            octets := make([]byte, 0)
            octets = append(octets, byte(value&0x7F))
            value >>= 7
            
            for value > 0 {
                octets = append([]byte{byte(0x80 | (value & 0x7F))}, octets...)
                value >>= 7
            }
            
            result = append(result, octets...)
        }
    }
    
    return result
}

// Helper function: Convert string to integer
func parseInt(s string) int {
    val, _ := strconv.Atoi(s)
    return val
}

// SNMP protocol packet construction
func SNMPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
    payload := []byte{
        0x30, 0x37, 0x02, 0x01, 0x01, 0x04, 0x06, 0x70, 
        0x75, 0x62, 0x6c, 0x69, 0x63, 0xA5, 0x2A, 0x02,
        0x04, 0x71, 0x6F, 0xE4, 0x53, 0x02, 0x01, 0x00,
        0x02, 0x01, 0x14, 0x30, 0x1C, 0x30, 0x0B, 0x06,
        0x07, 0x2B, 0x06, 0x01, 0x02, 0x01, 0x01, 0x01,
        0x05, 0x00, 0x30, 0x0D, 0x06, 0x09, 0x2B, 0x06,
        0x01, 0x02, 0x01, 0x01, 0x09, 0x01, 0x03, 0x05,
        0x00,
    }
    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}
