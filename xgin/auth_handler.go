package xgin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinares/xpkg/xoidc"
)

func AuthHandler(handler xoidc.OidcAuthFuncHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		mds := map[string]string{
			xoidc.AUTHORIZATION: c.GetHeader(xoidc.AUTHORIZATION),
		}
		_, err := xoidc.Authorize(c.Request.Context(), c.Request.URL.Path, mds, handler)
		if err != nil {
			ret := NewRet(err, nil)
			result, ctype := ret.GetRet()
			c.Data(200, string(ctype), []byte(result))
			c.Abort()
		}

	}
}
