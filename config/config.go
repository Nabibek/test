package config

import (
	"os"

	"github.com/spf13/viper"
)

// Структура для хранения конфигурации
type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Функция для загрузки конфигурации
func Load() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./cmd/server")
	viper.AddConfigPath("/app")
	viper.SetConfigType("yaml")

	// Чтение конфигурации
	viper.ReadInConfig()

	// Переменные окружения
	viper.AutomaticEnv() // Автоматическое считывание переменных окружения

	return &Config{
		ServerPort: getConfigValue("server.port", "8080"),
		DBHost:     getConfigValue("db.host", "db"),
		DBPort:     getConfigValue("db.port", "5432"),
		DBUser:     getConfigValue("db.user", "postgres"),
		DBPassword: getConfigValue("db.password", "password"),
		DBName:     getConfigValue("db.name", "kaspi_analyzer"),
	}
}

// Функция для получения значения из конфигурации с возможностью использования переменных окружения
func getConfigValue(key, defaultValue string) string {
	// Пытаемся получить значение из конфигурации
	if value := viper.GetString(key); value != "" {
		return value
	}

	// Если в конфиге нет значения, пытаемся считать из переменной окружения
	if value := os.Getenv(key); value != "" {
		return value
	}

	// Возвращаем значение по умолчанию
	return defaultValue
}
