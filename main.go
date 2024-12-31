package main

import (
	"PFM/ProxyFunc/PortVars"
	"PFM/route"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//检查config文件是否缺少
	PortVars.CheckAndCreateFile(PortVars.ConfigFilePath)
	//加载转发接口
	route.Proxy_Route(r)

	r.Run(":8281") // 监听并在 0.0.0.0:8080 上启动服务
}
