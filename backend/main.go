package main

import (
	"fmt"
	"time"

	"github.com/sarowl/GoGinTracker/backend/initializers"
	"github.com/sarowl/GoGinTracker/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	err := initializers.DB.AutoMigrate(&models.Users{}, &models.Order{}, &models.Sale{})
	if err != nil {
		fmt.Printf("[MIGRATION ERROR] Failed to run GORM migration: %v\n", err)
	} else {
		fmt.Println("[MIGRATION] GORM migrations completed successfully!")
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong-pong1234",
		})
	})

	r.POST("/api/users", func(c *gin.Context) {
		var user models.Users
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if user.FirebaseUID == "" {
			c.JSON(400, gin.H{"error": "firebase_uid is required"})
			return
		}

		result := initializers.DB.Where(models.Users{FirebaseUID: user.FirebaseUID}).
			Assign(models.Users{Name: user.Name, Email: user.Email}).
			FirstOrCreate(&user)

		if result.Error != nil {
			c.JSON(500, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(201, user)
	})

	// --- Orders API ---
	r.GET("/api/orders", func(c *gin.Context) {
		var orders []models.Order
		if err := initializers.DB.Find(&orders).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, orders)
	})

	r.POST("/api/orders", func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			fmt.Printf("[DEBUG API] JSON Binding Error: %v\n", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		order.TotalCost = float64(order.Quantity) * order.Cost
		if err := initializers.DB.Create(&order).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, order)
	})

	r.DELETE("/api/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := initializers.DB.Delete(&models.Order{}, id).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Order deleted"})
	})

	r.PUT("/api/orders/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Status    string `json:"status" binding:"required"`
			InvoiceID string `json:"invoiceId"`
			UpdatedBy string `json:"updatedBy"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if req.Status == "Delivered" {
			err := initializers.DB.Transaction(func(tx *gorm.DB) error {
				var order models.Order
				if err := tx.First(&order, id).Error; err != nil {
					return err
				}

				invoiceID := req.InvoiceID
				if invoiceID == "" {
					invoiceID = fmt.Sprintf("S-%d", (time.Now().UnixNano()/1e6)%100000)
				}

				sale := models.Sale{
					InvoiceID:   invoiceID,
					RefPO:       order.PONumber,
					Customer:    order.ShipTo,
					Description: order.Description,
					Amount:      order.TotalCost,
					Date:        time.Now().Format("2006-01-02"),
					DeliveredBy: req.UpdatedBy,
				}

				if err := tx.Create(&sale).Error; err != nil {
					return err
				}

				if err := tx.Delete(&order).Error; err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "Delivered", "transferred": true})
			return
		}

		var order models.Order
		if err := initializers.DB.First(&order, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}
		order.Status = req.Status
		if err := initializers.DB.Save(&order).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, order)
	})

	// --- Sales API ---
	r.GET("/api/sales", func(c *gin.Context) {
		var sales []models.Sale
		if err := initializers.DB.Find(&sales).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, sales)
	})

	r.Run()
}
