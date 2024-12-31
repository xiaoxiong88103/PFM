package ListProxy

import (
	"PFM/ProxyFunc/PortVars"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListPortForwards(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "查询成功", "data": PortVars.Rules})
}
