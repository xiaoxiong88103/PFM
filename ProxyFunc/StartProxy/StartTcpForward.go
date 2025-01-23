package StartProxy

import (
	"PFM/ProxyFunc/Vars"
	"PFM/ProxyFunc/WhiteList"
	"io"
	"log"
	"net"
	"time"
)

// 开启tcp规则转发
func StartTCPForward(rule Vars.PortForwardingRule) {

	// 尝试监听本地端口
	listener, err := net.Listen("tcp", ":"+rule.LocalPort)
	if err != nil {
		log.Printf("无法监听本地地址 %s: %v", rule.LocalPort, err)
		return
	}

	// 锁住 tcpListeners map，防止并发写入引发错误
	Vars.TcpListenersMu.Lock()
	Vars.TcpListeners[rule.ID] = listener
	Vars.TcpListenersMu.Unlock()

	// 确保 listener 关闭时也删除 map 中的记录
	defer func() {
		listener.Close()
		Vars.TcpListenersMu.Lock()
		delete(Vars.TcpListeners, rule.ID)
		Vars.TcpListenersMu.Unlock()
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
		if !WhiteList.IsIPAllowed(rule.LocalPort, clientAddr) {
			log.Printf("拒绝连接，客户端IP %s 不在端口 %s 的白名单中", clientAddr, rule.LocalPort)
			clientConn.Close()
			continue
		}

		go handleTCPConnection(clientConn, rule)
	}
}

// 处理tcp转发逻辑
func handleTCPConnection(clientConn net.Conn, rule Vars.PortForwardingRule) {
	defer clientConn.Close()

	var remoteConn net.Conn
	var err error
	for i := 0; i < 100; i++ {
		remoteConn, err = net.Dial("tcp", net.JoinHostPort(rule.RemoteIP, rule.RemotePort))
		if err == nil {
			break
		}
		log.Printf("连接到远程地址 %s:%s 失败 (第 %d 次重试): %v", rule.RemoteIP, rule.RemotePort, i+1, err)
		time.Sleep(time.Duration((i+1)*100) * time.Millisecond)
	}
	if err != nil {
		log.Printf("连接到远程地址 %s:%s 失败，超过最大重试次数: %v", rule.RemoteIP, rule.RemotePort, err)
		return
	}
	defer remoteConn.Close()

	go io.Copy(remoteConn, clientConn)
	io.Copy(clientConn, remoteConn)
}
