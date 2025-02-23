package route

import (
	"PFM/proxyFunc/proxy"
	"PFM/proxyFunc/whiteList"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProxyRoute(r *gin.Engine) {
	proxyRouter := r.Group("/proxy")
	// 设置端口转发的接口
	proxyRouter.POST("/setPort", proxy.SetPortForward)
	// 删除端口转发的接口
	proxyRouter.POST("/deletePort", proxy.DeletePortForward)
	// 查询端口转发的接口
	proxyRouter.GET("/getPort", proxy.ListPortForwards)
	whiteListRouter := r.Group("whiteList")
	whiteListRouter.GET("/reload", func(c *gin.Context) {
		_ = whiteList.LoadWhiteList()
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "重新加载成功黑白名单以加载最新配置"})
	})
	//Add
	whiteListRouter.POST("/add", whiteList.AddWhiteListHandler)
	// 查看白名单路由
	whiteListRouter.GET("/list", whiteList.ViewWhiteListHandler)
	// 查看全部白名单路由
	whiteListRouter.GET("/list/all", whiteList.ViewAllWhiteListsHandler)
	// 删除白名单路由
	whiteListRouter.POST("/delete", whiteList.DeleteWhiteListHandler)

}
