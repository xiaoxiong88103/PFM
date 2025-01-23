package Vars

var (
	WhiteList_files = "/opt/WhiteList.ini"
	WhiteList       = make(map[string][]string)
)

var WhiteList_Json struct {
	Port string `json:"port" binding:"required"`
	IP   string `json:"ip" binding:"required"`
}
