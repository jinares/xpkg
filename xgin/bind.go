package xgin

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func BindJSON(c *gin.Context, out interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	defer func() {
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}()
	return c.BindJSON(out)
}
func Bind(c *gin.Context, out interface{}) error {
	if strings.EqualFold(c.Request.Method, http.MethodGet) {
		return c.BindQuery(out)
	}
	if strings.EqualFold(c.ContentType(), string(APPLICATION_URLENCODED)) {
		return c.Bind(out)
	}
	return BindJSON(c, out)
}
