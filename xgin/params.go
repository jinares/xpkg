package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinares/xpkg/xtools"
	"strings"
)

// GetALLParams 获取所有请求参数
func GetALLParams(c *gin.Context) map[string]string {

	if strings.ToLower(c.Request.Method) == "post" {
		c.Request.ParseForm()
		queryJson := xtools.UrlDecode(c.Request.PostForm.Encode())

		return queryJson
	}
	return xtools.UrlDecode(c.Request.URL.RawQuery)
}

// GetQueryPar 获取所有请求参数
func GetQueryPar(c *gin.Context) map[string]string {

	if strings.ToLower(c.Request.Method) == "post" {
		c.Request.ParseForm()
		queryJson := xtools.UrlDecode(c.Request.PostForm.Encode())

		return queryJson
	}
	return xtools.UrlDecode(c.Request.URL.RawQuery)
}
func GetPar(c *gin.Context) map[string]string {
	return xtools.UrlDecode(c.Request.URL.RawQuery)
}

// func PostFormPar(c *gin.Context) map[string]string { 获取所有请求参数
func PostFormPar(c *gin.Context) map[string]string {

	if strings.ToLower(c.Request.Method) == "post" {
		c.Request.ParseForm()
		queryJson := xtools.UrlDecode(c.Request.PostForm.Encode())

		return queryJson
	}
	return map[string]string{}
}
