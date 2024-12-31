package PortVars

import (
	"net"
	"sync"
)

type PortForwardingRule struct {
	ID         string `json:"id"`
	Type       string `json:"type"` // tcp or udp
	RemoteIP   string `json:"remote_ip"`
	RemotePort string `json:"remote_port"`
	LocalPort  string `json:"local_port"`
	Comment    string `json:"comment"`
}

var (
	ConfigFilePath = "/opt/port_forwarding_rules.json"
	Rules          = make(map[string]PortForwardingRule) // 初始化全局规则
	RulesMu        sync.RWMutex                          // 用于保护 rules 的读写锁
	Proxy_wg       sync.WaitGroup                        // 用于管理 goroutines
	TcpListeners   = make(map[string]net.Listener)       // TCP 监听器
	UdpConns       = make(map[string]net.PacketConn)     // UDP 连接
	UdpConnsMu     sync.Mutex                            // 保护 udpConns 的互斥锁
	TcpListenersMu sync.Mutex
)
