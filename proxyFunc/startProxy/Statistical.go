package startProxy

import "sync"

type TrafficStat struct {
	Up   int64 // 上行
	Down int64 // 下行
}

var (
	PortTrafficStats = make(map[string]*TrafficStat)
	PortTrafficMu    sync.Mutex
)

// 更新统计记录 流量大小
func UpdatePortTraffic(port string, up, down int64) {
	PortTrafficMu.Lock()
	defer PortTrafficMu.Unlock()

	stat, exists := PortTrafficStats[port]
	if !exists {
		stat = &TrafficStat{}
		PortTrafficStats[port] = stat
	}
	stat.Up += up
	stat.Down += down
}
