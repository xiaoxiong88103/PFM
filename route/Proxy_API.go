package route

import (
	"PFM/ProxyFunc/DelProxy"
	"PFM/ProxyFunc/ListProxy"
	"PFM/route/ProxyRoute"
	"github.com/gin-gonic/gin"
)

func Proxy_Route(r *gin.Engine) {
	proxy := r.Group("/proxy")
	// 设置端口转发的接口
	proxy.POST("/setPort", ProxyRoute.SetPortForward)
	// 删除端口转发的接口
	proxy.POST("/deletePort", DelProxy.DeletePortForward)
	// 查询端口转发的接口
	proxy.GET("/get_port", ListProxy.ListPortForwards)

}
