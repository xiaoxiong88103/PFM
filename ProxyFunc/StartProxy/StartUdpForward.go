package StartProxy

import (
	"PFM/ProxyFunc/Vars"
	"PFM/ProxyFunc/WhiteList"
	"log"
	"net"
	"sync"
	"time"
)

// 处理udp转发逻辑
func StartUDPForward(rule Vars.PortForwardingRule) {
	defer handlePanic("StartUDPForward", rule.ID)

	// 创建本地 UDP 监听器
	localConn, err := net.ListenPacket("udp", ":"+rule.LocalPort)
	if err != nil {
		log.Printf("无法监听本地 UDP 地址 %s: %v", rule.LocalPort, err)
		return
	}
	Vars.UdpConnsMu.Lock()
	Vars.UdpConns[rule.ID] = localConn
	Vars.UdpConnsMu.Unlock()

	defer func() {
		localConn.Close()
		Vars.UdpConnsMu.Lock()
		delete(Vars.UdpConns, rule.ID)
		Vars.UdpConnsMu.Unlock()
	}()

	log.Printf("UDP 转发启动: %s -> %s:%s", rule.LocalPort, rule.RemoteIP, rule.RemotePort)

	clientToRemote := sync.Map{}
	lastActivity := sync.Map{}
	stopChannel := make(chan struct{}) // 用于优雅关闭

	// 定时清理超时连接
	go func() {
		defer handlePanic("UDP Cleaner", rule.ID)

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentTime := time.Now().Unix()
				lastActivity.Range(func(key, value interface{}) bool {
					if lastTime, ok := value.(int64); ok && currentTime-lastTime > 60 {
						if conn, exists := clientToRemote.Load(key); exists {
							log.Printf("连接超时，关闭 UDP 连接: %s", key)
							conn.(net.Conn).Close()
							clientToRemote.Delete(key)
							lastActivity.Delete(key)
						}
					}
					return true
				})
			case <-stopChannel:
				return
			}
		}
	}()

	// 读取本地 UDP 数据包并进行转发
	bufferPool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 4096)
		},
	}

	for {
		buffer := bufferPool.Get().([]byte)
		n, clientAddr, err := localConn.ReadFrom(buffer)
		if err != nil {
			if IsClosedConnErr(err) {
				log.Printf("本地 UDP 连接已关闭: %s", rule.LocalPort)
				bufferPool.Put(buffer)
				return
			}
			log.Printf("读取本地 UDP 连接失败: %v", err)
			bufferPool.Put(buffer)
			continue
		}

		clientKey := clientAddr.String()

		//校验白名单IP
		clientIP, _, _ := net.SplitHostPort(clientKey)
		// 校验 IP 是否在白名单中
		if !WhiteList.IsIPAllowed(rule.LocalPort, clientIP) {
			log.Printf("拒绝连接，客户端 IP %s 不在端口 %s 的白名单中", clientIP, rule.LocalPort)
			bufferPool.Put(buffer)
			continue
		}

		remoteConn, exists := clientToRemote.Load(clientKey)

		if !exists {
			// 如果没有现有的连接，则建立到远程的连接
			remoteConn, err = net.Dial("udp", net.JoinHostPort(rule.RemoteIP, rule.RemotePort))
			if err != nil {
				log.Printf("无法连接远程 UDP 地址 %s:%s: %v", rule.RemoteIP, rule.RemotePort, err)
				bufferPool.Put(buffer)
				continue
			}
			clientToRemote.Store(clientKey, remoteConn)
			lastActivity.Store(clientKey, time.Now().Unix())

			go func(clientAddr net.Addr, remoteConn net.Conn, clientKey string) {
				defer handlePanic("UDP Response Handler", clientKey)
				defer remoteConn.Close()

				remoteBuffer := bufferPool.Get().([]byte)
				defer bufferPool.Put(remoteBuffer)

				for {
					n, err := remoteConn.Read(remoteBuffer)
					if err != nil {
						log.Printf("读取远程 UDP 连接失败: %v", err)
						clientToRemote.Delete(clientKey)
						lastActivity.Delete(clientKey)
						return
					}

					_, err = localConn.WriteTo(remoteBuffer[:n], clientAddr)
					if err != nil {
						log.Printf("写入本地 UDP 连接失败: %v", err)
					}

					lastActivity.Store(clientKey, time.Now().Unix())
				}
			}(clientAddr, remoteConn.(net.Conn), clientKey)
		}

		_, err = remoteConn.(net.Conn).Write(buffer[:n])
		if err != nil {
			log.Printf("写入远程 UDP 连接失败: %v", err)
			clientToRemote.Delete(clientKey)
			lastActivity.Delete(clientKey)
			remoteConn.(net.Conn).Close()
		}

		lastActivity.Store(clientKey, time.Now().Unix())
		bufferPool.Put(buffer)
	}
}

// handlePanic 用于捕获和记录 panic，以确保 goroutine 崩溃不会影响主进程
func handlePanic(funcName, ruleID string) {
	if r := recover(); r != nil {
		log.Printf("端口转发规则 %s 在 %s 中发生崩溃: %v", ruleID, funcName, r)
	}
}
