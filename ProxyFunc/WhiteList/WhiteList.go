package WhiteList

import (
	"PFM/ProxyFunc/Vars"
	"PFM/util"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
)

// 初始化时加载白名单
func LoadWhiteList() error {
	path, err := util.InitConfigFiles(Vars.WhiteListFilePath, Vars.WhiteListWindowsFilePath, "")
	if err != nil {
		return err
	}
	Vars.WhiteListFilePath = path
	// 加载配置文件
	cfg, err := ini.Load(Vars.WhiteListFilePath)
	if err != nil {
		return err
	}
	Vars.WhiteList = make(map[string][]string)
	if section, err := cfg.GetSection("white_list"); err == nil {
		for key, value := range section.KeysHash() {
			Vars.WhiteList[key] = strings.Split(value, ",")
		}
	}

	log.Println("白名单加载成功:", Vars.WhiteList)
	return nil
}

// 校验IP是否在白名单中
func IsIPAllowed(port string, clientIP string) bool {
	allowedIPs, exists := Vars.WhiteList[port]
	if !exists {
		// 如果没有设置白名单，允许所有 IP
		return true
	}
	for _, ip := range allowedIPs {
		if ip == clientIP {
			return true
		}
	}
	return false
}

// AddWhiteListHandler 处理添加白名单的 POST 请求
func AddWhiteListHandler(c *gin.Context) {
	// 绑定请求体到 WhiteList_Json
	if err := c.ShouldBindJSON(&Vars.WhiteList_Json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "请求参数错误", "data": err.Error()})
		return
	}

	port := Vars.WhiteList_Json.Port
	newIPs := strings.Split(Vars.WhiteList_Json.IP, ",")

	// 加载配置文件
	cfg, err := ini.Load(Vars.WhiteListFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "加载配置文件失败", "data": err.Error()})
		return
	}

	// 获取或创建 [white_list] 部分
	section, err := cfg.GetSection("white_list")
	if err != nil {
		section, _ = cfg.NewSection("white_list") // 如果不存在，创建新部分
	}

	// 获取当前端口的白名单 IP 列表
	existingIPs := section.Key(port).Value()
	ipList := strings.Split(existingIPs, ",")
	ipSet := make(map[string]bool) // 用于快速查重

	for _, ip := range ipList {
		if ip != "" {
			ipSet[ip] = true
		}
	}

	// 分析请求中要添加的 IP 列表
	var alreadyAdded []string
	var newlyAdded []string

	for _, ip := range newIPs {
		if ipSet[ip] {
			alreadyAdded = append(alreadyAdded, ip)
		} else {
			newlyAdded = append(newlyAdded, ip)
			ipSet[ip] = true
		}
	}

	// 如果有新的 IP，需要更新配置文件
	if len(newlyAdded) > 0 {
		allIPs := make([]string, 0, len(ipSet))
		for ip := range ipSet {
			allIPs = append(allIPs, ip)
		}

		section.Key(port).SetValue(strings.Join(allIPs, ","))

		// 保存到配置文件
		if err := cfg.SaveTo(Vars.WhiteListFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "保存配置文件失败", "data": err.Error()})
			return
		}
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "处理完成",
		"data": gin.H{
			"already_added": alreadyAdded,
			"newly_added":   newlyAdded,
		},
	})
}

// ViewWhiteListHandler 查看白名单的处理函数
func ViewWhiteListHandler(c *gin.Context) {
	port := c.Query("port") // 从查询参数中获取端口号
	if port == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "缺少端口号参数", "data": nil})
		return
	}

	// 加载配置文件
	cfg, err := ini.Load(Vars.WhiteListFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "加载配置文件失败", "data": err.Error()})
		return
	}

	// 获取白名单部分
	section, err := cfg.GetSection("white_list")
	if err != nil || !section.HasKey(port) {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "未找到对应端口的白名单", "data": nil})
		return
	}

	// 获取指定端口的白名单 IP 列表
	whiteList := section.Key(port).Value()
	ipList := strings.Split(whiteList, ",")

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "获取成功", "data": ipList})
}

// DeleteWhiteListHandler 删除白名单的处理函数
func DeleteWhiteListHandler(c *gin.Context) {
	// 绑定请求体到 WhiteList_Json
	if err := c.ShouldBindJSON(&Vars.WhiteList_Json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "请求参数错误", "data": err.Error()})
		return
	}

	port := Vars.WhiteList_Json.Port
	deleteIPs := strings.Split(Vars.WhiteList_Json.IP, ",")

	// 加载配置文件
	cfg, err := ini.Load(Vars.WhiteListFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "加载配置文件失败", "data": err.Error()})
		return
	}

	// 获取或创建 [white_list] 部分
	section, err := cfg.GetSection("white_list")
	if err != nil || !section.HasKey(port) {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "未找到对应端口的白名单", "data": nil})
		return
	}

	// 获取当前端口的白名单 IP 列表
	existingIPs := section.Key(port).Value()
	ipList := strings.Split(existingIPs, ",")
	ipSet := make(map[string]bool)

	// 将现有 IP 放入集合
	for _, ip := range ipList {
		if ip != "" {
			ipSet[ip] = true
		}
	}

	// 删除请求中的 IP
	var deletedIPs []string
	var notFoundIPs []string
	for _, ip := range deleteIPs {
		if ipSet[ip] {
			deletedIPs = append(deletedIPs, ip)
			delete(ipSet, ip)
		} else {
			notFoundIPs = append(notFoundIPs, ip)
		}
	}

	// 更新配置文件
	allIPs := make([]string, 0, len(ipSet))
	for ip := range ipSet {
		allIPs = append(allIPs, ip)
	}

	if len(allIPs) > 0 {
		section.Key(port).SetValue(strings.Join(allIPs, ","))
	} else {
		section.DeleteKey(port) // 如果删除后没有 IP，移除该端口的白名单
	}

	if err := cfg.SaveTo(Vars.WhiteListFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "保存配置文件失败", "data": err.Error()})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除完成",
		"data": gin.H{
			"deleted":   deletedIPs,
			"not_found": notFoundIPs,
		},
	})
}
