package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Sphirium/learning-projects/wb-tech-demo-lo/internal/models"
	"github.com/go-redis/redis/v8"
)

type OrderCacheInterface interface {
	Get(orderUID string) (*models.Order, error)
	Set(order *models.Order) error
	Close() error
}

type OrderCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewOrderCache(addr, password string) OrderCacheInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &OrderCache{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *OrderCache) Set(order *models.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, "order:"+order.OrderUID, data, 24*time.Hour).Err()
}

func (c *OrderCache) Get(orderUID string) (*models.Order, error) {
	data, err := c.client.Get(c.ctx, "order:"+orderUID).Result()
	if err != nil {
		return nil, err
	}

	var order models.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (c *OrderCache) Close() error {
	return c.client.Close()
}
