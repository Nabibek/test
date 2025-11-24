package models

import "time"

type KaspiDataRequest struct {
	ProductID string        `json:"product_id"`
	Offers    ProductOffers `json:"offers"`
}
type ProductOffers struct {
	Offers      []Offer `json:"offers"`
	Total       int     `json:"total"`
	OffersCount int     `json:"offersCount"`
}
type Offer struct {
	MasterSku               string                 `json:"masterSku"`
	MasterCategory          string                 `json:"masterCategory"`
	MerchantId              string                 `json:"merchantId"`
	MerchantName            string                 `json:"merchantName"`
	MerchantSku             string                 `json:"merchantSku"`
	MerchantReviewsQuantity int                    `json:"merchantReviewsQuantity"`
	MerchantRating          float64                `json:"merchantRating"`
	MerchantSegmentId       float64                `json:"merchantSegmentId"`
	PurchaseCount           int                    `json:"purchaseCount"`
	Title                   string                 `json:"title"`
	Price                   float64                `json:"price"`
	PriceBeforeDiscount     float64                `json:"priceBeforeDiscount"`
	Discount                int                    `json:"discount"`
	DeliveryType            string                 `json:"deliveryType"`
	DeliveryDuration        string                 `json:"deliveryDuration"`
	KaspiDelivery           bool                   `json:"kaspiDelivery"`
	Preorder                int                    `json:"preorder"`
	DeliveryOptions         map[string]interface{} `json:"deliveryOptions"`
}
type ProductInfo struct {
	ProductID string    `json:"product_id"`
	Sellers   []Seller  `json:"sellers"`
	Timestamp time.Time `json:"timestamp"`
}
