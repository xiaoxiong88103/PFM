package SaveJson

import (
	"PFM/ProxyFunc/PortVars"
	"encoding/json"
	"fmt"
	"os"
)

// SavePortForwardingRules 保存规则到文件Json中
func SavePortForwardingRules(rules map[string]PortVars.PortForwardingRule) error {
	PortVars.RulesMu.Lock()
	defer PortVars.RulesMu.Unlock()

	tmpFilePath := PortVars.ConfigFilePath + ".tmp"
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("无法创建临时文件: %v", err)
	}
	defer os.Remove(tmpFilePath)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(rules); err != nil {
		return fmt.Errorf("无法写入规则文件: %v", err)
	}

	return os.Rename(tmpFilePath, PortVars.ConfigFilePath)
}
