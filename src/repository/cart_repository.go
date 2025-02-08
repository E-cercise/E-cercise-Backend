package repository

import (
	"errors"
	"fmt"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	//FindByID(tx *gorm.DB, mID string) (*model.MuscleGroup, error)
	AddLineItem(userID uuid.UUID, lineEquipment *model.LineEquipment) error
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
