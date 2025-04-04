package startProxy

import (
	"PFM/proxyFunc/vars"
	"PFM/proxyFunc/whiteList"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// 开启tcp规则转发
func StartTCPForward(rule vars.PortForwardingRule) {

	// 尝试监听本地端口
	listener, err := net.Listen("tcp", ":"+rule.LocalPort)
	if err != nil {
		log.Printf("无法监听本地地址 %s: %v", rule.LocalPort, err)
		return
	}

	// 锁住 tcpListeners map，防止并发写入引发错误
	vars.TcpListenersMu.Lock()
	vars.TcpListeners[rule.ID] = listener
	vars.TcpListenersMu.Unlock()

	// 确保 listener 关闭时也删除 map 中的记录
	defer func() {
		listener.Close()
		vars.TcpListenersMu.Lock()
		delete(vars.TcpListeners, rule.ID)
		vars.TcpListenersMu.Unlock()
	}()

	log.Printf("TCP 转发启动: %s -> %s:%s", rule.LocalPort, rule.RemoteIP, rule.RemotePort)
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if IsClosedConnErr(err) {
				log.Printf("TCP 监听器已关闭: %s", rule.LocalPort)
				return
			}
			log.Printf("接受客户端连接失败: %v", err)
			continue
		}

		// 获取客户端 IP
		clientAddr := clientConn.RemoteAddr().(*net.TCPAddr).IP.String()
		// 校验 IP 是否在白名单中
		if !whiteList.IsIPAllowed(rule.LocalPort, clientAddr) {
			log.Printf("拒绝连接，客户端IP %s 不在端口 %s 的白名单中", clientAddr, rule.LocalPort)
			clientConn.Close()
			continue
		}

		go handleTCPConnection(clientConn, rule)
	}
}

// 处理 TCP 转发连接（带详细日志）
func handleTCPConnection(clientConn net.Conn, rule vars.PortForwardingRule) {
	defer clientConn.Close()

	startTime := time.Now()
	startStr := startTime.Format("2006-01-02 15:04:05")
	clientAddr := clientConn.RemoteAddr().String()

	// 打印连接建立日志
	log.Printf("[连接建立] 时间: %s | 本地端口: %s | 客户端: %s | 目标: %s:%s",
		startStr, rule.LocalPort, clientAddr, rule.RemoteIP, rule.RemotePort)

	// ================= 连接远程端 ====================
	var remoteConn net.Conn
	var err error
	for i := 0; i < 100; i++ {
		remoteConn, err = net.Dial("tcp", net.JoinHostPort(rule.RemoteIP, rule.RemotePort))
		if err == nil {
			break
		}
		log.Printf("连接远程失败 %s:%s (第 %d 次): %v", rule.RemoteIP, rule.RemotePort, i+1, err)
		time.Sleep(time.Duration((i+1)*100) * time.Millisecond)
	}
	if err != nil {
		log.Printf("连接远程失败超限，终止连接: %v", err)
		return
	}
	defer remoteConn.Close()

	// ================= 开始转发 ====================
	var upBytes, downBytes int64
	var wg sync.WaitGroup
	wg.Add(2)

	// 客户端 → 远程（上行）
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("上行转发异常恢复: %v", r)
			}
			wg.Done()
		}()
		n, _ := io.Copy(remoteConn, clientConn)
		atomic.AddInt64(&upBytes, n)
	}()

	// 远程 → 客户端（下行）
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("下行转发异常恢复: %v", r)
			}
			wg.Done()
		}()
		n, _ := io.Copy(clientConn, remoteConn)
		atomic.AddInt64(&downBytes, n)
	}()

	// 等待转发完成
	wg.Wait()

	// ================= 打印关闭日志 ====================
	duration := time.Since(startTime)
	log.Printf("[连接关闭] 时间: %s | 本地端口: %s | 客户端: %s | 目标: %s:%s | 用时: %s | 上行: %dB | 下行: %dB | 总: %dB",
		time.Now().Format("2006-01-02 15:04:05"),
		rule.LocalPort,
		clientAddr,
		rule.RemoteIP, rule.RemotePort,
		duration,
		upBytes, downBytes, upBytes+downBytes)

	// 累加全局端口统计
	UpdatePortTraffic(rule.LocalPort, upBytes, downBytes)
}
