package main

import (
	"PFM/ProxyFunc/Proxy"
	"PFM/ProxyFunc/Vars"
	"PFM/ProxyFunc/WhiteList"
)

func init() {
	//检查config文件是否缺少
	_ = Vars.CheckAndCreateFileJson(Vars.ConfigFilePath)
	_ = Vars.CheckAndCreateINI(Vars.WhiteListFilePath)
	//开机转发恢复
	Proxy.InitReloadProxy()
	//白名单检查列
	_ = WhiteList.LoadWhiteList()
}
