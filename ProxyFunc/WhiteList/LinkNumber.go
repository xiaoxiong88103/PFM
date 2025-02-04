package WhiteList

import (
	"github.com/go-ini/ini"
	"log"
	"strconv"
	"sync"
)

// 用来记录 IP 的连接次数
var connectionCounts sync.Map

// 更新 IP 的连接次数，并返回最新的次数
func UpdateConnectionCount(port string) int {
	count := 1
	if v, ok := connectionCounts.Load(port); ok {
		count = v.(int) + 1
	}
	connectionCounts.Store(port, count)
	return count
}

// 清除 IP 的连接次数
func ResetConnectionCount(port string) {
	connectionCounts.Delete(port)
}

// 查询 IP 的当前连接次数
func QueryConnectionCount(port string) int {
	if v, ok := connectionCounts.Load(port); ok {
		return v.(int)
	}
	return 0 // 如果没有记录，则返回 0
}

// 从配置文件获取指定端口的最大允许次数，并检查当前连接数是否超限
func IsPortWithinLimit(filePath string, port string) bool {
	// 加载配置文件
	cfg, err := ini.Load(filePath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 获取 [ProxyNumber] 节的配置
	section, err := cfg.GetSection("ProxyNumber")
	if err != nil {
		log.Fatalf("未找到 ProxyNumber 节: %v", err)
	}

	// 获取端口对应的最大允许次数
	maxCountStr := section.Key(port).String()
	if maxCountStr == "" {
		log.Printf("未找到端口 %s 的最大连接限制配置，允许连接", port)
		return true // 如果没有配置该端口，则不限制
	}

	maxCount, err := strconv.Atoi(maxCountStr)
	if err != nil {
		log.Fatalf("端口 %s 的最大连接次数配置错误: %v", port, err)
	}

	// 获取当前连接次数
	currentCount := 0
	if v, ok := connectionCounts.Load(port); ok {
		currentCount = v.(int)
	}

	// 检查是否超过限制
	if currentCount >= maxCount {
		log.Printf("端口 %s 当前连接次数 %d，已超过最大限制 %d", port, currentCount, maxCount)
		return false // 超过限制，返回 false
	}

	// 未超出限制，增加连接次数并返回 true
	connectionCounts.Store(port, currentCount+1)
	return true
}
