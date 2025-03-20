package protocol
import(
    "encoding/binary"
    "math/rand"
    "net"
    "time"
)

func NTPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
    payload := []byte{0x17, 0x00, 0x03, 0x2a, 0x00, 0x00, 0x00, 0x00}
    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}

func NTPResponseBuffer() []byte {
    // 标准NTP monlist响应最多包含600个条目
    const maxEntries = 600
   
    // 每个monlist条目在真实实现中实际上是8字节的IP+端口+计数器+其他字段
    // 完整的格式为：4字节IP + 2字节端口 + 2字节其他信息
    const entrySize = 8
    
    // 计算头部大小：标准NTP mode 7头部为8字节，扩展信息可能有额外字段
    const headerSize = 8
    
    // 再加上每个条目的详细信息字段(72字节)，这才是真实的monlist响应格式
    // 包括：最后访问时间、首次访问时间、MRU索引、packet count等
    const detailsSize = 72
    
    // 每个完整条目的总大小
    const fullEntrySize = entrySize + detailsSize
    
    // 计算总响应大小
    const responseSize = headerSize + (maxEntries * fullEntrySize)
    
    // 分配缓冲区
    buf := make([]byte, responseSize)
    
    // 设置NTP mode 7头部
    buf[0] = 0x27                      // LI=0, VN=4, Mode=7 (私有模式)
    buf[1] = 0x02                      // 实现特定代码 (monlist命令)
    buf[2] = 0x00                      // 响应码(0表示成功)
    buf[3] = 0x00                      // 序列号和认证标志
    binary.BigEndian.PutUint16(buf[4:6], uint16(responseSize)) // 响应大小
    binary.BigEndian.PutUint16(buf[6:8], uint16(maxEntries))   // 条目数量
    
    // 获取当前时间的Unix时间戳
    nowUnix := time.Now().Unix()
    
    // NTP时间戳与Unix时间戳的偏移量 (从1900年到1970年的秒数)
    const ntpEpochOffset int64 = 2208988800
    
    // 初始化随机数生成器
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    
    // 填充600个客户端记录
    for i := 0; i < maxEntries; i++ {
        baseOffset := headerSize + (i * fullEntrySize)
        
        // 1. 基本条目 (IP和端口) - 8字节
        // 生成随机但有效的IPv4地址，避开保留地址
        var ip net.IP
        for {
            first := byte(r.Intn(223) + 1) // 避开0和224-255
            
            // 避开常见的保留IP范围
            if (first == 10) || // 10.0.0.0/8
               (first == 127) || // 127.0.0.0/8
               (first == 169 && byte(r.Intn(256)) == 254) || // 169.254.0.0/16
               (first == 172 && byte(r.Intn(256)) >= 16 && byte(r.Intn(256)) <= 31) || // 172.16.0.0/12
               (first == 192 && byte(r.Intn(256)) == 168) { // 192.168.0.0/16
                continue
            }
            
            ip = net.IPv4(
                first,
                byte(r.Intn(256)),
                byte(r.Intn(256)),
                byte(r.Intn(256)),
            )
            break
        }
        
        // 生成随机但有效的端口
        port := uint16(r.Intn(65535-1024) + 1024)
        
        // 复制IP和端口
        copy(buf[baseOffset:baseOffset+4], ip.To4())
        binary.BigEndian.PutUint16(buf[baseOffset+4:baseOffset+6], port)
        
        // 添加计数器和标志字段，完成基本条目
        binary.BigEndian.PutUint16(buf[baseOffset+6:baseOffset+8], uint16(r.Intn(1000)))
        
        // 2. 详细信息字段 - 72字节
        detailsOffset := baseOffset + entrySize
        
        // 填充最后访问时间 (NTP时间戳格式，8字节)
        lastAccessTime := nowUnix + ntpEpochOffset
        binary.BigEndian.PutUint32(buf[detailsOffset:detailsOffset+4], uint32(lastAccessTime))
        binary.BigEndian.PutUint32(buf[detailsOffset+4:detailsOffset+8], uint32(r.Intn(1000000)))
        
        // 填充首次访问时间 (NTP时间戳格式，8字节)
        // 修复了int64和int类型不匹配的问题
        firstAccessTime := nowUnix - int64(r.Intn(86400)) + ntpEpochOffset
        binary.BigEndian.PutUint32(buf[detailsOffset+8:detailsOffset+12], uint32(firstAccessTime))
        binary.BigEndian.PutUint32(buf[detailsOffset+12:detailsOffset+16], uint32(r.Intn(1000000)))
        
        // 填充其余详细信息字段 (包括MRU索引、packet counts等)
        for j := 16; j < detailsSize; j += 4 {
            binary.BigEndian.PutUint32(buf[detailsOffset+j:detailsOffset+j+4], uint32(r.Intn(10000)))
        }
    }
    return buf
}