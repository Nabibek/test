package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"Mini-Quicko/internal/core/models"
	"Mini-Quicko/internal/core/ports"
)

type service struct {
	repo ports.Repository
}

func NewService(repo ports.Repository) ports.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) SaveKaspiData(ctx context.Context, request *models.KaspiDataRequest) (*models.ProductAnalysis, error) {
	// Конвертируем офферы в sellers
	sellers := make([]models.Seller, len(request.Offers.Offers))
	for i, offer := range request.Offers.Offers {
		sellers[i] = models.Seller{
			ID:        offer.MerchantId,
			Name:      offer.MerchantName,
			Price:     offer.Price,
			Rating:    offer.MerchantRating,
			Reviews:   offer.MerchantReviewsQuantity,
			Purchases: offer.PurchaseCount,
			SKU:       offer.MerchantSku,
			Segment:   offer.MerchantSegmentId,
		}
	}

	// Сохраняем информацию о продукте
	productInfo := &models.ProductInfo{
		ProductID: request.ProductID,
		Sellers:   sellers,
		Timestamp: time.Now(),
	}

	if err := s.repo.SaveProductInfo(ctx, productInfo); err != nil {
		log.Printf("Warning: failed to save product info: %v", err)
	}

	// Сохраняем историю цен
	for _, seller := range sellers {
		history := &models.PriceHistory{
			ProductID: request.ProductID,
			SellerID:  seller.ID,
			Price:     seller.Price,
			Timestamp: time.Now(),
		}
		if err := s.repo.SavePriceHistory(ctx, history); err != nil {
			log.Printf("Warning: failed to save price history for seller %s: %v", seller.ID, err)
		}
	}

	// Анализируем цены
	analysis := s.analyzePrices(request.ProductID, sellers)
	analysis.TotalOffers = request.Offers.Total
	analysis.AnalysisTime = time.Now().Format(time.RFC3339)

	return analysis, nil
}

func (s *service) AnalyzeProduct(ctx context.Context, productID string) (*models.ProductAnalysis, error) {
	// Получаем последние данные из БД
	productInfo, err := s.repo.GetProductInfo(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product info: %w", err)
	}

	if len(productInfo.Sellers) == 0 {
		return nil, fmt.Errorf("no data found for product %s", productID)
	}

	// Анализируем цены
	analysis := s.analyzePrices(productID, productInfo.Sellers)
	analysis.AnalysisTime = time.Now().Format(time.RFC3339)

	return analysis, nil
}

func (s *service) analyzePrices(productID string, sellers []models.Seller) *models.ProductAnalysis {
	if len(sellers) == 0 {
		return &models.ProductAnalysis{
			ProductID:      productID,
			DumpingSellers: []models.Seller{},
			Sellers:        []models.Seller{},
			AnalysisTime:   time.Now().Format(time.RFC3339),
		}
	}

	// Вычисляем минимальную и среднюю цену
	minPrice := sellers[0].Price
	maxPrice := sellers[0].Price
	sumPrice := 0.0

	// Собираем статистику по сегментам
	segmentStats := make(map[float64]struct {
		count int
		total float64
		min   float64
	})

	for _, seller := range sellers {
		if seller.Price < minPrice {
			minPrice = seller.Price
		}
		if seller.Price > maxPrice {
			maxPrice = seller.Price
		}
		sumPrice += seller.Price

		// Статистика по сегментам
		stat := segmentStats[seller.Segment]
		stat.count++
		stat.total += seller.Price
		if stat.count == 1 || seller.Price < stat.min {
			stat.min = seller.Price
		}
		segmentStats[seller.Segment] = stat
	}

	avgPrice := sumPrice / float64(len(sellers))

	// Определяем демпингующих продавцов
	var dumpingSellers []models.Seller
	for _, seller := range sellers {
		// Проверяем демпинг относительно сегмента
		if stat, exists := segmentStats[seller.Segment]; exists && stat.count > 1 {
			segmentAvg := stat.total / float64(stat.count)
			// Демпинг если цена ниже 90% от средней по сегменту или ниже минимальной по сегменту
			if seller.Price < segmentAvg*0.9 || seller.Price < stat.min*1.05 {
				dumpingSellers = append(dumpingSellers, seller)
			}
		}
	}

	// Вычисляем оптимальную цену
	optimalPrice := s.calculateOptimalPrice(minPrice, avgPrice, sellers, segmentStats)

	return &models.ProductAnalysis{
		ProductID:      productID,
		MinPrice:       minPrice,
		AvgPrice:       avgPrice,
		OptimalPrice:   optimalPrice,
		DumpingSellers: dumpingSellers,
		Sellers:        sellers,
		AnalysisTime:   time.Now().Format(time.RFC3339),
	}
}

func (s *service) calculateOptimalPrice(minPrice, avgPrice float64, sellers []models.Seller, segmentStats map[float64]struct {
	count int
	total float64
	min   float64
},
) float64 {
	// Логика расчета оптимальной цены:
	// 1. Берем среднюю цену между минимальной и средней
	// 2. Учитываем сегмент продавца
	// 3. Добавляем небольшую премию за хороший рейтинг

	basePrice := (minPrice + avgPrice) / 2

	// Находим средний рейтинг всех продавцов
	totalRating := 0.0
	for _, seller := range sellers {
		totalRating += seller.Rating
	}
	avgRating := totalRating / float64(len(sellers))

	// Корректируем цену на основе рейтинга (продавцы с высоким рейтингом могут брать больше)
	ratingMultiplier := 1.0 + (avgRating-4.0)*0.05 // +5% за каждый балл выше 4.0

	optimal := basePrice * ratingMultiplier

	// Округляем до кратного 1000
	return math.Round(optimal/1000) * 1000
}

func (s *service) GetPriceHistory(ctx context.Context, productID string) ([]models.PriceHistory, error) {
	return s.repo.GetPriceHistory(ctx, productID, 100) // последние 100 записей
}

func (s *service) GetProductInfo(ctx context.Context, productID string) (*models.ProductInfo, error) {
	return s.repo.GetProductInfo(ctx, productID)
}

func (s *service) HealthCheck(ctx context.Context) error {
	return s.repo.HealthCheck(ctx)
}
