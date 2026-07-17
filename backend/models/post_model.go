package models

import (
	"time"
)

type Users struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	Email       string     `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	FirebaseUID string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"firebase_uid"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Users) TableName() string {
	return "users"
}

type Order struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	PONumber    uint64     `gorm:"uniqueIndex;not null" json:"poNumber"`
	ShipTo      string     `gorm:"type:varchar(255);not null" json:"shipTo"`
	Description string     `gorm:"type:text" json:"description"`
	Quantity    int        `gorm:"type:integer;not null;default:1" json:"quantity"`
	Cost        float64    `gorm:"type:decimal(12,2);not null" json:"cost"`
	TotalCost   float64    `gorm:"type:decimal(12,2);not null" json:"totalCost"`
	CancelDate  string     `gorm:"type:varchar(50);not null" json:"cancelDate"`
	Status      string     `gorm:"type:varchar(50);not null;default:'Pending'" json:"status"`
	InsertedBy  string     `gorm:"type:varchar(255)" json:"insertedBy"` // Firebase UID of creator
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type Sale struct {
	ID          uint       `gorm:"primaryKey" json:"dbId"`
	InvoiceID   string     `gorm:"uniqueIndex;type:varchar(100);not null" json:"id"`
	RefPO       uint64     `gorm:"not null" json:"refPo"`
	Customer    string     `gorm:"type:varchar(255);not null" json:"customer"`
	Description string     `gorm:"type:text" json:"description"`
	Amount      float64    `gorm:"type:decimal(12,2);not null" json:"amount"`
	Date        string     `gorm:"type:varchar(50);not null" json:"date"`
	DeliveredBy string     `gorm:"type:varchar(255)" json:"deliveredBy"` // Firebase UID of delivery marker
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Sale) TableName() string {
	return "sales"
}
