package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sphirium/wb-tech-demo-lo/internal/cache"
	"github.com/Sphirium/wb-tech-demo-lo/internal/config"
	"github.com/Sphirium/wb-tech-demo-lo/internal/consumer"
	"github.com/Sphirium/wb-tech-demo-lo/internal/handler"
	"github.com/Sphirium/wb-tech-demo-lo/internal/repository"
	"github.com/Sphirium/wb-tech-demo-lo/internal/service"
	"github.com/Sphirium/wb-tech-demo-lo/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	logg := logger.New(cfg.LogLevel)

	// Подключение к БД
	db, err := gorm.Open(postgres.Open(cfg.PostgresURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Применяем миграции
	if err := applyMigrations(db, "./migrations"); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	// Создаём зависимости
	repo := repository.NewOrderRepository(db)
	cache := cache.NewOrderCache(cfg.RedisAddr, cfg.RedisPassword)
	serv := service.NewOrderService(repo, cache)

	// Восстанавливаем кеш
	if err := serv.RestoreCacheFromDB(); err != nil {
		logg.Warn("Failed to restore cache from DB: %v", err)
	}

	// Kafka consumer
	consumer := consumer.NewKafkaConsumer(cfg.KafkaBroker, cfg.KafkaTopic, serv)
	consumer.Start()
	defer consumer.Close()

	// HTTP
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	handler := handler.NewOrderHandler(serv)
	r.Get("/order/{order_uid}", handler.GetOrder)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	logg.Info("Server starting on port %s", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, r))
}

// applyMigrations выполняет все .up.sql миграции из папки
func applyMigrations(db *gorm.DB, migrationDir string) error {
	files, err := filepath.Glob(filepath.Join(migrationDir, "*.up.sql"))
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file, err)
		}

		sql := string(content)
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", file, err)
		}

		log.Printf("✅ Applied migration: %s", file)
	}

	return nil
}
