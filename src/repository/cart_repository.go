package repository

import (
	"errors"
	"fmt"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	AddLineItem(userID uuid.UUID, lineEquipment *model.LineEquipment) error
	DeleteLineItem(lineEquipmentID uuid.UUID) (int64, error)
	GetCart(userID uuid.UUID) (*model.Cart, error)
	ClearAllLineItems(userID uuid.UUID) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddLineItem(userID uuid.UUID, lineEquipment *model.LineEquipment) error {
	var cart model.Cart

	if err := r.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cart not found for user ID: %s", userID)
		}
		return err
	}

	lineEquipment.CartID = cart.ID

	if err := r.db.Create(lineEquipment).Error; err != nil {
		return fmt.Errorf("failed to add line equipment to cart: %v", err)
	}

	return nil
}

func (r *cartRepository) DeleteLineItem(lineEquipmentID uuid.UUID) (int64, error) {
	res := r.db.Delete(model.LineEquipment{}, lineEquipmentID)

	return res.RowsAffected, res.Error
}

func (r *cartRepository) GetCart(userID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart

	if err := r.db.Preload("LineEquipments").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("cart not found for user ID: %s", userID)
		}
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) ClearAllLineItems(userID uuid.UUID) error {
	var cart model.Cart

	if err := r.db.Preload("LineEquipments").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cart not found for user ID: %s", userID)
		}
		return err
	}

	if err := r.db.Delete(cart.LineEquipments).Error; err != nil {
		return fmt.Errorf("failed to clear all line equipments in cart: %v", err)
	}

	return nil
}