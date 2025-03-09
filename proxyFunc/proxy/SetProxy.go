package proxy

import (
	"PFM/proxyFunc/saveJson"
	"PFM/proxyFunc/startProxy"
	"PFM/proxyFunc/vars"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetPortForward(c *gin.Context) {
	var req vars.PortForwardingRule
	// 解析并绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "请求格式错误", "data": nil})
		return
	}
	// 检查并初始化 rules
	if vars.Rules == nil {
		vars.Rules = make(map[string]vars.PortForwardingRule)
	}
	// 校验请求体字段
	if req.ID == "" || req.Type == "" || req.RemoteIP == "" || req.RemotePort == "" || req.LocalPort == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "ID、类型、IP 或端口字段缺失或为空", "data": nil})
		return
	}
	// 1. **先检查本地端口冲突**:
	//    如果已有相同类型的转发规则使用了相同的本地端口，则返回冲突
	vars.RulesMu.RLock() // 读锁
	for _, existingRule := range vars.Rules {
		if existingRule.Type == req.Type && existingRule.LocalPort == req.LocalPort {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 1,
				"msg":  fmt.Sprintf("本地端口 %s 已被占用，请更换端口", req.LocalPort),
				"data": nil,
			})
			vars.RulesMu.RUnlock() // 别忘了在 return 之前释放读锁
			return
		}
	}
	vars.RulesMu.RUnlock()

	// 检查规则是否已存在
	if _, exists := vars.Rules[req.ID]; exists {
		fmt.Println("已存在的规则:", vars.Rules[req.ID])
		//c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "规则已存在", "data": nil})
		// 兼容修改端口功能
		vars.Rules[req.ID] = req
		id := req.ID // 获取传入的 ID
		if rule, exists := vars.Rules[id]; exists {
			// 重新启动转发
			StartForwarding(rule)
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "端口转发已重新启动", "data": rule})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "规则不存在", "data": nil})
		}
		return
	}

	// 启动转发
	StartForwarding(req)

	// 保存规则到内存和文件
	vars.Rules[req.ID] = req
	if err := saveJson.SavePortForwardingRules(vars.Rules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "保存端口转发规则失败", "data": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "端口转发设置成功", "data": nil})
}

// 开启转发的规则
func StartForwarding(rule vars.PortForwardingRule) {
	if rule.Type == "tcp" {
		go startProxy.StartTCPForward(rule)
	} else if rule.Type == "udp" {
		go startProxy.StartUDPForward(rule)
	}

}
