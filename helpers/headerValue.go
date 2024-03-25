package helpers

import "github.com/gin-gonic/gin"

func GetContentType(c *gin.Context) string {
	contentType := c.GetHeader("Content-Type")
	return contentType
}

func IsJSONContentType(c *gin.Context) bool {
	return c.ContentType() == "application/json"
}
