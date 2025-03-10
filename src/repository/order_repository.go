package repository

import (
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(tx *gorm.DB, order *model.Order) error
	SaveOrder(tx *gorm.DB, order *model.Order) error
	GetOrderDetail(orderID uuid.UUID) (*model.Order, error)
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

func (r *orderRepository) GetOrderDetail(orderID uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("LineEquipments").Find(&order, "id = ?", orderID).Error
	return &order, err
}
