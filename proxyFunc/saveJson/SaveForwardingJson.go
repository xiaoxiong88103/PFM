package saveJson

import (
	"PFM/proxyFunc/vars"
	"encoding/json"
	"fmt"
	"os"
)

// SavePortForwardingRules 保存规则到文件Json中
func SavePortForwardingRules(rules map[string]vars.PortForwardingRule) error {
	vars.RulesMu.Lock()
	defer vars.RulesMu.Unlock()

	// 重写文件保存逻辑
	//thisTime := time.Now().Format("2006-01-02-15:04:05")
	tmpFilePath := vars.ConfigFilePath + ".tmp"
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("无法创建临时文件: %v", err)
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(rules); err != nil {
		return fmt.Errorf("无法写入规则文件: %v", err)
	}
	_ = file.Close()
	_ = os.Remove(vars.ConfigFilePath)
	return os.Rename(tmpFilePath, vars.ConfigFilePath)
}
