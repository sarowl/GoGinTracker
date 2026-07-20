package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sarowl/GoGinTracker/backend/initializers"
	"github.com/sarowl/GoGinTracker/backend/models"
)

func SyncUser(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.FirebaseUID == "" {
		c.JSON(400, gin.H{"error": "firebase_uid is required"})
		return
	}

	var existing models.Users
	err := initializers.DB.Where("firebase_uid = ?", user.FirebaseUID).First(&existing).Error
	if err == nil {
		if user.Name != "" {
			existing.Name = user.Name
		}
		if user.Email != "" {
			existing.Email = user.Email
		}
		if user.Role != "" {
			existing.Role = user.Role
		}
		if err := initializers.DB.Save(&existing).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, existing)
		return
	}

	if user.Role == "" {
		user.Role = "Employee"
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

func GetUserProfile(c *gin.Context) {
	uid := c.Param("uid")
	var user models.Users
	if err := initializers.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User profile not found"})
		return
	}
	c.JSON(200, user)
}
