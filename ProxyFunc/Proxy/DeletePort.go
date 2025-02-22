package Proxy

import (
	"PFM/ProxyFunc/SaveJson"
	"PFM/ProxyFunc/Vars"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "保存端口转发规则失败", "data": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "端口转发删除成功", "data": id})
}
