package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//func BasicAuth() gin.HandlerFunc {
//	return gin.BasicAuth(gin.Accounts{
//		"key": "value",
//	})
//}

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, password, hasAuth := c.Request.BasicAuth()
		if !hasAuth || user != "key" || password != "value" {
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "UNAUTHORIZED",
			})
			return
		}
		c.Next()
	}
}
