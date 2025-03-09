package proxy

import (
	"PFM/proxyFunc/vars"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func LoadPortForwardingRules() (map[string]vars.PortForwardingRule, error) {
	vars.Rules = make(map[string]vars.PortForwardingRule)
	file, err := os.Open(vars.ConfigFilePath)
	if os.IsNotExist(err) {
		// 如果文件不存在，则创建一个空的配置文件
		file, err = os.Create(vars.ConfigFilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(vars.Rules)
		if err != nil {
			return nil, err
		}
		return vars.Rules, nil
	}
	if os.IsNotExist(err) {
		return vars.Rules, nil // 文件不存在则返回空规则列表
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&vars.Rules)
	if err != nil {
		return nil, err
	}
	return vars.Rules, nil
}

func InitReloadProxy() {
	// 加载现有的端口转发规则
	var err error
	_, err = LoadPortForwardingRules()
	if err != nil {
		log.Fatalf("无法加载端口转发规则: %v", err)
	}

	// 恢复端口转发规则
	for _, rule := range vars.Rules {
		StartForwarding(rule)
		time.Sleep(1 * time.Millisecond) // 延迟1毫秒
	}

}

// 重新将 停止的 转发开起来 传入id 即可
func RestartPortForward(c *gin.Context) {
	var req struct {
		ID string `json:"id"`
	}

	// 解析并绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "请求格式错误", "data": nil})
		return
	}

	id := req.ID // 获取传入的 ID
	if rule, exists := vars.Rules[id]; exists {
		// 重新启动转发
		StartForwarding(rule)
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "端口转发已重新启动", "data": rule})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "规则不存在", "data": nil})
	}
}
