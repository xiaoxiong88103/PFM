package proxy

import (
	"PFM/proxyFunc/vars"
	"encoding/json"
	"log"
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
