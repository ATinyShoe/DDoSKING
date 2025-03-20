package server

import (
	"auth/packetbuilder"
	"auth/packetbuilder/protocol"
	"fmt"
	"time"

	"github.com/google/gopacket/layers"
)

type Message struct {
	Payload []byte // 载荷
	SrcPort int    // 源端口号
}

// 睡眠响应时间（期间收到的包直接丢弃），停止睡眠后再收到包则直接发响应包。然后进入新一轮
func (s *Server) queryAggregation() {
	setTimeOut := 15 // 超时设置
	messageChan := make(chan Message)
	domain := "example.com"
	state := "waiting" // 状态: waiting, sleeping, ready

	go ListenUDP(s.DstIP, 53, messageChan)

	for {
		select {
		case msg := <-messageChan:
			switch state {
			case "waiting":
				// 收到第一个包，开始睡眠
				fmt.Printf("收到报文，开始超时时间 %v 秒的阻塞\n", setTimeOut)
				state = "sleeping"

				go func() {
					// 睡眠指定时间
					time.Sleep(time.Duration(setTimeOut) * time.Second)
					// 睡眠结束后，将状态改为准备接收
					state = "ready"
					fmt.Println("阻塞结束，等待新的DNS查询")
				}()

			case "sleeping":
				// 在睡眠期间收到的包直接丢弃
				fmt.Println("阻塞期间收到报文，忽略")

			case "ready":
				// 睡眠结束后收到的第一个包，立即发送响应
				query, err := protocol.DNSParse(msg.Payload)
				if err != nil {
					fmt.Println("解析DNS报文失败,", err)
					state = "waiting" // 重置状态
					continue
				}

				// 构造响应
				response := protocol.DNSResponse()
				response.ID = query.ID
				response.RA = false
				response.Questions[0] = query.Questions[0]
				response.Questions[0].Name = []byte(domain)
				// 为了绕过域名压缩技术，尽可能使得NS记录长且每个域名完全不同
				response.Answers = []layers.DNSResourceRecord{
					{
						Name:  []byte(domain),
						Type:  layers.DNSTypeNS,
						Class: layers.DNSClassIN,                                                                                                                                                                                                                                                       // IN类
						TTL:   36000,                                                                                                                                                                                                                                                                   // TTL，单位秒
						NS:    []byte("fd92yovrr5t3oaa9jahtdjauh4rtguv1da3ovzge399f2mvagbvxg4r000cjllr.v0fg52suitru09b8sarrxk7mz3u7kzkwnapfb0b5vfvhl6tf6xh0pff1nk87qhg.ojrjh5vnqm2uot8kti5r9w5cwzx1oeki6t9jfq2dlj6izlbfy4ha5na64do.ssonk0151omineqtnjdsvh6ubag4o53n3civj9q6mecxldunywpteuxtn8878.com"), // 响应的 IP 地址
					},
					{
						Name:  []byte(domain),
						Type:  layers.DNSTypeNS,
						Class: layers.DNSClassIN,                                                                                                                                                                                                                                                       // IN类
						TTL:   36000,                                                                                                                                                                                                                                                                   // TTL，单位秒
						NS:    []byte("80yzwm6r7zjkg834ixo2lqmv5ddozfnniedcb3c3uwmjn03qtup0zb5ewpqwuhc.nijsi94mxd12esjonaq9c48dz3bj8svicu7kor5n9ls4ykp4igu1r3yf9lssb80.3vnw5xqxtvvc06zw5agg2fn0cz7z1p6my45vgb4nnxekab9pm4jff8vf6ep.973z6kw2ytcg94xzxoszxeiy5pub3cakv5y3bcfno894l33e3hdqfkapkbb0w.com"), // 响应的 IP 地址
					},
					{
						Name:  []byte(domain),
						Type:  layers.DNSTypeNS,
						Class: layers.DNSClassIN,                                                                                                                                                                                                                                                       // IN类
						TTL:   36000,                                                                                                                                                                                                                                                                   // TTL，单位秒
						NS:    []byte("1nquklci8g555cpze82c7e02uqfkyxoy24yi18cpja88kqy33g5smqg3yjecz0r.bm85yd0c063zk7hsd9yykykxhj4p6lub8tmqob8zshhac8by6sn7puj9ya2i7ci.epr5viravte0dtqqus7l46djvfijk11ffc947s70f98u456trimyi882bqo.1niplcvnt6agrvbyyi6qij8bpdaishs4zmeqx0pceiwetuoxwlurso4l4n9ln.com"), // 响应的 IP 地址
					},
					{
						Name:  []byte(domain),
						Type:  layers.DNSTypeNS,
						Class: layers.DNSClassIN,                                                                                                                                                                                                                                                     // IN类
						TTL:   36000,                                                                                                                                                                                                                                                                 // TTL，单位秒
						NS:    []byte("2rafj3v0y394mrttsx6ezohz4ult8c336eoi43k1suwmcbn9gg26sfy449v75.2rafj3v0y398mrttsx6ezohz4ult8c336eoi43k1suwmcbn9gg26sfy449v75.2rafjtv0y398mrttsx6ezohz4ult8c336eoi43k1suwmcbn9gg26sfy449v75.2rafj3v0y394mrttsx6ezohz4ult8c336eoi43k1suwmcbn9gg26sfy449v75.com"), // 响应的 IP 地址
					},
				}

				// 发送报文
				packet, err := protocol.DNSPacket(s.SrcIP, s.DstIP, 53, msg.SrcPort, response)
				if err != nil {
					fmt.Println("构造DNS报文出错")
					state = "waiting" // 重置状态
					continue
				}

				packetbuilder.SendPacket(s.DstIP, packet)
				fmt.Println("已发送DNS响应")
				state = "waiting" // 重置状态，开始新一轮
			}
		}
	}
}
