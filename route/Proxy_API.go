package route

import (
	"PFM/ProxyFunc/Proxy"
	"PFM/ProxyFunc/WhiteList"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Proxy_Route(r *gin.Engine) {
	proxy := r.Group("/proxy")
	// 设置端口转发的接口
	proxy.POST("/setPort", Proxy.SetPortForward)
	// 删除端口转发的接口
	proxy.POST("/deletePort", Proxy.DeletePortForward)
	// 查询端口转发的接口
	proxy.GET("/get_port", Proxy.ListPortForwards)

	whiteList := r.Group("whiteList")

	whiteList.GET("/reload", func(c *gin.Context) {
		WhiteList.LoadWhiteList()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "重新加载成功黑白名单以加载最新配置"})
	})

	//Add
	whiteList.POST("/add", WhiteList.AddWhiteListHandler)
	// 查看白名单路由
	whiteList.GET("/list", WhiteList.ViewWhiteListHandler)

	// 删除白名单路由
	whiteList.POST("/delete", WhiteList.DeleteWhiteListHandler)

}
