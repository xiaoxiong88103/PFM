package vars

import (
	"net"
	"sync"
)

// 转发字段
type PortForwardingRule struct {
	ID         string `json:"id"`          //获取ID索引用的
	Type       string `json:"type"`        // tcp or udp
	RemoteIP   string `json:"remote_ip"`   // 获取远程IP字段
	RemotePort string `json:"remote_port"` //获取远程端口转发的要
	LocalPort  string `json:"local_port"`  //本地的端口绑定
	Comment    string `json:"comment"`     //备注接口负责备注作用的
}

// todo 返回给前端的结构体做特殊处理有些东西不适合给前端看到
type PortForwardingRuleList struct {
	Id         string `json:"id"`          //获取ID索引用的
	Type       string `json:"type"`        // tcp or udp
	Status     uint   `json:"status"`      // 1-true 2-false
	RemoteIP   string `json:"remote_ip"`   // 获取远程IP字段
	RemotePort string `json:"remote_port"` //获取远程端口转发的要
	LocalPort  string `json:"local_port"`  //本地的端口绑定
	Comment    string `json:"comment"`     //备注接口负责备注作用的
}

// 转发逻辑
var (
	configFileName        = "port_rules.json"
	ConfigFilePath        = "/opt/PFM/" + configFileName
	ConfigWindowsFilePath = "./conf/" + configFileName
	Rules                 = make(map[string]PortForwardingRule) // 初始化全局规则
	RulesMu               sync.RWMutex                          // 用于保护 rules 的读写锁
	Proxy_wg              sync.WaitGroup                        // 用于管理 goroutines
	TcpListeners          = make(map[string]net.Listener)       // TCP 监听器
	UdpConns              = make(map[string]net.PacketConn)     // UDP 连接
	UdpConnsMu            sync.Mutex                            // 保护 udpConns 的互斥锁
	TcpListenersMu        sync.Mutex
)
