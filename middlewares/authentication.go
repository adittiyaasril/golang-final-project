package middlewares

import (
	"final-project/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := helpers.VerifyToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthenticated",
				"message": "Invalid token",
			})
			return
		}

		c.Set("userData", claims)
		c.Next()
	}
}
