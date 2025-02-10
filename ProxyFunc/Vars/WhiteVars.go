package Vars

var (
	WhiteListFileName        = "WhiteList.ini"
	WhiteListFilePath        = "/opt/PFM/" + WhiteListFileName
	WhiteListWindowsFilePath = "./conf/" + WhiteListFileName
	WhiteList                = make(map[string][]string)
)

var WhiteList_Json struct {
	Port string `json:"port" binding:"required"`
	IP   string `json:"ip" binding:"required"`
}
