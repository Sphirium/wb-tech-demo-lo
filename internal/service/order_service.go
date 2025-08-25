package service

import (
	"encoding/json"
	"time"

	"github.com/Sphirium/wb-tech-demo-lo/internal/cache"
	"github.com/Sphirium/wb-tech-demo-lo/internal/models"
	"github.com/Sphirium/wb-tech-demo-lo/internal/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	repo  repository.OrderRepositoryInterface
	cache cache.OrderCacheInterface
}

func NewOrderService(repo repository.OrderRepositoryInterface, cache cache.OrderCacheInterface) *OrderService {
	return &OrderService{repo: repo, cache: cache}
}

func (s *OrderService) SaveOrder(data []byte) error {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return err
	}

	if order.OrderUID == "" {
		order.OrderUID = uuid.New().String()
	}

	if order.Delivery != nil {
		order.Delivery.OrderID = order.OrderUID
	}
	if order.Payment != nil {
		order.Payment.OrderID = order.OrderUID
		order.Payment.PaymentDtTime = time.Unix(order.Payment.PaymentDt, 0)
	}
	for i := range order.Items {
		order.Items[i].OrderID = order.OrderUID
	}

	return s.repo.Create(&order)
}

func (s *OrderService) GetOrderByUID(orderUID string) (*models.Order, error) {
	if order, err := s.cache.Get(orderUID); err == nil {
		return order, nil
	}

	order, err := s.repo.FindByOrderUID(orderUID)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(order)
	return order, nil
}

func (s *OrderService) RestoreCacheFromDB() error {
	uids, err := s.repo.GetAllOrderUIDs()
	if err != nil {
		return err
	}

	for _, uid := range uids {
		order, err := s.repo.FindByOrderUID(uid)
		if err != nil {
			continue
		}
		_ = s.cache.Set(order)
	}
	return nil
}
