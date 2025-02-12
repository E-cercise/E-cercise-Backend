package service

import (
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/data/response"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartService interface {
	AddEquipmentToCart(req request.CartItemPostRequest, userID uuid.UUID) error
	DeleteLineEquipmentInCart(lineEquipmentID uuid.UUID) (string, error)
	GetAllLineEquipmentInCart(userID uuid.UUID) (*response.GetCartItemResponse, error)
	ModifyLineEquipmentInCart(req request.CartItemPutRequest) error
}

type cartService struct {
	db            *gorm.DB
	cartRepo      repository.CartRepository
	equipmentRepo repository.EquipmentRepository
}

func NewCartService(db *gorm.DB, cartRepo repository.CartRepository, equipmentRepo repository.EquipmentRepository) CartService {
	return &cartService{db: db, cartRepo: cartRepo, equipmentRepo: equipmentRepo}
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

	_, err = s.equipmentRepo.FindByID(equipmentID)
	if err != nil {
		logger.Log.WithError(err).Error("cant find equipment ID:", equipmentID)
		return fmt.Errorf("equipment ID: %v not found", equipmentID)
	}

	_, err = s.equipmentRepo.FindOptionByID(eqpOptID)
	if err != nil {
		logger.Log.WithError(err).Error("cant find equipment Option ID:", eqpOptID)
		return fmt.Errorf("equipmentOptionID: %v not found", eqpOptID)
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

func (s *cartService) GetAllLineEquipmentInCart(userID uuid.UUID) (*response.GetCartItemResponse, error) {

	cart, err := s.cartRepo.GetCart(userID)

	if err != nil {
		logger.Log.WithError(err).Error("error finding cart in db with userID: ", userID)
		return nil, err
	}

	var resp response.GetCartItemResponse
	total := 0.0
	for _, line := range cart.LineEquipments {

		equipment, err := s.equipmentRepo.FindByID(line.EquipmentID)
		if err != nil {
			logger.Log.WithError(err).Error("error during find equipment ID", equipment.ID)
			return nil, err
		}

		equipmentOption, err := s.equipmentRepo.FindOptionByID(line.EquipmentOptionID)
		if err != nil {
			logger.Log.WithError(err).Error("error during find equipmentOpyion ID", equipmentOption.ID)
			return nil, err
		}

		lineTotal := float64(line.Quantity) * equipmentOption.Price
		total += lineTotal

		resp.LineEquipments = append(resp.LineEquipments, response.LineEquipment{
			EquipmentName:   fmt.Sprintf("%v: %v", equipment.Name, equipmentOption.Name),
			LineEquipmentID: line.ID.String(),
			Quantity:        line.Quantity,
			Total:           lineTotal,
		})

	}

	resp.TotalPrice = total

	return &resp, nil
}

func (s *cartService) ModifyLineEquipmentInCart(req request.CartItemPutRequest) error {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range req.Items {
		err := s.cartRepo.ModifyLineItem(tx, uuid.MustParse(item.LineEquipmentID), item.Quantity)
		if err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("cant modify line equipment with ID:", item.LineEquipmentID)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
