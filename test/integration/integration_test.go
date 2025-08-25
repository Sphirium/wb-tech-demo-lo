package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sphirium/wb-tech-demo-lo/internal/cache"
	"github.com/Sphirium/wb-tech-demo-lo/internal/handler"
	"github.com/Sphirium/wb-tech-demo-lo/internal/models"
	"github.com/Sphirium/wb-tech-demo-lo/internal/repository"
	"github.com/Sphirium/wb-tech-demo-lo/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestServer создаёт тестовый HTTP-сервер с реальными зависимостями
func setupTestServer() (*http.ServeMux, *gorm.DB, func()) {
	// Подключаемся к тестовой БД
	dsn := "host=localhost user=wbuser password=wbpass dbname=wb_orders port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database for testing: " + err.Error())
	}

	// Очищаем и мигрируем
	db.Migrator().DropTable(&models.Order{}, &models.Item{})
	db.AutoMigrate(&models.Order{}, &models.Item{})

	repo := repository.NewOrderRepository(db)
	cache := cache.NewOrderCache("localhost:6379", "")
	orderService := service.NewOrderService(repo, cache)

	// Создаём роутер с обработчиками
	r := chi.NewRouter()
	orderHandler := handler.NewOrderHandler(orderService)

	// Эндпоинт для приёма заказа (аналог Kafka)
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err := orderService.SaveOrder(body); err != nil {
			http.Error(w, "Invalid order data", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status": "created"}`))
	})

	// Эндпоинт для получения заказа
	r.Get("/order/{order_uid}", orderHandler.GetOrder)

	// Обёртка в http.ServeMux
	mux := http.NewServeMux()
	mux.Handle("/", r)

	// Функция очистки
	teardown := func() {
		cache.Close()
		// Очистка таблиц
		db.Migrator().DropTable(&models.Order{}, &models.Item{})
	}

	return mux, db, teardown
}

// Тест: создание заказа и его получение
func TestCreateAndRetrieveOrder(t *testing.T) {
	// Настройка
	mux, _, teardown := setupTestServer()
	defer teardown()

	server := httptest.NewServer(mux)
	defer server.Close()

	const testOrderUID = "b76eaf44-b342-4b1a-a8bf-ab1c4df1918e"

	testOrder := models.Order{
		OrderUID:    testOrderUID,
		TrackNumber: "TRACK123456",
		Entry:       "WBIL",
		Delivery: &models.Delivery{
			OrderID: testOrderUID,
			Name:    "John Doe",
			Phone:   "+1234567890",
			Zip:     "12345",
			City:    "New York",
			Address: "5th Ave 10",
			Region:  "NY",
			Email:   "john@example.com",
		},
		Payment: &models.Payment{
			Transaction:  testOrderUID,
			OrderID:      testOrderUID,
			Currency:     "USD",
			Provider:     "stripe",
			Amount:       1000,
			DeliveryCost: 200,
			GoodsTotal:   800,
			PaymentDt:    1637907727,
			Bank:         "alpha",
		},
		Items: []models.Item{
			{
				ChrtID:      1001,
				OrderID:     testOrderUID,
				Name:        "Test Item",
				Price:       800,
				TotalPrice:  800,
				NMID:        98765,
				Brand:       "TestBrand",
				Sale:        0,
				Size:        "M",
				Status:      202,
				TrackNumber: "TRACK123456",
				RID:         "a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8",
			},
		},
		Locale:          "en",
		CustomerID:      "cust_001",
		DeliveryService: "meest",
		DateCreated:     "2023-01-01T10:00:00Z",
		Shardkey:        "1",
		SMID:            1,
		OofShard:        "1",
	}

	jsonData, _ := json.Marshal(testOrder)

	// Этап 1: Отправляем заказ через /create
	resp, err := http.Post(server.URL+"/create", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Этап 2: Получаем заказ по ID
	resp, err = http.Get(server.URL + "/order/" + testOrderUID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Читаем тело ответа
	body, _ := io.ReadAll(resp.Body)
	var retrievedOrder models.Order
	err = json.Unmarshal(body, &retrievedOrder)
	assert.NoError(t, err)

	// Проверяем ключевые поля
	assert.Equal(t, testOrder.OrderUID, retrievedOrder.OrderUID)
	assert.Equal(t, testOrder.Delivery.Name, retrievedOrder.Delivery.Name)
	assert.Equal(t, testOrder.Payment.Amount, retrievedOrder.Payment.Amount)
	assert.Equal(t, 1, len(retrievedOrder.Items))
	assert.Equal(t, "Test Item", retrievedOrder.Items[0].Name)

	// Проверка кеширования: второй запрос должен быть быстрее
	start := time.Now()
	resp, err = http.Get(server.URL + "/order/" + testOrderUID)
	duration := time.Since(start)
	assert.Less(t, duration, 100*time.Millisecond, "Second request should be fast (from cache)")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Тест: заказ не найден
func TestOrderNotFound(t *testing.T) {
	mux, _, teardown := setupTestServer()
	defer teardown()

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/order/unknown-order")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// Тест: невалидный JSON
func TestInvalidJSON(t *testing.T) {
	mux, _, teardown := setupTestServer()
	defer teardown()

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Post(server.URL+"/create", "application/json", bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Тест: пустой order_uid — генерируется UUID
func TestEmptyOrderUIDGeneratesUUID(t *testing.T) {
	mux, db, teardown := setupTestServer()
	defer teardown()

	server := httptest.NewServer(mux)
	defer server.Close()

	testOrder := map[string]interface{}{
		"order_uid":    "", // ✅ пустой
		"track_number": "TRACK999",
		"delivery": map[string]string{
			"name": "Test",
		},
		"payment": map[string]interface{}{
			"amount": 500,
		},
		"items": []map[string]interface{}{
			{
				"chrt_id":     1,
				"name":        "Dummy",
				"price":       500,
				"total_price": 500,
			},
		},
		"locale": "en",
	}

	jsonData, _ := json.Marshal(testOrder)
	resp, err := http.Post(server.URL+"/create", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Извлекаем созданный заказ из БД
	var order models.Order
	err = db.Where("track_number = ?", "TRACK999").Preload("Items").First(&order).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, order.OrderUID)
	assert.NotEqual(t, uuid.Nil.String(), order.OrderUID)
	fmt.Printf("Generated OrderUID: %s\n", order.OrderUID)
}
