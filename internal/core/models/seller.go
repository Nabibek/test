package models

type Seller struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Rating    float64 `json:"rating"`
	Reviews   int     `json:"reviews"`
	Purchases int     `json:"purchases"`
	SKU       string  `json:"sku"`
	Segment   float64 `json:"segment"`
}
