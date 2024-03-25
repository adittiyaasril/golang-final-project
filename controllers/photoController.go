package controllers

import (
	"final-project/database"
	"final-project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreatePhoto(c *gin.Context) {
	var photo models.Photo
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := photo.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, _ := c.Get("userData")
	claims := userData.(jwt.MapClaims)
	userId := uint64(claims["id"].(float64))
	photo.UserID = userId

	if err := database.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create photo"})
		return
	}

	c.JSON(http.StatusCreated, photo)
}

func GetPhotos(c *gin.Context) {
	var photos []models.Photo
	if err := database.DB.Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch photos"})
		return
	}

	c.JSON(http.StatusOK, photos)
}

func GetPhotoByID(c *gin.Context) {
	photoID := c.Param("photoId")

	id, err := strconv.ParseUint(photoID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	var photo models.Photo
	if err := database.DB.Where("id = ?", id).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	c.JSON(http.StatusOK, photo)
}

func UpdatePhoto(c *gin.Context) {
	photoID := c.Param("photoId")
	id, err := strconv.ParseUint(photoID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	var updatedPhoto models.Photo
	if err := c.ShouldBindJSON(&updatedPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	var existingPhoto models.Photo
	if err := database.DB.Where("id = ?", id).First(&existingPhoto).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	userData, _ := c.Get("userData")
	claims := userData.(jwt.MapClaims)
	userId := uint64(claims["id"].(float64))
	if existingPhoto.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this photo"})
		return
	}

	if err := existingPhoto.ValidateUpdate(updatedPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedPhoto.Title != "" {
		existingPhoto.Title = updatedPhoto.Title
	}
	if updatedPhoto.Caption != "" {
		existingPhoto.Caption = updatedPhoto.Caption
	}
	if updatedPhoto.PhotoURL != "" {
		existingPhoto.PhotoURL = updatedPhoto.PhotoURL
	}

	if err := database.DB.Save(&existingPhoto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo"})
		return
	}

	c.JSON(http.StatusOK, existingPhoto)
}

func DeletePhoto(c *gin.Context) {
	photoID := c.Param("photoId")
	id, err := strconv.ParseUint(photoID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	var photo models.Photo
	if err := database.DB.Where("id = ?", id).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	userData, _ := c.Get("userData")
	claims := userData.(jwt.MapClaims)
	userId := uint64(claims["id"].(float64))
	if photo.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this photo"})
		return
	}

	if err := database.DB.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}
