package service

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	AddEquipmentToCart(req request.CheckoutCartRequest, userID uuid.UUID) error
}

type orderService struct {
	db            *gorm.DB
	cartRepo      repository.CartRepository
	equipmentRepo repository.EquipmentRepository
}

func NewOrderService(db *gorm.DB, cartRepo repository.CartRepository, equipmentRepo repository.EquipmentRepository) OrderService {
	return &orderService{db: db, cartRepo: cartRepo, equipmentRepo: equipmentRepo}
}
