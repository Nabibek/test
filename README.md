# Mini-Quicko (Kaspi Price Analyzer)

![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13-336791?logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-3.8-2496ED?logo=docker)

Мини-сервис для анализа цен на товары Kaspi Marketplace. Сервис получает данные о товарах, анализирует цены конкурентов, определяет демпинг и рекомендует оптимальную цену.

## 🚀 Возможности

- **Анализ цен**: Минимальная, средняя и оптимальная цена товара.
- **Выявление демпинга**: Автоматическое определение продавцов с заниженными ценами.
- **История изменений**: Отслеживание изменения цен во времени.
- **REST API**: Простой HTTP API для интеграции.
- **PostgreSQL**: Надежное хранение данных.
- **Docker**: Простой запуск в контейнере.

## 📊 Аналитика

Сервис анализирует:
- Минимальную и среднюю цену товара.
- Демпингующих продавцов (цена ниже 90% от средней по сегменту).
- Оптимальную цену с учетом рейтинга продавца и сегмента.
- Историю изменения цен.

## 🛠 Технологии

- **Backend**: Go 1.19+
- **База данных**: PostgreSQL
- **HTTP Router**: Gorilla Mux
- **Контейнеризация**: Docker + Docker Compose

## 📦 Установка и запуск

### Требования

- Docker и Docker Compose.
- Go 1.19+ (для локальной разработки).

### Быстрый запуск с Docker

1. **Клонируйте репозиторий**:
```bash
git clone <repository-url>
cd Mini-Quicko
```
2. **Запустите сервис**:
```bash
docker-compose up --build
```
3. **Проверьте работу сервиса**:
```bash

curl http://localhost:8080/health
```
### Локальная разработка

1.**Установите зависимости**:
```bash

go mod download
```
2. **Запустите PostgreSQL**:
```bash

docker-compose up db -d
```
3. **Запустите приложение**:
```bash

go run cmd/server/main.go
```
### 📡 API Endpoints

1. **Health Check**

        GET /health

        Проверка работоспособности сервиса.

    Response:

{
  "status": "healthy"
}

2. **Сохранение данных Kaspi**
```http
    POST /products/save-kaspi-data
```
    Сохранение данных о товаре с Kaspi Marketplace.

Request Body:
```json
{
  "product_id": "121806358",
  "offers": {
    "offers": [
      {
        "masterSku": "121806358",
        "merchantId": "30358551",
        "merchantName": "Xiaomi Official Store",
        "merchantRating": 4.9,
        "merchantReviewsQuantity": 3588,
        "purchaseCount": 428,
        "price": 179990.0,
        "priceBeforeDiscount": 199990.0,
        "discount": 10,
        "deliveryType": "TO_DOOR",
        "deliveryDuration": "TILL_2_DAYS"
      }
    ],
    "total": 27,
    "offersCount": 27
  }
}
```
Response:
```json
{
  "product_id": "121806358",
  "min_price": 179990,
  "avg_price": 184990,
  "optimal_price": 182000,
  "dumping_sellers": [
    {
      "id": "30358551",
      "name": "Xiaomi Official Store",
      "price": 179990,
      "rating": 4.9,
      "reviews": 3588,
      "purchases": 428,
      "sku": "55421",
      "segment": 3
    }
  ],
  "sellers": [...],
  "total_offers": 27,
  "analysis_time": "2024-01-15T10:30:00Z"
}
```
3. **Анализ продукта**
```http
    GET /products/{productId}/analyze
```
    Получение анализа цен для сохраненного продукта.

Response: Аналогично POST /products/save-kaspi-data.

4. **История цен**
```http
    GET /products/{productId}/history
```
    Получение истории изменения цен.

Response:
```json
[
  {
    "id": 1,
    "product_id": "121806358",
    "seller_id": "30358551",
    "price": 179990,
    "timestamp": "2024-01-15T10:30:00Z"
  }
]
```
5. **Информация о продукте**
```http
    GET /products/{productId}/info
```
    Получение текущей информации о продукте.

Response:
```json
    {
      "product_id": "121806358",
      "sellers": [...],
      "timestamp": "2024-01-15T10:30:00Z"
    }
```
### 🗄 Структура проекта
```
Mini-Quicko/
├── cmd/                    # Точка входа приложения
│   └── server/
├── config/                 # Конфигурация
├── internal/
│   ├── core/               # Модели данных и порты
│   ├── handlers/           # HTTP обработчики
│   ├── repository/         # Работа с БД (PostgreSQL)
│   └── service/            # Бизнес-логика
├── docker/                 # Docker конфигурации
├── Dockerfile              # Конфигурация Docker образа
├── docker-compose.yml      # Конфигурация Docker Compose
└── README.md               # Этот файл
```
### 🔧 Конфигурация

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| SERVER_PORT | 8080 | Порт HTTP сервера |
| DB_HOST | localhost | Хост PostgreSQL |
| DB_PORT | 5432 | Порт PostgreSQL |
| DB_USER | postgres | Пользователь БД |
| DB_PASSWORD | password | Пароль БД |
| DB_NAME | kaspi_analyzer | Имя базы данных |

### 🐳 Docker
Сборка образа
```bash
docker build -t mini-quicko .
```
**Запуск с Docker Compose**
```bash
docker-compose up -d
```
**Просмотр логов**
```bash
docker-compose logs -f app
```
### 🎯 Пример использования

**Сохранение данных с Kaspi**
```
curl -X POST http://localhost:8080/products/save-kaspi-data \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "121806358",
    "offers": {
        "offers": [
            {
                "masterSku": "121806358",
                "merchantId": "30358551",
                "merchantName": "Xiaomi Official Store",
                "merchantRating": 4.9,
                "merchantReviewsQuantity": 3588,
                "purchaseCount": 428,
                "price": 179990.0,
                "priceBeforeDiscount": 199990.0,
                "discount": 10,
                "deliveryType": "TO_DOOR",
                "deliveryDuration": "TILL_2_DAYS"
            }
        ],
        "total": 27,
        "offersCount": 27
    }
}'
```
### Получение анализа
```
curl http://localhost:8080/products/121806358/analyze
```
### Просмотр истории цен
```
curl http://localhost:8080/products/121806358/history
```
### 🔍 Логика анализа
Определение демпинга

Продавец считается демпингующим, если:

    Цена ниже 90% от средней цены в его сегменте.

    Цена ниже минимальной цены в сегменте + 5%.

**Расчет оптимальной цены**
```bash

optimal_price = (min_price + avg_price) / 2 * rating_multiplier
rating_multiplier = 1.0 + (avg_rating - 4.0) * 0.05
```