package model

import "time"

type Checkout struct {
	Email    string            `json:"email"`
	Address  string            `json:"address"`
	Products []ProductQuantity `json:"products"`
}

type ProductQuantity struct {
	ID       string `json:"id"`
	Quantity int32  `json:"quantity"`
}

type Order struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	Address           string     `json:"address"`
	GrandTotal        int64      `json:"grandTotal"`
	Passcode          *string    `json:"passcode,omitempty"`
	PaidAt            *time.Time `json:"paidAt,omitempty"`
	PaidBank          *string    `json:"paidBank,omitempty"`
	PaidAccountNumber *string    `json:"paidAccountNumber,omitempty"`
}

type OrderDetail struct {
	ID        string `json:"id"`
	OrderID   string `json:"orderId"`
	ProductID string `json:"productId"`
	Quantity  int32  `json:"quantity"`
	Price     int64  `json:"price"`
	Total     int64  `json:"total"`
}
