package protocol

import (
	"bytes"
	"fmt"
)

// MEMCACHEDResponseBuffer generates Memcached protocol responses
// Amplification factor: up to 51,000x
// Returns multiple response packets for cyclic sending to simulate real attack behavior
func MEMCACHEDResponseBuffer() [][]byte {
	// Create response packets - typically only 2-3 packets in real scenarios
	responses := make([][]byte, 0, 5)
	
	// Basic server info packet - this part is relatively realistic
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
	
	// Command statistics packet - this part is also relatively realistic
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
	
	// Memory statistics packet - this part is also relatively realistic
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
	
	// Add a custom statistics packet for more realistic size
	customStats := new(bytes.Buffer)
	// Custom statistics are usually numbers or short strings
	addMemcachedStat(customStats, "custom_stat_last_reset", "1633458712")
	addMemcachedStat(customStats, "custom_stat_cache_hit_ratio", "99.48")
	// Some values might be slightly longer, but not too long
	addMemcachedStat(customStats, "custom_stat_config", "max_memory=64G;threads=4;lru_crawler=yes;maxconns=1024")
	// A few values might be relatively long
	addMemcachedStat(customStats, "custom_monitoring_info", "node=cache-01;cluster=us-west;dc=portland;env=prod;owner=dbteam")
	customStats.WriteString("END\r\n")
	responses = append(responses, customStats.Bytes())
	
	// Add a more realistic slab statistics packet
	slabStats := new(bytes.Buffer)
	// Real servers typically have around 8-12 slab classes
	for slabClass := 1; slabClass <= 8; slabClass++ {
		prefix := fmt.Sprintf("slab_%d", slabClass)
		
		// Standard slab statistics - these are numeric values
		addMemcachedStat(slabStats, prefix+"_chunk_size", fmt.Sprintf("%d", 96*slabClass))
		addMemcachedStat(slabStats, prefix+"_chunks_per_page", fmt.Sprintf("%d", 10240/slabClass))
		addMemcachedStat(slabStats, prefix+"_total_pages", fmt.Sprintf("%d", slabClass*5))
		addMemcachedStat(slabStats, prefix+"_total_chunks", fmt.Sprintf("%d", slabClass*50))
		addMemcachedStat(slabStats, prefix+"_used_chunks", fmt.Sprintf("%d", slabClass*30))
		addMemcachedStat(slabStats, prefix+"_free_chunks", fmt.Sprintf("%d", slabClass*20))
		addMemcachedStat(slabStats, prefix+"_free_chunks_end", fmt.Sprintf("%d", slabClass*10))
		addMemcachedStat(slabStats, prefix+"_mem_requested", fmt.Sprintf("%d", slabClass*500000))
		
		// Possible 1-2 custom metrics
		if slabClass < 4 {
			addMemcachedStat(slabStats, prefix+"_hit_ratio", fmt.Sprintf("%.2f", 75.5+float64(slabClass*5)))
		}
	}
	slabStats.WriteString("END\r\n")
	responses = append(responses, slabStats.Bytes())
	
	return responses
}

// addMemcachedStat adds a single Memcached statistic line
func addMemcachedStat(buf *bytes.Buffer, name, value string) {
	statLine := fmt.Sprintf("STAT %s %s\r\n", name, value)
	buf.WriteString(statLine)
}

func MEMCACHEDPacket(srcIP, dstIP string, srcPort, dstPort int) ([]byte, error) {
	payload := append(
		[]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00}, 
		[]byte("gets p h e\n")...,                            // Memcached "gets" command
	)
	return BuildUDPPacket(srcIP, dstIP, srcPort, dstPort, payload)
}