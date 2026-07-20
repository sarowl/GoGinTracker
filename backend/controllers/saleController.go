package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sarowl/GoGinTracker/backend/initializers"
	"github.com/sarowl/GoGinTracker/backend/models"
)

type SaleDTO struct {
	DBID        uint    `json:"dbId"`
	ID          string  `json:"id"`
	RefPO       uint64  `json:"refPo"`
	Customer    string  `json:"customer"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	DeliveredBy string  `json:"deliveredBy"`
}

func GetSales(c *gin.Context) {
	var deliveredOrders []models.Order
	initializers.DB.Where("status = ?", "Delivered").Find(&deliveredOrders)

	sales := make([]SaleDTO, 0, len(deliveredOrders))
	for _, ord := range deliveredOrders {
		invID := ""
		if ord.InvoiceID != nil && *ord.InvoiceID != "" {
			invID = *ord.InvoiceID
		} else {
			invID = fmt.Sprintf("S-%d", ord.ID)
		}

		cust := ord.Customer
		if cust == "" {
			cust = ord.ShipTo
		}

		delBy := ""
		if ord.DeliveredBy != nil {
			delBy = *ord.DeliveredBy
		}

		dtStr := ord.UpdatedAt.Format("2006-01-02")
		if ord.DeliveredAt != nil {
			dtStr = ord.DeliveredAt.Format("2006-01-02")
		}

		sales = append(sales, SaleDTO{
			DBID:        ord.ID,
			ID:          invID,
			RefPO:       ord.PONumber,
			Customer:    cust,
			Description: ord.Description,
			Amount:      ord.TotalCost,
			Date:        dtStr,
			DeliveredBy: delBy,
		})
	}

	var legacySales []models.Sale
	initializers.DB.Find(&legacySales)
	for _, s := range legacySales {
		exists := false
		for _, existing := range sales {
			if existing.ID == s.InvoiceID {
				exists = true
				break
			}
		}
		if !exists {
			sales = append(sales, SaleDTO{
				DBID:        s.ID,
				ID:          s.InvoiceID,
				RefPO:       s.RefPO,
				Customer:    s.Customer,
				Description: s.Description,
				Amount:      s.Amount,
				Date:        s.Date,
				DeliveredBy: s.DeliveredBy,
			})
		}
	}

	c.JSON(200, sales)
}
