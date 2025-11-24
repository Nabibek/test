package models

type ProductAnalysis struct {
	ProductID      string   `json:"product_id"`
	MinPrice       float64  `json:"min_price"`
	AvgPrice       float64  `json:"avg_price"`
	OptimalPrice   float64  `json:"optimal_price"`
	DumpingSellers []Seller `json:"dumping_sellers"`
	Sellers        []Seller `json:"sellers"`
	TotalOffers    int      `json:"total_offers"`
	AnalysisTime   string   `json:"analysis_time"`
}
