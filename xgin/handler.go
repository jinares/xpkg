package xgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
	"net/http"
)

type (
	HandlerFunc func(c *gin.Context) (ret IRet)
)

// DoHandler 返回JSON数据
func DoHandler(handler ...HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		if len(handler) < 1 {

			c.JSON(http.StatusOK, Error(c, xtools.XErr(codes.NotFound, "404")))
			return
		}
		for _, item := range handler {
			cb, err := CopyBody(c)
			if err != nil {
				fmt.Println(fmt.Sprintf("copy-gin-body-err:%s", err.Error()))
				return
			}
			//如果返回结果为空 则执行下一个Handler
			ret := item(c)
			if ret != nil {
				ret = Trace(c, ret)
				data, contenttype := ret.GetRet()
				c.Data(http.StatusOK, string(contenttype), []byte(data))
				return
			}
			//如果 context 已经终止 退出
			if c.IsAborted() {
				break
			}
			c.Request.Body = cb
		}

		c.JSON(http.StatusOK, Succ(c, true))

	}
}
func BaseHandler(handler ...HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		for _, item := range handler {
			cb, err := CopyBody(c)
			if err != nil {
				fmt.Println(fmt.Sprintf("copy-gin-body-err:%s", err.Error()))
				return
			}
			ret := item(c)
			//如果返回结果为空 则执行下一个Handler
			if ret != nil {
				ret = Trace(c, ret)
				data, contenttype := ret.GetRet()
				c.Data(http.StatusOK, string(contenttype), []byte(data))
				c.Abort()
				return
			}
			//如果 context 已经终止 退出
			if c.IsAborted() {
				break
			}
			c.Request.Body = cb

		}

	}
}
