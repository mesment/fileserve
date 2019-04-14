package handler

import (
	"time"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mesment/fileserver/utils"
	pkge "github.com/mesment/fileserver/pkg/errors"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = pkge.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = pkge.INVALID_PARAMS
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				code = pkge.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = pkge.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != pkge.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code" : code,
				"msg" : pkge.GetMsg(code),
				"data" : data,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
