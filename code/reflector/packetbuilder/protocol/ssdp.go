package protocol

import (
    "fmt"
    "strings"
    "time"
)

// SSDPResponseBuffer generates SSDP response data - returns multiple response packets for cyclic sending
// Returns multiple response packets instead of merging them into one to simulate attacks more realistically
func SSDPResponseBuffer() [][]byte {
    // Generate a long but valid UUID
    uuid := "uuid:" + strings.Repeat("f", 8) + "-" +
        strings.Repeat("e", 4) + "-" +
        strings.Repeat("d", 4) + "-" +
        strings.Repeat("a", 4) + "-" +
        strings.Repeat("b", 12)

    // More realistic device type list
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

    // Realistic device description URL paths
    paths := []string{
        "gatedesc.xml",
        "igd.xml",
        "upnp/IGD.xml",
        "tr64desc.xml",
        "RootDevice.xml",
        "device.xml",
        "DeviceDescription.xml",
    }

    // Realistic server identifiers
    serverTypes := []string{
        "UNIX/5.0 UPnP/1.0 Cisco/1.0",
        "Linux/3.10.39 UPnP/1.0 Technicolor/1.0",
        "Debian/4.0 UPnP/1.0 MiniUPnPd/1.8",
        "Windows NT/6.1 UPnP/1.0 MiniUPnPd/1.9",
        "Ubuntu/18.04 UPnP/1.1 MiniDLNA/1.2.1",
    }

    // Create a slice to store all response data
    responses := make([][]byte, len(serviceTypes))

    // Generate responses for each service type
    for i, st := range serviceTypes {
        // Create a random but valid LOCATION URL for each service type
        ipPart := fmt.Sprintf("192.168.%d.%d", (i*7)%256, (i*13)%256)
        port := 1024 + (i*1000)%9000

        pathIndex := i % len(paths)
        pathExt := paths[pathIndex]

        // Add random but valid query parameters to increase URL length
        queryParams := fmt.Sprintf("mac=%s&device=%s&timestamp=%d",
            strings.Replace(uuid, "uuid:", "", 1),
            strings.Repeat("router", i+1),
            time.Now().Unix()+int64(i*100))

        location := fmt.Sprintf("http://%s:%d/%s?%s",
            ipPart,
            port,
            pathExt,
            queryParams)

        // Select server identifier
        serverIndex := i % len(serverTypes)
        server := serverTypes[serverIndex]

        // Construct USN (Unique Service Name)
        usn := fmt.Sprintf("%s::%s", uuid, st)

        // Build the complete response
        response := fmt.Sprintf(
            "HTTP/1.1 200 OK\r\n"+
                "CACHE-CONTROL: max-age=1800\r\n"+
                "DATE: %s\r\n"+ // RFC1123 formatted date
                "EXT: \r\n"+
                "LOCATION: %s\r\n"+
                "OPT: \"http://schemas.upnp.org/upnp/1/0/\"; ns=01\r\n"+ // Valid optional header
                "01-NLS: %s\r\n"+ // Notification ID
                "SERVER: %s\r\n"+
                "ST: %s\r\n"+
                "USN: %s\r\n"+
                "BOOTID.UPNP.ORG: %d\r\n"+ // UPnP version 1.1 extension
                "CONFIGID.UPNP.ORG: %d\r\n"+ // UPnP version 1.1 extension
                "CONTENT-LENGTH: %d\r\n\r\n"+
                "%s", // Add some XML content to increase size
            time.Now().Format(time.RFC1123),
            location,
            strings.Replace(uuid, "uuid:", "", 1),
            server,
            st,
            usn,
            time.Now().Unix()%9999,
            (time.Now().Unix()/10)%999,
            1000, // Content length
            strings.Repeat("<device><serviceList><service></service></serviceList></device>", 50),
        )

        // Add the response to the slice
        responses[i] = []byte(response)
    }

    return responses
}

// SSDP packet construction
func SSDPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
    payload := []byte(
        "M-SEARCH * HTTP/1.1\r\n" +
            "HOST: 239.255.255.250:1900\r\n" +
            "MAN: \"ssdp:discover\"\r\n" +
            "MX: 2\r\n" +
            "ST: ssdp:all\r\n\r\n")
    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}
