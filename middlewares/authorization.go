package middlewares

import (
	"final-project/database"
	"final-project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		photoID := c.Param("photoId")
		id, err := strconv.ParseUint(photoID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
			c.Abort()
			return
		}

		var photo models.Photo
		if err := database.DB.Where("id = ?", id).First(&photo).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			c.Abort()
			return
		}

		userData, _ := c.Get("userData")
		claims := userData.(jwt.MapClaims)
		userId := uint64(claims["id"].(float64))
		if photo.UserID != userId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthorizeComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID := c.Param("commentId")
		id, err := strconv.ParseUint(commentID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			c.Abort()
			return
		}

		var comment models.Comment
		if err := database.DB.Where("id = ?", id).First(&comment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			c.Abort()
			return
		}

		userData, exists := c.Get("userData")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		claims := userData.(jwt.MapClaims)
		userId := uint64(claims["id"].(float64))

		if comment.UserID != userId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthorizeSocialMedia() gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, exists := c.Get("userData")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims := userData.(jwt.MapClaims)
		userID := uint64(claims["id"].(float64))

		socialMediaID := c.Param("socialMediaId")
		socialMediaIDUint, err := strconv.ParseUint(socialMediaID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid social media ID"})
			c.Abort()
			return
		}

		var socialMedia models.SocialMedia
		if err := database.DB.Where("id = ?", socialMediaIDUint).First(&socialMedia).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
			c.Abort()
			return
		}

		if socialMedia.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}
