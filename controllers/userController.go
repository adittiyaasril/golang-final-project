package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {

	if !helpers.IsJSONContentType(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
		return
	}

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	if err := database.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "chk_users_age") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Age must be greater than 8"})
			return
		} else if strings.Contains(err.Error(), "email") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email has already been registered"})
			return
		} else if strings.Contains(err.Error(), "username") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username has already been taken"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	responseData := gin.H{
		"age":             user.Age,
		"email":           user.Email,
		"id":              user.ID,
		"profileImageURL": user.ProfileImageURL,
		"username":        user.Username,
	}

	c.JSON(http.StatusCreated, responseData)
}

func Login(c *gin.Context) {
	if !helpers.IsJSONContentType(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
		return
	}

	var loginReq models.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := loginReq.ValidateLogin(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", loginReq.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := helpers.ComparePasswordHash(loginReq.Password, existingUser.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := helpers.GenerateToken(existingUser.ID, existingUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UpdateUser(c *gin.Context) {

	if !helpers.IsJSONContentType(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
		return
	}

	userID := c.Param("userId")

	userData, exists := c.Get("userData")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, ok := userData.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userIDInt64, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if userIDInt64 != int64(claims["id"].(float64)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
		return
	}

	var updateUserRequest models.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := updateUserRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if updateUserRequest.Username != "" {
		user.Username = updateUserRequest.Username
	}
	if updateUserRequest.Email != "" {
		user.Email = updateUserRequest.Email
	}
	if updateUserRequest.Age != 0 {
		user.Age = updateUserRequest.Age
	}
	if updateUserRequest.ProfileImageURL != "" {
		user.ProfileImageURL = updateUserRequest.ProfileImageURL
	}

	if err := database.DB.Save(&user).Error; err != nil {
		if strings.Contains(err.Error(), "chk_users_age") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Age must be greater than 8"})
			return
		} else if strings.Contains(err.Error(), "email") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email has already been registered"})
			return
		} else if strings.Contains(err.Error(), "username") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username has already been taken"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	responseData := gin.H{
		"id":              user.ID,
		"username":        user.Username,
		"email":           user.Email,
		"age":             user.Age,
		"profileImageURL": user.ProfileImageURL,
	}
	c.JSON(http.StatusOK, responseData)
}

func DeleteUser(c *gin.Context) {

	userData, exists := c.Get("userData")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, ok := userData.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := uint64(claims["id"].(float64))

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
