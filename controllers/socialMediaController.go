package controllers

import (
	"final-project/database"
	"final-project/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateSocialMedia(c *gin.Context) {
	var socialMedia models.SocialMedia
	if err := c.ShouldBindJSON(&socialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	userData, exists := c.Get("userData")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims := userData.(jwt.MapClaims)
	userID := uint64(claims["id"].(float64))

	socialMedia.UserID = userID

	if err := socialMedia.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&socialMedia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social media"})
		return
	}

	c.JSON(http.StatusCreated, socialMedia)
}

func GetSocialMedias(c *gin.Context) {
	var socialMedias []models.SocialMedia
	if err := database.DB.Find(&socialMedias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social medias"})
		return
	}

	c.JSON(http.StatusOK, socialMedias)
}

func UpdateSocialMedia(c *gin.Context) {
	socialMediaID := c.Param("socialMediaId")
	var socialMedia models.SocialMedia
	if err := database.DB.Where("id = ?", socialMediaID).First(&socialMedia).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
		return
	}

	var updatedSocialMedia models.SocialMedia
	if err := c.ShouldBindJSON(&updatedSocialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := updatedSocialMedia.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&socialMedia).Updates(updatedSocialMedia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social media"})
		return
	}

	c.JSON(http.StatusOK, updatedSocialMedia)
}

func DeleteSocialMedia(c *gin.Context) {
	socialMediaID := c.Param("socialMediaId")
	var socialMedia models.SocialMedia
	if err := database.DB.Where("id = ?", socialMediaID).First(&socialMedia).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
		return
	}

	if err := database.DB.Delete(&socialMedia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete social media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Social media deleted successfully"})
}
