package main

import (
	"Mini-Quicko/config"
	"Mini-Quicko/internal/handlers"
	"Mini-Quicko/internal/repository"
	"Mini-Quicko/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к БД
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	repo, err := repository.NewPostgresRepository(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()

	// Инициализация сервиса
	service := service.NewService(repo)

	// Инициализация handlers
	handler := handlers.NewHTTPHandler(service)

	// Настройка роутера
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Health check для Docker
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := service.HealthCheck(r.Context()); err != nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
