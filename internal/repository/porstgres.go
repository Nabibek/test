package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"Mini-Quicko/internal/core/models"
	"Mini-Quicko/internal/core/ports"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connStr string) (ports.Repository, error) {
	var db *sql.DB
	var err error

	// Пытаемся подключиться с retry
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Failed to open database (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping database (attempt %d): %v", i+1, err)
			db.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		// Успешное подключение
		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return &PostgresRepository{db: db}, nil
}

func createTables(db *sql.DB) error {
	// Таблица истории цен
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS price_history (
			id SERIAL PRIMARY KEY,
			product_id VARCHAR(255) NOT NULL,
			seller_id VARCHAR(255) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(product_id, seller_id, timestamp)
		)
	`)
	if err != nil {
		return err
	}

	// Таблица информации о продуктах
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS product_info (
			product_id VARCHAR(255) NOT NULL,
			seller_id VARCHAR(255) NOT NULL,
			seller_name VARCHAR(255),
			price DECIMAL(10,2) NOT NULL,
			rating DECIMAL(3,2),
			reviews INTEGER,
			purchases INTEGER,
			sku VARCHAR(255),
			segment DECIMAL(3,1),
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (product_id, seller_id, timestamp)
		)
	`)
	return err
}

func (r *PostgresRepository) SavePriceHistory(ctx context.Context, history *models.PriceHistory) error {
	query := `
		INSERT INTO price_history (product_id, seller_id, price, timestamp)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query,
		history.ProductID,
		history.SellerID,
		history.Price,
		history.Timestamp,
	)
	return err
}

func (r *PostgresRepository) GetPriceHistory(ctx context.Context, productID string, limit int) ([]models.PriceHistory, error) {
	query := `
		SELECT id, product_id, seller_id, price, timestamp
		FROM price_history
		WHERE product_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, productID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.PriceHistory
	for rows.Next() {
		var h models.PriceHistory
		if err := rows.Scan(&h.ID, &h.ProductID, &h.SellerID, &h.Price, &h.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}

func (r *PostgresRepository) GetProductInfo(ctx context.Context, productID string) (*models.ProductInfo, error) {
	query := `
		SELECT seller_id, seller_name, price, rating, reviews, purchases, sku, segment, timestamp
		FROM product_info
		WHERE product_id = $1 AND timestamp = (
			SELECT MAX(timestamp) FROM product_info WHERE product_id = $1
		)
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productInfo := &models.ProductInfo{
		ProductID: productID,
		Sellers:   []models.Seller{},
	}

	for rows.Next() {
		var seller models.Seller
		var timestamp time.Time
		if err := rows.Scan(&seller.ID, &seller.Name, &seller.Price, &seller.Rating, &seller.Reviews, &seller.Purchases, &seller.SKU, &seller.Segment, &timestamp); err != nil {
			return nil, err
		}
		productInfo.Sellers = append(productInfo.Sellers, seller)
		productInfo.Timestamp = timestamp // Будет установлено время последней записи
	}

	return productInfo, nil
}

func (r *PostgresRepository) SaveProductInfo(ctx context.Context, productInfo *models.ProductInfo) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, seller := range productInfo.Sellers {
		query := `
			INSERT INTO product_info (product_id, seller_id, seller_name, price, rating, reviews, purchases, sku, segment, timestamp)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`
		_, err := tx.ExecContext(ctx, query,
			productInfo.ProductID,
			seller.ID,
			seller.Name,
			seller.Price,
			seller.Rating,
			seller.Reviews,
			seller.Purchases,
			seller.SKU,
			seller.Segment,
			productInfo.Timestamp,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) HealthCheck(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
