package models

import "time"

type PriceHistory struct {
	ID        int       `json:"id" db:"id"`
	ProductID string    `json:"product_id" db:"product_id"`
	SellerID  string    `json:"seller_id" db:"seller_id"`
	Price     float64   `json:"price" db:"price"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}
