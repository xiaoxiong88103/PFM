package startProxy

import (
	"net"
	"strings"
)

// 判断网络是否被关闭Error返回
func IsClosedConnErr(err error) bool {
	if opErr, ok := err.(*net.OpError); ok && strings.Contains(opErr.Err.Error(), "use of closed network connection") {
		return true
	}
	return false
}
