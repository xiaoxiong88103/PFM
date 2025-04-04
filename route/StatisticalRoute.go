package route

import (
	"PFM/proxyFunc/startProxy"
	"github.com/gin-gonic/gin"
)

func StatisticalRoute(r *gin.Engine) {

	r.GET("/tcp/stats", func(c *gin.Context) {
		startProxy.PortTrafficMu.Lock()
		defer startProxy.PortTrafficMu.Unlock()

		result := make(map[string]gin.H)
		for port, stat := range startProxy.PortTrafficStats {
			result[port] = gin.H{
				"up_bytes":   stat.Up,
				"down_bytes": stat.Down,
				"total":      stat.Up + stat.Down,
			}
		}
		c.JSON(200, result)
	})

}
