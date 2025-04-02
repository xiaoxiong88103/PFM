package main

import (
	"PFM/proxyFunc/proxy"
	"PFM/proxyFunc/vars"
	"PFM/proxyFunc/whiteList"
	"log"
)

func init() {
	//检查config文件是否缺少
	err := vars.CheckAndCreateFileJson(vars.ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	inierr := vars.CheckAndCreateINI(vars.WhiteListFilePath)
	if inierr != nil {
		log.Fatal(inierr)
	}
	//开机转发恢复
	proxy.InitReloadProxy()
	//白名单检查列
	loaderr := whiteList.LoadWhiteList()
	if loaderr != nil {
		log.Fatal(loaderr)
	}
}
