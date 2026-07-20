package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sarowl/GoGinTracker/backend/initializers"
	"github.com/sarowl/GoGinTracker/backend/models"
)

func GetExpenses(c *gin.Context) {
	var expenses []models.Expense
	if err := initializers.DB.Order("expense_date desc").Find(&expenses).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, expenses)
}

func CreateExpense(c *gin.Context) {
	var dto struct {
		ExpenseNo       string  `json:"expenseNo"`
		VendorInvoiceNo string  `json:"vendorInvoiceNo"`
		Category        string  `json:"category"`
		Description     string  `json:"description"`
		Amount          float64 `json:"amount"`
		Vendor          string  `json:"vendor"`
		PaymentMethod   string  `json:"paymentMethod"`
		ReceiptURL      string  `json:"receiptUrl"`
		ExpenseDate     string  `json:"expenseDate"`
		RecordedBy      string  `json:"recordedBy"`
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var expDate time.Time
	if dto.ExpenseDate != "" {
		parsedDate, err := time.Parse("2006-01-02", dto.ExpenseDate)
		if err == nil {
			expDate = parsedDate
		} else {
			expDate = time.Now()
		}
	} else {
		expDate = time.Now()
	}

	expNo := dto.ExpenseNo
	if expNo == "" {
		expNo = fmt.Sprintf("EXP-%d", (time.Now().UnixNano()/1e6)%100000)
	}

	expense := models.Expense{
		ExpenseNo:       expNo,
		VendorInvoiceNo: dto.VendorInvoiceNo,
		Category:        dto.Category,
		Description:     dto.Description,
		Amount:          dto.Amount,
		Vendor:          dto.Vendor,
		PaymentMethod:   dto.PaymentMethod,
		ReceiptURL:      dto.ReceiptURL,
		ExpenseDate:     expDate,
		RecordedBy:      dto.RecordedBy,
	}

	if err := initializers.DB.Create(&expense).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, expense)
}

func DeleteExpense(c *gin.Context) {
	role := c.GetHeader("X-User-Role")
	if role != "Owner" {
		c.JSON(403, gin.H{"error": "Employees are not permitted to delete expenses."})
		return
	}
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.Expense{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("[DEBUG API] Soft-deleted Expense ID %s (marked deleted_at)\n", id)
	c.JSON(200, gin.H{"message": "Expense soft deleted"})
}
