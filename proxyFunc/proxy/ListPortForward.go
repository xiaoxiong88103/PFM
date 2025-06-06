package proxy

import (
	"PFM/proxyFunc/vars"
	"PFM/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListPortForwards(c *gin.Context) {
	Forwards := make([]*vars.PortForwardingRuleList, 0)
	// 检查所有规则，找出已经停止的转发
	for id, rule := range vars.Rules {
		it := &vars.PortForwardingRuleList{
			Id:         rule.ID,
			Type:       rule.Type,
			RemoteIP:   rule.RemoteIP,
			RemotePort: rule.RemotePort,
			Status:     util.FalseEum,
			LocalPort:  rule.LocalPort,
			Comment:    rule.Comment,
		}
		switch rule.Type {
		case util.UdpEum:
			vars.UdpConnsMu.Lock()
			if _, ok := vars.UdpConns[id]; ok {
				it.Status = util.TrueEum
			}
			vars.UdpConnsMu.Unlock()
		case util.TcpEum:
			if _, ok := vars.TcpListeners[id]; ok {
				it.Status = util.TrueEum
			}
		}
		Forwards = append(Forwards, it)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "查询成功", "data": Forwards})
}

func GetStoppedPortForwards(c *gin.Context) {
	stoppedForwards := make([]vars.PortForwardingRule, 0)
	// 检查所有规则，找出已经停止的转发
	for id, rule := range vars.Rules {
		if rule.Type == "tcp" {
			if _, ok := vars.TcpListeners[id]; !ok {
				stoppedForwards = append(stoppedForwards, rule)
			}
		} else if rule.Type == "udp" {
			vars.UdpConnsMu.Lock()
			if _, ok := vars.UdpConns[id]; !ok {
				stoppedForwards = append(stoppedForwards, rule)
			}
			vars.UdpConnsMu.Unlock()
		}
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "获取已停止的端口转发成功", "data": stoppedForwards})
}

func GetActivePortForwards(c *gin.Context) {
	activeForwards := make([]vars.PortForwardingRule, 0)
	// 检查所有规则，找出正在运行的转发
	for id, rule := range vars.Rules {
		if rule.Type == "tcp" {
			if _, ok := vars.TcpListeners[id]; ok {
				activeForwards = append(activeForwards, rule)
			}
		} else if rule.Type == "udp" {
			vars.UdpConnsMu.Lock()
			if _, ok := vars.UdpConns[id]; ok {
				activeForwards = append(activeForwards, rule)
			}
			vars.UdpConnsMu.Unlock()
		}
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "获取正在运行的端口转发成功", "data": activeForwards})
}
