package protocol

import (
    "encoding/binary"
    "math/rand"
    "net"
    "time"
    "strconv"
)

func NTPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
    payload := []byte{0x17, 0x00, 0x03, 0x2a, 0x00, 0x00, 0x00, 0x00}
    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}

func NTPResponseBuffer() []byte {
    // More reasonable number of entries - real servers usually have 10-30 active clients
    const maxEntries = 12

    // Keep the original structure definition
    const entrySize = 8
    const headerSize = 8
    const detailsSize = 72
    const fullEntrySize = entrySize + detailsSize

    // Calculate total response size
    const responseSize = headerSize + (maxEntries * fullEntrySize)

    // Allocate buffer
    buf := make([]byte, responseSize)

    // Set NTP mode 7 header
    buf[0] = 0x27                      // LI=0, VN=4, Mode=7 (private mode)
    buf[1] = 0x02                      // Implementation-specific code (monlist command)
    buf[2] = 0x00                      // Response code (0 indicates success)
    buf[3] = 0x00                      // Sequence number and authentication flags
    binary.BigEndian.PutUint16(buf[4:6], uint16(responseSize)) // Response size
    binary.BigEndian.PutUint16(buf[6:8], uint16(maxEntries))   // Number of entries

    // Get the current Unix timestamp
    nowUnix := time.Now().Unix()

    // Offset between NTP timestamp and Unix timestamp (seconds from 1900 to 1970)
    const ntpEpochOffset int64 = 2208988800

    // Initialize random number generator
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    // Create a more realistic client IP distribution
    // 1. Some local network clients
    // 2. A few public NTP servers/pools
    // 3. Some random internet clients
    clientIPs := []net.IP{
        // Local network clients
        net.ParseIP("192.168.1.100"),
        net.ParseIP("192.168.1.101"),
        net.ParseIP("192.168.1.105"),
        net.ParseIP("10.0.0.15"),
        net.ParseIP("10.0.0.23"),

        // Known NTP server/pool IPs (examples)
        net.ParseIP("162.159.200.1"),   // time.cloudflare.com
        net.ParseIP("216.239.35.4"),    // time.google.com

        // Random internet clients
        net.ParseIP("80.81.217." + strconv.Itoa(r.Intn(255))),
        net.ParseIP("203.114." + strconv.Itoa(r.Intn(255)) + "." + strconv.Itoa(r.Intn(255))),
        net.ParseIP("45." + strconv.Itoa(r.Intn(255)) + "." + strconv.Itoa(r.Intn(255)) + "." + strconv.Itoa(r.Intn(255))),
        net.ParseIP("138." + strconv.Itoa(r.Intn(255)) + "." + strconv.Itoa(r.Intn(255)) + "." + strconv.Itoa(r.Intn(255))),
        net.ParseIP("190." + strconv.Itoa(r.Intn(255)) + "." + strconv.Itoa(r.Intn(255)) + "." + strconv.Itoa(r.Intn(255))),
    }

    // Populate entries with more realistic data
    for i := 0; i < maxEntries; i++ {
        baseOffset := headerSize + (i * fullEntrySize)

        // 1. Basic entry (IP and port) - 8 bytes
        // Use our predefined IP list
        ip := clientIPs[i%len(clientIPs)]

        // Generate more realistic NTP client ports
        // NTP clients typically use high random ports
        var port uint16
        if i < 2 {
            // Some clients might be other NTP servers using the standard NTP port
            port = 123
        } else {
            // Most clients use high ephemeral ports
            port = uint16(r.Intn(10000) + 50000)
        }

        // Copy IP and port
        copy(buf[baseOffset:baseOffset+4], ip.To4())
        binary.BigEndian.PutUint16(buf[baseOffset+4:baseOffset+6], port)

        // Add counter and flag fields - make connection counts more realistic
        // Local clients typically have higher connection counts
        var connCount uint16
        if i < 5 {  // Local network clients
            connCount = uint16(r.Intn(500) + 500)  // 500-1000
        } else {
            connCount = uint16(r.Intn(100) + 1)    // 1-100
        }
        binary.BigEndian.PutUint16(buf[baseOffset+6:baseOffset+8], connCount)

        // 2. Details field - 72 bytes
        detailsOffset := baseOffset + entrySize

        // Generate more realistic timestamps - different distributions within the last 24 hours
        var lastAccessDelta int64

        if i < 5 {  // Local clients recently active
            lastAccessDelta = int64(r.Intn(3600))  // Within the last hour
        } else if i < 8 {  // Some clients active a few hours ago
            lastAccessDelta = int64(r.Intn(7200) + 3600)  // 1-3 hours ago
        } else {  // Some clients active a few days ago
            lastAccessDelta = int64(r.Intn(86400) + 7200)  // 3 hours to 1 day ago
        }

        // Fill last access time (NTP timestamp format, 8 bytes)
        lastAccessTime := nowUnix - lastAccessDelta + ntpEpochOffset
        binary.BigEndian.PutUint32(buf[detailsOffset:detailsOffset+4], uint32(lastAccessTime))
        binary.BigEndian.PutUint32(buf[detailsOffset+4:detailsOffset+8], uint32(r.Intn(100000)))

        // Fill first access time (NTP timestamp format, 8 bytes)
        // First access should be earlier than last access
        firstAccessDelta := lastAccessDelta + int64(r.Intn(86400*7))  // 1-7 days earlier
        firstAccessTime := nowUnix - firstAccessDelta + ntpEpochOffset
        binary.BigEndian.PutUint32(buf[detailsOffset+8:detailsOffset+12], uint32(firstAccessTime))
        binary.BigEndian.PutUint32(buf[detailsOffset+12:detailsOffset+16], uint32(r.Intn(100000)))

        // Fill the remaining details fields (including MRU index, packet counts, etc.)
        // Use more realistic packet counts, typically proportional to connection time
        for j := 16; j < detailsSize; j += 4 {
            var packetCount uint32
            if j == 16 { // Assume this is the packet count field
                // Local clients have more packets
                if i < 5 {
                    packetCount = uint32(r.Intn(5000) + 1000)
                } else {
                    packetCount = uint32(r.Intn(500) + 10)
                }
            } else {
                packetCount = uint32(r.Intn(100))
            }
            binary.BigEndian.PutUint32(buf[detailsOffset+j:detailsOffset+j+4], packetCount)
        }
    }
    return buf
}
