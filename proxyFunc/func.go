package proxyFunc

import (
	"PFM/proxyFunc/LogsWrite"
)

func InitPublic() {
	LogsWrite.LoadConfig()
}
