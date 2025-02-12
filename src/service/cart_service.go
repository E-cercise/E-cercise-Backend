package service

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartService interface {
	AddEquipmentToCart(req request.CartItemPostRequest, userID uuid.UUID) error
	DeleteLineEquipmentInCart(lineEquipmentID uuid.UUID) (string, error)
}

type cartService struct {
	db       *gorm.DB
	cartRepo repository.CartRepository
}

func NewCartService(db *gorm.DB, cartRepo repository.CartRepository) CartService {
	return &cartService{db: db, cartRepo: cartRepo}
}

func (s *cartService) AddEquipmentToCart(req request.CartItemPostRequest, userID uuid.UUID) error {

	equipmentID, err := uuid.Parse(req.EquipmentID)
	if err != nil {
		logger.Log.WithError(err).Error("error parsing request body")
		return err
	}

	eqpOptID, err := uuid.Parse(req.EquipmentOptionID)
	if err != nil {
		logger.Log.WithError(err).Error("error parsing equipmentID ")
		return err
	}

	newLineEquipment := model.LineEquipment{
		EquipmentID:       equipmentID,
		EquipmentOptionID: eqpOptID,
		Quantity:          req.Quantity,
	}

	err = s.cartRepo.AddLineItem(userID, &newLineEquipment)
	if err != nil {
		logger.Log.WithError(err).Error("error adding item into cart", "item", newLineEquipment)
		return err
	}

	return nil
}

func (s *cartService) DeleteLineEquipmentInCart(lineEquipmentID uuid.UUID) (string, error) {

	recordCount, err := s.cartRepo.DeleteLineItem(lineEquipmentID)
	if err != nil {
		logger.Log.WithError(err).Error("error deleting line item ID: ", lineEquipmentID)
		return "error", err
	}

	if recordCount == 0 {
		logger.Log.Warning("user trying to delete line item that doesnt exists")
		return "204", nil
	}

	return "success", nil
}
