package repository

import (
	"github.com/Sphirium/wb-tech-demo-lo/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepositoryInterface interface {
	Create(order *models.Order) error
	FindByOrderUID(orderUID string) (*models.Order, error)
	GetAllOrderUIDs() ([]string, error)
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepositoryInterface {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		if order.Delivery != nil {
			// ON CONFLICT для delivery
			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(order.Delivery).Error; err != nil {
				return err
			}
		}

		if order.Payment != nil {
			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(order.Payment).Error; err != nil {
				return err
			}
		}

		if len(order.Items) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&order.Items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *OrderRepository) FindByOrderUID(orderUID string) (*models.Order, error) {
	var order models.Order
	if err := r.db.
		Preload("Delivery").
		Preload("Payment").
		Preload("Items").
		Where("order_uid = ?", orderUID).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetAllOrderUIDs() ([]string, error) {
	var uids []string
	err := r.db.Model(&models.Order{}).Pluck("order_uid", &uids).Error
	return uids, err
}
