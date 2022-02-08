package xgin

import (
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

// RealIP X-Real-Ip
func RealIP(ctx *gin.Context) string {
	r := ctx.Request
	//xProxyAuth := strings.TrimSpace(r.Header.Get("X-Proxy-Auth"))
	xRealIP := strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	remoteAddr := r.RemoteAddr

	ips := strings.Split(xRealIP, ",")
	if len(ips) > 0 && ips[0] != "" {
		oip := strings.TrimSpace(ips[0])
		if strings.Count(oip, ":") < 2 {
			ip, _, err := net.SplitHostPort(oip)
			if err != nil {
				//不带端口
				ip = oip
			}
			return ip
		}
		//考虑是否是ipv6
		ip := net.ParseIP(oip).To16()
		if ip != nil {
			return ip.String()
		}
	}
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return ip
}
