package repository

import (
	"github.com/E-cercise/E-cercise/src/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(tx *gorm.DB, order *model.Order) error
	SaveOrder(tx *gorm.DB, order *model.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(tx *gorm.DB, order *model.Order) error {
	return tx.Create(order).Error
}

func (r *orderRepository) SaveOrder(tx *gorm.DB, order *model.Order) error {
	return tx.Save(order).Error
}
