package main

import (
	"PFM/route"
	"PFM/util"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
		// 处理请求
		c.Next()
	}
}

func main() {
	r := gin.Default()
	//加载转发接口
	route.Proxy_Route(r)

	r.Use(middleware())
	// 托管文件
	r.Use(static.Serve("/", static.LocalFile(util.WebPanelPublicPath, true)))

	r.Run(":8281") // 监听并在 0.0.0.0:8080 上启动服务
}
