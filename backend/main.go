package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sarowl/GoGinTracker/backend/controllers"
	"github.com/sarowl/GoGinTracker/backend/initializers"
	"github.com/sarowl/GoGinTracker/backend/middleware"
	"github.com/sarowl/GoGinTracker/backend/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	err := initializers.DB.AutoMigrate(&models.Users{}, &models.Order{}, &models.Sale{}, &models.Expense{}, &models.Company{})
	if err != nil {
		fmt.Printf("[MIGRATION ERROR] Failed to run GORM migration: %v\n", err)
	} else {
		fmt.Println("[MIGRATION] GORM migrations completed successfully!")
	}
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong-pong1234",
		})
	})

	// User Routes
	r.POST("/api/users", controllers.SyncUser)
	r.GET("/api/users/:uid", controllers.GetUserProfile)

	// Order Routes
	r.GET("/api/orders", controllers.GetOrders)
	r.POST("/api/orders", controllers.CreateOrder)
	r.DELETE("/api/orders/:id", controllers.DeleteOrder)
	r.PUT("/api/orders/:id/status", controllers.UpdateOrderStatus)

	// Sales Routes
	r.GET("/api/sales", controllers.GetSales)

	// Expense Routes
	r.GET("/api/expenses", controllers.GetExpenses)
	r.POST("/api/expenses", controllers.CreateExpense)
	r.DELETE("/api/expenses/:id", controllers.DeleteExpense)

	// Company Routes
	r.POST("/api/companies/create", controllers.CreateCompany)
	r.POST("/api/companies/join", controllers.JoinCompany)

	r.Run()
}
