package xoidc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func GetUserAgent(mds map[string]string) string {
	if ua := mds["grpcgateway-user-agent"]; ua != "" {
		return ua
	}
	return mds["user-agent"]
}

// GetIPFromMeta returns IP address from request.
// Only when it used use proxy
func GetIPFromMeta(mds map[string]string) string {
	if ip := mds["x-forwarded-for"]; ip != "" {
		parts := strings.Split(ip, ",")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

func GetDeviceID(mds map[string]string) string {
	if id := mds["x-device-id"]; id != "" {
		return id
	}
	return mds["grpcgateway-x-device-id"]
}

func VerifyHas(par string, data []string) error {
	if len(data) == 0 {
		return nil
	}
	datamap := make(map[string]bool)
	for _, s := range data {
		if s == "all" {
			return nil
		}
		if s == "" {
			continue
		}
		datamap[s] = true
	}
	if len(datamap) == 0 {
		return nil
	}
	var has bool = false
	for _, s := range data {
		if datamap[s] {
			has = true
			break
		}
	}
	if !has {
		return status.Errorf(codes.PermissionDenied, "user  want [%s], but go [%s]", par, data)
	}
	return nil
}
