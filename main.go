package main

import (
	"PFM/route"
	"PFM/util"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 托管文件
	r.Use(static.Serve("/", static.LocalFile(util.WebPanelPublicPath, true)))
	//加载转发接口
	route.Proxy_Route(r)
	r.Run(":8281")
}
