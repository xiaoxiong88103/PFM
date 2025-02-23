package proxy

import (
	"PFM/proxyFunc/vars"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListPortForwards(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "查询成功", "data": vars.Rules})
}
