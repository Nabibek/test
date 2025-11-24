package ports

import (
	"Mini-Quicko/internal/core/models"
	"context"
)

type Repository interface {
	SavePriceHistory(ctx context.Context, history *models.PriceHistory) error
	GetPriceHistory(ctx context.Context, productID string, limit int) ([]models.PriceHistory, error)
	GetProductInfo(ctx context.Context, productID string) (*models.ProductInfo, error)
	SaveProductInfo(ctx context.Context, productInfo *models.ProductInfo) error
	HealthCheck(ctx context.Context) error
	Close() error
}
