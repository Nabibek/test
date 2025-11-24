package ports

import (
	"Mini-Quicko/internal/core/models"
	"context"
)

type Service interface {
	AnalyzeProduct(ctx context.Context, productID string) (*models.ProductAnalysis, error)
	GetPriceHistory(ctx context.Context, productID string) ([]models.PriceHistory, error)
	GetProductInfo(ctx context.Context, productID string) (*models.ProductInfo, error)
	SaveKaspiData(ctx context.Context, request *models.KaspiDataRequest) (*models.ProductAnalysis, error)
	HealthCheck(ctx context.Context) error
}
