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
	CreateOrder(req request.CheckoutCartRequest, userID uuid.UUID) error
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

func (s *orderService) CreateOrder(req request.CheckoutCartRequest, userID uuid.UUID) error {
	logger.Log.Info("Starting order creation process", map[string]interface{}{"user_id": userID})

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			logger.Log.WithError(fmt.Errorf("%v", r)).Error("Panic recovered during order creation")
			tx.Rollback()
		}
	}()

	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get cart for user", map[string]interface{}{"user_id": userID})
		tx.Rollback()
		return fmt.Errorf("failed to get cart for user %s: %v", userID, err)
	}

	cartItemsMap := make(map[uuid.UUID]model.LineEquipment)
	for _, item := range cart.LineEquipments {
		cartItemsMap[item.ID] = item
	}

	// Initialize order
	order := &model.Order{
		ID:             uuid.New(),
		UserID:         userID,
		LineEquipments: []model.LineEquipment{},
		OrderStatus:    enum.OrderPlaced,
	}

	var totalPrice float64

	for _, checkoutItem := range req.Items {
		lineEquipmentID, err := uuid.Parse(checkoutItem.LineEquipmentID)
		if err != nil {
			logger.Log.WithError(err).Error("Invalid line equipment ID", map[string]interface{}{"line_equipment_id": checkoutItem.LineEquipmentID})
			tx.Rollback()
			return fmt.Errorf("invalid line equipment ID: %s", checkoutItem.LineEquipmentID)
		}

		cartItem, exists := cartItemsMap[lineEquipmentID]
		if !exists {
			errMsg := fmt.Sprintf("Line item not found in cart: %s", checkoutItem.LineEquipmentID)
			logger.Log.Error(errMsg, map[string]interface{}{"user_id": userID, "line_equipment_id": checkoutItem.LineEquipmentID})
			tx.Rollback()
			return errors.New(errMsg)
		}

		if checkoutItem.Quantity > cartItem.Quantity {
			errMsg := fmt.Sprintf("Insufficient quantity for line item %s (requested: %d, available: %d)", lineEquipmentID, checkoutItem.Quantity, cartItem.Quantity)
			logger.Log.Error(errMsg, map[string]interface{}{"user_id": userID, "line_equipment_id": checkoutItem.LineEquipmentID})
			tx.Rollback()
			return errors.New(errMsg)
		}

		equipmentOption, err := s.equipmentRepo.FindOptionByID(cartItem.EquipmentOptionID)
		if err != nil {
			logger.Log.WithError(err).Error("Failed to fetch equipment option", map[string]interface{}{"equipment_option_id": cartItem.EquipmentOptionID})
			tx.Rollback()
			return fmt.Errorf("failed to fetch equipment option %s: %v", cartItem.EquipmentOptionID, err)
		}

		if equipmentOption.RemainingProducts < checkoutItem.Quantity {
			errMsg := fmt.Sprintf("Insufficient stock for equipment option %s", cartItem.EquipmentOptionID)
			logger.Log.Error(errMsg, map[string]interface{}{"equipment_option_id": cartItem.EquipmentOptionID})
			tx.Rollback()
			return errors.New(errMsg)
		}

		equipmentOption.RemainingProducts -= checkoutItem.Quantity
		if err := tx.Save(equipmentOption).Error; err != nil {
			logger.Log.WithError(err).Error("Failed to update inventory", map[string]interface{}{"equipment_option_id": cartItem.EquipmentOptionID})
			tx.Rollback()
			return fmt.Errorf("failed to update inventory: %v", err)
		}

		totalPrice += float64(checkoutItem.Quantity) * equipmentOption.Price
		remainingQuantity := cartItem.Quantity - checkoutItem.Quantity

		if remainingQuantity == 0 {
			// Instead of deleting, we simply change cartID to nil and attach orderID
			cartItem.CartID = nil
			cartItem.OrderID = &order.ID
			if err := tx.Save(cartItem).Error; err != nil {
				logger.Log.WithError(err).Error("Failed to update cart item to order item", map[string]interface{}{"cart_item_id": cartItem.ID})
				tx.Rollback()
				return fmt.Errorf("failed to update cart item to order item: %v", err)
			}
			logger.Log.Info("Updated cart item to be part of the order", map[string]interface{}{"cart_item_id": cartItem.ID, "order_id": order.ID})
		} else {
			// Update the cart item quantity
			if err := s.cartRepo.ModifyLineItem(tx, cartItem.ID, remainingQuantity); err != nil {
				logger.Log.WithError(err).Error("Failed to update cart item quantity", map[string]interface{}{"cart_item_id": cartItem.ID, "remaining_quantity": remainingQuantity})
				tx.Rollback()
				return fmt.Errorf("failed to update cart item quantity: %v", err)
			}

			newOrderLineItem := model.LineEquipment{
				ID:                uuid.New(),
				OrderID:           &order.ID,
				CartID:            nil,
				EquipmentID:       cartItem.EquipmentID,
				EquipmentOptionID: cartItem.EquipmentOptionID,
				Quantity:          checkoutItem.Quantity,
			}

			// Add new line item to order
			order.LineEquipments = append(order.LineEquipments, newOrderLineItem)
			logger.Log.Info("Created new line item for order", map[string]interface{}{"order_id": order.ID, "line_equipment_id": newOrderLineItem.ID})
		}
	}

	order.TotalPrice = totalPrice

	if err := s.orderRepo.CreateOrder(tx, order); err != nil {
		logger.Log.WithError(err).Error("Failed to create order", map[string]interface{}{"order_id": order.ID})
		tx.Rollback()
		return fmt.Errorf("failed to create order: %v", err)
	}

	logger.Log.Info("Order created successfully", map[string]interface{}{"order_id": order.ID, "total_price": totalPrice})

	if err := tx.Commit().Error; err != nil {
		logger.Log.WithError(err).Error("error committing transaction")
		tx.Rollback()
		return err
	}
	return nil

}
