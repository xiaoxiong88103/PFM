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
	proxy.GET("/getPort", Proxy.ListPortForwards)

	whiteList := r.Group("whiteList")
	//重新加载黑白名单
	whiteList.GET("/reload", func(c *gin.Context) {
		WhiteList.LoadWhiteList()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "重新加载成功黑白名单以加载最新配置"})
	})
	//Add whiteList
	whiteList.POST("/add", WhiteList.AddWhiteListHandler)
	// 查看白名单路由
	whiteList.GET("/list", WhiteList.ViewWhiteListHandler)
	// 查看全部白名单路由
	whiteList.GET("/list/all", WhiteList.ViewAllWhiteListsHandler)
	// 删除白名单路由
	whiteList.POST("/delete", WhiteList.DeleteWhiteListHandler)
	//白名单的限制端口次数查询
	whiteList_number := whiteList.Group("/number")
	whiteList_number.GET("/status", func(c *gin.Context) {
		port := c.Query("port")
		port_number := WhiteList.QueryConnectionCount(port)
		c.JSON(http.StatusOK, gin.H{"port_number": port_number})
	})
	whiteList_number.GET("/clear", func(c *gin.Context) {
		port := c.Query("port")
		WhiteList.ResetConnectionCount(port)
		c.String(200, "恭喜清理成功")
	})

}
