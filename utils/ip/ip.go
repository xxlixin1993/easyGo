package ip

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

// ErrIPLimited ip受限
var ErrIPLimited = errors.New("IP has been limited")

// IsInternalIP 检查是否是内网Ip
func IsInternalIP(ipString string) bool {
	ipByte := net.ParseIP(ipString)
	ipArr := ipByte.To4()
	if ipByte == nil {
		return false
	}
	if ipArr[0] == 10 || (ipArr[0] == 172 && ipArr[1] == 16) || (ipArr[0] == 192 && ipArr[1] == 168) || ipString == "127.0.0.1" {
		return true
	}
	return false
}

// RemoteAddr 获取客户端ip
func RemoteAddr(r *http.Request) string {
	ipStr := r.Header.Get("X-Real-IP")
	if ipStr == "" {
		ipStr = InternalIP(r)
	}
	return ipStr
}

// InternalIP 获取内网ip 此方法只能拿到nginx机器的内网ip
func InternalIP(r *http.Request) string {
	ip := r.RemoteAddr
	if ip == "" {
		return ""
	}
	tmp := strings.Split(ip, ":")
	return tmp[0]
}
