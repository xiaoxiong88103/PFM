package route

import (
	"PFM/proxyFunc/proxy"

	"github.com/gin-gonic/gin"
)

func ProxyRoute(r *gin.Engine) {
	proxyRouter := r.Group("/proxy")
	// 设置端口转发的接口
	proxyRouter.POST("/setPort", proxy.SetPortForward)
	// 将关闭后的接口 重新打开
	proxyRouter.POST("/restartPort", proxy.RestartPortForward)
	// 删除端口转发的接口
	proxyRouter.POST("/deletePort", proxy.DeletePortForward)
	// 查询端口转发的接口
	proxyRouter.GET("/getPort", proxy.ListPortForwards)
	//查询转发正在转发的接口
	proxyRouter.GET("/getActivePort", proxy.GetActivePortForwards)
	// 停止转发的接口 暂停某个id的转发
	proxyRouter.POST("/stopPort", proxy.StopPortForward)
	// 查询停止转发的接口目前暂停转发的接口
	proxyRouter.GET("/getStopPort", proxy.GetStoppedPortForwards)

}
