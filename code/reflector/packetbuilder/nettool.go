package packetbuilder

import (
    "bufio"       // 用于逐行读取文件
    "fmt"         // 格式化输出和错误处理
    "os"          // 文件操作
    "os/exec"
    "time"
    "net"
    "strings"     // 字符串操作
    "github.com/google/gopacket/pcap" // 用于网络接口管理
)

// 读取文件列表
func LoadIPList(filePath string) ([]string, error) {
    var ampList []string

    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("无法打开文件: %v", err)
    }
    defer file.Close()

    // 逐行读取文件中的 IP 地址
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            ampList = append(ampList, line)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("读取文件失败: %v", err)
    }

    return ampList, nil
}

// FindInterface 查找发往目标IP的网络接口和推荐源IP地址
func FindInterface(targetIP string) (string, string, error) {
    cmd := exec.Command("ip", "route", "get", targetIP)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", "", fmt.Errorf("执行命令失败: %v", err)
    }
    
    // 解析输出，示例输出：10.100.0.151 dev ix100 src 10.100.0.150 uid 0 
    parts := strings.Fields(string(output))
    var dev, src string
    for i := 0; i < len(parts); i++ {
        switch parts[i] {
        case "dev":
            if i+1 < len(parts) {
                dev = parts[i+1]
            }
        case "src":
            if i+1 < len(parts) {
                src = parts[i+1]
            }
        }
    }
    
    if dev == "" {
        return "", "", fmt.Errorf("未找到接口")
    }
    
    // 验证源IP是否有效
    if src != "" && net.ParseIP(src) == nil {
        return dev, "", nil // 忽略无效源IP
    }
    
    return dev, src, nil
}


// FindMAC 根据目标IP返回下一跳的MAC地址字符串
func FindMAC(targetIP string) (string, error) {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", targetIP)
    cmd.Run()

    // 等待一小段时间确保ARP缓存已经更新
	time.Sleep(500 * time.Millisecond)

	// 使用 ip route get 命令获取下一跳IP
	cmd = exec.Command("ip", "route", "get", targetIP)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get route: %v", err)
	}

	// 解析输出，获取下一跳IP
	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("no route found for IP: %s", targetIP)
	}

	// 解析下一跳IP
	var nextHopIP string
	for _, line := range lines {
		if strings.Contains(line, "via") {
			// 如果有 via，说明目标IP在外部网络
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "via" && i+1 < len(fields) {
					nextHopIP = fields[i+1]
					break
				}
			}
			break
		} else if strings.Contains(line, "dev") {
			// 如果没有 via，说明目标IP在局域网内
			nextHopIP = targetIP
			break
		}
	}

	if nextHopIP == "" {
		return "", fmt.Errorf("could not determine next hop IP for: %s", targetIP)
	}

	// 使用 arp 命令获取MAC地址
	cmd = exec.Command("arp", "-n", nextHopIP)
	output, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get ARP entry: %v", err)
	}

	// 解析ARP输出，获取MAC地址
	lines = strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("no ARP entry found for IP: %s", nextHopIP)
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 3 {
		return "", fmt.Errorf("invalid ARP entry format for IP: %s", nextHopIP)
	}

	macAddress := fields[2]
	return macAddress, nil
}


// GetSourceMAC 获取发送到目标IP的源MAC地址
func GetSrcMAC(targetIP string) (net.HardwareAddr, error) {
    ifaceName, _, err := FindInterface(targetIP)
    if err != nil {
        return nil, fmt.Errorf("failed to get route interface: %w", err)
    }

    iface, err := net.InterfaceByName(ifaceName)
    if err != nil {
        return nil, fmt.Errorf("failed to get interface %s: %w", ifaceName, err)
    }

    if iface.HardwareAddr == nil || len(iface.HardwareAddr) == 0 {
        return nil, fmt.Errorf("interface %s has no MAC address", ifaceName)
    }

    return iface.HardwareAddr, nil
}

// 调用网络接口发送报文
func SendPacket(targetIP string, packet []byte) {
    // 获取网络接口并打开
    interfaceName, _, err := FindInterface(targetIP)
    if err != nil {
        fmt.Printf("没有可用的网络接口: %v\n", err)
        return
    }
    handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
    if err != nil {
        fmt.Printf("无法打开设备: %v\n", err)
        return
    }
    defer handle.Close()

    if err := handle.WritePacketData(packet); err != nil {
        fmt.Printf("发送伪造数据包失败: %v\n", err)
        return	
    }
}


