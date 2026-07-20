package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sarowl/GoGinTracker/backend/initializers"
	"github.com/sarowl/GoGinTracker/backend/models"
)

func CreateCompany(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		OwnerUID string `json:"ownerUid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	code := fmt.Sprintf("CMP-%06d", (time.Now().UnixNano()/1e3)%1000000)

	comp := models.Company{
		Name:     req.Name,
		Code:     code,
		OwnerUID: req.OwnerUID,
	}

	if err := initializers.DB.Create(&comp).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var user models.Users
	if err := initializers.DB.Where("firebase_uid = ?", req.OwnerUID).First(&user).Error; err == nil {
		user.CompanyCode = code
		initializers.DB.Save(&user)
	}

	c.JSON(201, gin.H{
		"company": comp,
		"user":    user,
	})
}

func JoinCompany(c *gin.Context) {
	var req struct {
		Code    string `json:"code" binding:"required"`
		UserUID string `json:"userUid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var comp models.Company
	if err := initializers.DB.Where("code = ?", req.Code).First(&comp).Error; err != nil {
		c.JSON(404, gin.H{"error": "Invalid Company Code. Please check with your business owner."})
		return
	}

	var user models.Users
	if err := initializers.DB.Where("firebase_uid = ?", req.UserUID).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User profile not found."})
		return
	}

	user.CompanyCode = comp.Code
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"company": comp,
		"user":    user,
	})
}
