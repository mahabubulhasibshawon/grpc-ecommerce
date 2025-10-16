// internal/domain/models.go
package domain

import "time"

type User struct {
	ID       int64
	Username string
	Password string
}

type Order struct {
	ConsignmentID     string
	CreatedAt         time.Time
	Description       string
	MerchantOrderID   string
	RecipientName     string
	RecipientAddress  string
	RecipientPhone    string
	OrderAmount       float64
	TotalFee          float64
	Instruction       string
	OrderTypeID       int64
	CODFee            float64
	PromoDiscount     float64
	Discount          float64
	DeliveryFee       float64
	Status            string
	OrderType         string
	ItemType          int64
	StoreName         string
	StoreContactPhone string
	CODAmount         float64
	DeliveryCharge    float64
	UserID            int64
	StoreID           int64
	RecipientCity     int64
	RecipientZone     int64
	RecipientArea     int64
	DeliveryType      int64
	ItemQuantity      int64
	ItemWeight        float64
	AmountToCollect   float64
}