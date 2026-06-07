package utils

import (
	"net"
	"strings"
)

// IpVerify ip验证
func IpVerify(ipStr string) bool {
	ipStr = strings.TrimSpace(ipStr)
	if net.ParseIP(ipStr) != nil {
		return true
	} else {
		return false
	}
}
