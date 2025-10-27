package LogsWrite

import (
	"PFM/proxyFunc/vars"
	"github.com/go-ini/ini"
	"log"
)

func LoadConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("读取 config.ini 失败: %v，使用默认配置", err)
		vars.LogsWrite = true // 默认启用日志
		return
	}

	vars.LogsWrite = cfg.Section("config").Key("logs").MustBool(true)
	log.Printf("日志功能: %v", vars.LogsWrite)
}

// Logf 统一日志入口，根据 LogsWrite 决定是否输出
func Logf(format string, v ...interface{}) {
	if vars.LogsWrite {
		log.Printf(format, v...)
	}
}
