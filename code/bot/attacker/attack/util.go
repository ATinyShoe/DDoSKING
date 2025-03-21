package attack

import(
	"fmt"
	"time"
	"math/rand"
)
var STOP = make(chan struct{}) // 用于停止攻击

// 设置随机种子
func init() {
    rand.Seed(time.Now().UnixNano())
}

// 伪造地址和端口
func RandIPv4() string {
	rand.Seed(time.Now().UnixNano()) // 设置随机种子
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(256), // 每个段的范围是 0-255
		rand.Intn(256),
		rand.Intn(256),
		rand.Intn(256),
	)
}

func RandPort() int {
	return rand.Intn(65535-1) + 1    // 生成 1 到 65535 的随机端口
}

// ResetStopChannel 重置停止通道
func ResetStopChannel() {
	// 重置停止通道
	select {
	case <-STOP:
		// 通道已经关闭，需要重建
		STOP = make(chan struct{})
	default:
		// 通道仍然打开，不需要做任何事情
	}
}