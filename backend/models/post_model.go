package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255)" json:"name"`
	Email       string         `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	FirebaseUID string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"firebase_uid"`
	Role        string         `gorm:"type:varchar(50);not null;default:'Employee'" json:"role"`
	CompanyCode string         `gorm:"type:varchar(50)" json:"company_code,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Users) TableName() string {
	return "users"
}

type Company struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Code      string         `gorm:"uniqueIndex;type:varchar(50);not null" json:"code"`
	OwnerUID  string         `gorm:"type:varchar(255);not null" json:"ownerUid"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (Company) TableName() string {
	return "companies"
}

type Order struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PONumber    uint64         `gorm:"uniqueIndex;not null" json:"poNumber"`
	InvoiceID   *string        `gorm:"uniqueIndex;type:varchar(100)" json:"invoiceId,omitempty"`

	Customer    string         `gorm:"type:varchar(255)" json:"customer"`
	ShipTo      string         `gorm:"type:varchar(255)" json:"shipTo"`
	Description string         `gorm:"type:text" json:"description"`
	Quantity    int            `gorm:"type:integer;not null;default:1" json:"quantity"`
	Cost        float64        `gorm:"type:decimal(12,2);not null" json:"cost"`
	TotalCost   float64        `gorm:"type:decimal(12,2);not null" json:"totalCost"`

	Status      string         `gorm:"type:varchar(50);not null;default:'Pending'" json:"status"`

	OrderDate   time.Time      `gorm:"type:date" json:"orderDate"`
	CancelDate  string         `gorm:"type:varchar(50)" json:"cancelDate,omitempty"`
	DeliveredAt *time.Time     `json:"deliveredAt,omitempty"`
	PaidAt      *time.Time     `json:"paidAt,omitempty"`

	InsertedBy  string         `gorm:"type:varchar(255)" json:"insertedBy"`
	DeliveredBy *string        `gorm:"type:varchar(255)" json:"deliveredBy,omitempty"`

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type Sale struct {
	ID          uint           `gorm:"primaryKey" json:"dbId"`
	InvoiceID   string         `gorm:"uniqueIndex;type:varchar(100);not null" json:"id"`
	RefPO       uint64         `gorm:"not null" json:"refPo"`
	Customer    string         `gorm:"type:varchar(255);not null" json:"customer"`
	Description string         `gorm:"type:text" json:"description"`
	Amount      float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Date        string         `gorm:"type:varchar(50);not null" json:"date"`
	DeliveredBy string         `gorm:"type:varchar(255)" json:"deliveredBy"` // Firebase UID of delivery marker
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Sale) TableName() string {
	return "sales"
}

type Expense struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ExpenseNo       string         `gorm:"uniqueIndex;type:varchar(50);not null" json:"expenseNo"`
	VendorInvoiceNo string         `gorm:"type:varchar(100)" json:"vendorInvoiceNo,omitempty"`
	Category        string         `gorm:"type:varchar(100);not null" json:"category"`
	Description     string         `gorm:"type:text" json:"description"`
	Amount          float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Vendor          string         `gorm:"type:varchar(255)" json:"vendor,omitempty"`
	PaymentMethod   string         `gorm:"type:varchar(50)" json:"paymentMethod,omitempty"`
	ReceiptURL      string         `gorm:"type:text" json:"receiptUrl,omitempty"`
	ExpenseDate     time.Time      `gorm:"type:date;not null" json:"expenseDate"`
	RecordedBy      string         `gorm:"type:varchar(255);not null" json:"recordedBy"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (Expense) TableName() string {
	return "expenses"
}
