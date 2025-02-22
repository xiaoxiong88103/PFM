package vars

var (
	WhiteListFileName        = "whitelist_rules.ini"
	WhiteListFilePath        = "/opt/PFM/" + WhiteListFileName
	WhiteListWindowsFilePath = "./conf/" + WhiteListFileName
	WhiteList                = make(map[string][]string)
)

var WhiteListJson struct {
	Port string `json:"port" binding:"required"`
	IP   string `json:"ip" binding:"required"`
}
