package protocol

import (
    "bufio"       // For reading files line by line
    "fmt"         // For formatted output and error handling
    "os"          // For file operations
    "os/exec"
    "time"
    "net"
    "strings"     // For string operations
    "github.com/google/gopacket/pcap" // For network interface management
)

// LoadIPList reads a list of IPs from a file
func LoadIPList(filePath string) ([]string, error) {
    var ipList []string

    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    // Read IP addresses line by line
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            ipList = append(ipList, line)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("failed to read file: %v", err)
    }

    return ipList, nil
}

// FindInterface finds the network interface and recommended source IP for a target IP
func FindInterface(targetIP string) (string, string, error) {
    cmd := exec.Command("ip", "route", "get", targetIP)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", "", fmt.Errorf("failed to execute command: %v", err)
    }
    
    // Parse output, example: 10.100.0.151 dev ix100 src 10.100.0.150 uid 0 
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
        return "", "", fmt.Errorf("interface not found")
    }
    
    // Validate source IP
    if src != "" && net.ParseIP(src) == nil {
        return dev, "", nil // Ignore invalid source IP
    }
    
    return dev, src, nil
}

// FindMAC returns the MAC address of the next hop for a target IP
func FindMAC(targetIP string) (string, error) {
    cmd := exec.Command("ping", "-c", "1", "-W", "1", targetIP)
    cmd.Run()

    // Wait briefly to ensure ARP cache is updated
    time.Sleep(500 * time.Millisecond)

    // Use ip route get to get the next hop IP
    cmd = exec.Command("ip", "route", "get", targetIP)
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to get route: %v", err)
    }

    // Parse output to get the next hop IP
    lines := strings.Split(string(output), "\n")
    if len(lines) == 0 {
        return "", fmt.Errorf("no route found for IP: %s", targetIP)
    }

    // Parse the next hop IP
    var nextHopIP string
    for _, line := range lines {
        if strings.Contains(line, "via") {
            // If "via" exists, the target IP is in an external network
            fields := strings.Fields(line)
            for i, field := range fields {
                if field == "via" && i+1 < len(fields) {
                    nextHopIP = fields[i+1]
                    break
                }
            }
            break
        } else if strings.Contains(line, "dev") {
            // If no "via", the target IP is in the local network
            nextHopIP = targetIP
            break
        }
    }

    if nextHopIP == "" {
        return "", fmt.Errorf("could not determine next hop IP for: %s", targetIP)
    }

    // Use arp command to get the MAC address
    cmd = exec.Command("arp", "-n", nextHopIP)
    output, err = cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to get ARP entry: %v", err)
    }

    // Parse ARP output to get the MAC address
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

// GetSrcMAC retrieves the source MAC address for sending to a target IP
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

// SendPacket sends a packet to the target IP using the network interface
func SendPacket(targetIP string, packet []byte) {
    // Get the network interface and open it
    interfaceName, _, err := FindInterface(targetIP)
    if err != nil {
        fmt.Printf("No available network interface: %v\n", err)
        return
    }
    handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
    if err != nil {
        fmt.Printf("Failed to open device: %v\n", err)
        return
    }
    defer handle.Close()

    if err := handle.WritePacketData(packet); err != nil {
        fmt.Printf("Failed to send forged packet: %v\n", err)
        return	
    }
}
