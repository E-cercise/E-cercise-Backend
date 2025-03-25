package repository

import (
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(tx *gorm.DB, order *model.Order) error
	SaveOrder(tx *gorm.DB, order *model.Order) error
	FindByID(orderID uuid.UUID) (*model.Order, error)
	UpdateOrderStatusByID(orderID uuid.UUID, orderStatus enum.OrderStatus) error
	FindByStatus(userID uuid.UUID, orderStatus enum.OrderStatus) ([]model.Order, error)
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

func (r *orderRepository) FindByID(orderID uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("LineEquipments").Find(&order, "id = ?", orderID).Error
	return &order, err
}

func (r *orderRepository) UpdateOrderStatusByID(orderID uuid.UUID, orderStatus enum.OrderStatus) error {
	var order model.Order
	return r.db.Model(&order).Where("id = ?", orderID).Update("order_status", orderStatus).Error
}

func (r *orderRepository) FindByStatus(userID uuid.UUID, orderStatus enum.OrderStatus) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("LineEquipments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC").Limit(1)
	}).Preload("LineEquipments.EquipmentOption").
		Preload("LineEquipments.EquipmentOption.Equipment").
		Where("user_id = ? AND order_status = ?", userID, orderStatus).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
