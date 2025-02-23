package main

import (
	"PFM/proxyFunc/proxy"
	"PFM/proxyFunc/vars"
	"PFM/proxyFunc/whiteList"
)

func init() {
	//检查config文件是否缺少
	_ = vars.CheckAndCreateFileJson(vars.ConfigFilePath)
	_ = vars.CheckAndCreateINI(vars.WhiteListFilePath)
	//开机转发恢复
	proxy.InitReloadProxy()
	//白名单检查列
	_ = whiteList.LoadWhiteList()
}
