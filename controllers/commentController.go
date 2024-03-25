package controllers

import (
	"net/http"
	"strconv"

	"final-project/database"
	"final-project/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
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

	comment.UserID = userID

	if err := comment.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var photo models.Photo
	if err := database.DB.First(&photo, comment.PhotoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func GetComments(c *gin.Context) {
	var comments []models.Comment
	if err := database.DB.Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func UpdateComment(c *gin.Context) {
	commentID := c.Param("commentId")
	id, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var updatedComment models.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if updatedComment.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data provided for update"})
		return
	}

	var existingComment models.Comment
	if err := database.DB.Where("id = ?", id).First(&existingComment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	existingComment.Message = updatedComment.Message

	if err := database.DB.Save(&existingComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, existingComment)
}

func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")
	id, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var comment models.Comment
	if err := database.DB.Where("id = ?", id).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
