package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sarowl/GoGinTracker/backend/initializers"
	"github.com/sarowl/GoGinTracker/backend/models"
)

func GetOrders(c *gin.Context) {
	var orders []models.Order
	query := initializers.DB
	if c.Query("all") != "true" {
		query = query.Where("status != ?", "Delivered")
	}
	if err := query.Find(&orders).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, orders)
}

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		fmt.Printf("[DEBUG API] JSON Binding Error: %v\n", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	order.TotalCost = float64(order.Quantity) * order.Cost
	if order.Customer == "" {
		order.Customer = order.ShipTo
	}
	if order.OrderDate.IsZero() {
		order.OrderDate = time.Now()
	}
	if order.Status == "" {
		order.Status = "Pending"
	}

	if err := initializers.DB.Create(&order).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, order)
}

func DeleteOrder(c *gin.Context) {
	role := c.GetHeader("X-User-Role")
	if role != "Owner" {
		c.JSON(403, gin.H{"error": "Employees are not permitted to delete orders."})
		return
	}
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.Order{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("[DEBUG API] Soft-deleted Order ID %s (marked deleted_at)\n", id)
	c.JSON(200, gin.H{"message": "Order soft deleted"})
}

func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status    string `json:"status" binding:"required"`
		InvoiceID string `json:"invoiceId"`
		UpdatedBy string `json:"updatedBy"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[DEBUG API] Status Update Bind Error: %v\n", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := initializers.DB.First(&order, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	order.Status = req.Status

	if req.Status == "Delivered" {
		invoiceID := req.InvoiceID
		if invoiceID == "" {
			if order.InvoiceID != nil && *order.InvoiceID != "" {
				invoiceID = *order.InvoiceID
			} else {
				invoiceID = fmt.Sprintf("S-%d", (time.Now().UnixNano()/1e6)%100000)
			}
		}
		order.InvoiceID = &invoiceID

		now := time.Now()
		order.DeliveredAt = &now

		if req.UpdatedBy != "" {
			order.DeliveredBy = &req.UpdatedBy
		}

		saleRecord := models.Sale{
			InvoiceID:   invoiceID,
			RefPO:       order.PONumber,
			Customer:    order.Customer,
			Description: order.Description,
			Amount:      order.TotalCost,
			Date:        now.Format("2006-01-02"),
			DeliveredBy: req.UpdatedBy,
		}
		initializers.DB.Where("invoice_id = ?", invoiceID).FirstOrCreate(&saleRecord)
	}

	if err := initializers.DB.Save(&order).Error; err != nil {
		fmt.Printf("[DEBUG API] DB Save Error: %v\n", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[DEBUG API] Order #%d status successfully updated to '%s' (Invoice: %v)\n", order.PONumber, req.Status, req.InvoiceID)
	c.JSON(200, order)
}
