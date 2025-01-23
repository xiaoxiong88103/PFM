package Proxy

import (
	"PFM/ProxyFunc/SaveJson"
	"PFM/ProxyFunc/Vars"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//// 删除队列：用于存放待删除的规则 ID
//var deleteQueue = make(chan string, 100) // 缓冲队列，用于存放待删除的端口转发 ID
//// 后台处理删除队列中的任务
//func processDeleteQueue() {
//	for id := range deleteQueue {
//		// 处理删除逻辑
//		Vars.RulesMu.Lock()
//		if rule, exists := Vars.Rules[id]; exists {
//			// 停止转发
//			if rule.Type == "tcp" {
//				if listener, ok := Vars.TcpListeners[id]; ok {
//					listener.Close()
//				}
//			} else if rule.Type == "udp" {
//				Vars.UdpConnsMu.Lock()
//				if conn, ok := Vars.UdpConns[id]; ok {
//					conn.Close()
//				}
//				Vars.UdpConnsMu.Unlock()
//			}
//
//			time.Sleep(1 * time.Millisecond) // 延迟1毫秒
//			// 从规则列表中删除
//			delete(Vars.Rules, id)
//		}
//		Vars.RulesMu.Unlock()
//		Vars.Proxy_wg.Done() // 删除任务完成，减少等待组计数
//	}
//}
//
//// 删除端口转发的接口
//func DeletePortForward(c *gin.Context) {
//	var req struct {
//		ID string `json:"id"`
//	}
//
//	// 解析并绑定请求体
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "请求格式错误", "data": nil})
//		return
//	}
//
//	// 创建一个切片存储要删除的 ID
//	idsToDelete := []string{req.ID, req.ID + "-2", req.ID + "-3"}
//
//	// 遍历所有要删除的 ID
//	for _, id := range idsToDelete {
//		Vars.Proxy_wg.Add(1) // 增加等待组计数
//
//		// 将任务添加到删除队列
//		deleteQueue <- id
//	}
//
//	// 启动处理删除队列的 goroutine（如果还没有启动）
//	go processDeleteQueue()
//
//	// 等待所有删除任务完成
//	Vars.Proxy_wg.Wait()
//
//	// 保存修改后的规则到文件
//	if err := SaveJson.SavePortForwardingRules(Vars.Rules); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "保存端口转发规则失败", "data": err.Error()})
//		return
//	}
//
//	// 返回成功响应
//	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "端口转发删除成功", "data": idsToDelete})
//}
// 删除端口转发的接口（直接删除，不使用队列）

func DeletePortForward(c *gin.Context) {
	var req struct {
		ID string `json:"id"`
	}

	// 解析并绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "请求格式错误", "data": nil})
		return
	}

	id := req.ID // 直接使用传入的 ID

	Vars.RulesMu.Lock()
	defer Vars.RulesMu.Unlock() // 确保在函数退出时释放锁

	if rule, exists := Vars.Rules[id]; exists {
		// 停止转发
		if rule.Type == "tcp" {
			if listener, ok := Vars.TcpListeners[id]; ok {
				listener.Close()
			}
		} else if rule.Type == "udp" {
			Vars.UdpConnsMu.Lock()
			if conn, ok := Vars.UdpConns[id]; ok {
				conn.Close()
			}
			Vars.UdpConnsMu.Unlock()
		}

		time.Sleep(1 * time.Millisecond) // 延迟1毫秒
		// 从规则列表中删除
		delete(Vars.Rules, id)
	} else {
		// 如果 ID 不存在，返回错误信息
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "规则不存在", "data": nil})
		return
	}

	// 保存修改后的规则到文件
	if err := SaveJson.SavePortForwardingRules(Vars.Rules); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "保存端口转发规则失败", "data": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "端口转发删除成功", "data": id})
}
