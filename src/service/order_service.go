package service

import (
	"errors"
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(req request.CheckoutCartRequest, user model.User) error
}

type orderService struct {
	db            *gorm.DB
	cartRepo      repository.CartRepository
	equipmentRepo repository.EquipmentRepository
	orderRepo     repository.OrderRepository
}

func NewOrderService(db *gorm.DB, cartRepo repository.CartRepository, equipmentRepo repository.EquipmentRepository, orderRepo repository.OrderRepository) OrderService {
	return &orderService{db: db, cartRepo: cartRepo, equipmentRepo: equipmentRepo, orderRepo: orderRepo}
}

func (s *orderService) CreateOrder(req request.CheckoutCartRequest, user model.User) error {
	logger.Log.Info("Starting order creation process", map[string]interface{}{"user_id": user.ID})

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			logger.Log.WithError(fmt.Errorf("%v", r)).Error("Panic recovered during order creation")
			tx.Rollback()
		}
	}()

	cart, err := s.cartRepo.GetCart(user.ID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get cart for user", map[string]interface{}{"user_id": user.ID})
		tx.Rollback()
		return fmt.Errorf("failed to get cart for user %s: %v", user.ID, err)
	}

	cartItemsMap := make(map[uuid.UUID]model.LineEquipment)
	for _, item := range cart.LineEquipments {
		cartItemsMap[item.ID] = item
	}

	// Initialize order
	order := &model.Order{
		UserID:          user.ID,
		DeliveryAddress: user.Address,
		LineEquipments:  []model.LineEquipment{},
		OrderStatus:     enum.OrderPending,
	}

	if err := s.orderRepo.CreateOrder(tx, order); err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("Failed to create order", map[string]interface{}{"order_id": order.ID})
		return fmt.Errorf("failed to create order: %v", err)
	}

	var totalPrice float64
	for _, lineEquipmentID := range req.LineEquipments {

		cartItem, exists := cartItemsMap[lineEquipmentID]
		if !exists {
			errMsg := fmt.Sprintf("Line item not found in cart: %s", lineEquipmentID)
			logger.Log.Error(errMsg, map[string]interface{}{"user_id": user.ID, "line_equipment_id": lineEquipmentID})
			tx.Rollback()
			return errors.New(errMsg)
		}

		equipmentOption, err := s.equipmentRepo.FindOptionByID(cartItem.EquipmentOptionID)
		if err != nil {
			logger.Log.WithError(err).Error("Failed to fetch equipment option", map[string]interface{}{"equipment_option_id": cartItem.EquipmentOptionID})
			tx.Rollback()
			return fmt.Errorf("failed to fetch equipment option %s: %v", cartItem.EquipmentOptionID, err)
		}

		equipmentOption.RemainingProducts -= cartItem.Quantity
		if err := tx.Save(equipmentOption).Error; err != nil {
			logger.Log.WithError(err).Error("Failed to update inventory", map[string]interface{}{"equipment_option_id": cartItem.EquipmentOptionID})
			tx.Rollback()
			return fmt.Errorf("failed to update inventory: %v", err)
		}

		totalPrice += float64(cartItem.Quantity) * equipmentOption.Price
		cartItem.CartID = nil
		cartItem.OrderID = &order.ID
		if err := tx.Save(cartItem).Error; err != nil {
			logger.Log.WithError(err).Error("Failed to update cart item to order item", map[string]interface{}{"cart_item_id": cartItem.ID})
			tx.Rollback()
			return fmt.Errorf("failed to update cart item to order item: %v", err)
		}
		logger.Log.Info("Updated cart item to be part of the order", map[string]interface{}{"cart_item_id": cartItem.ID, "order_id": order.ID})
	}

	order.TotalPrice = totalPrice
	if err := s.orderRepo.SaveOrder(tx, order); err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("Failed to save order", map[string]interface{}{"order_id": order.ID})
		return fmt.Errorf("failed to save order: %v", err)
	}

	logger.Log.Info("Order created successfully", map[string]interface{}{"order_id": order.ID, "total_price": totalPrice})

	if err := tx.Commit().Error; err != nil {
		logger.Log.WithError(err).Error("error committing transaction")
		tx.Rollback()
		return err
	}
	return nil

}
