package service

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/data/response"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"gorm.io/gorm"
	"strings"
)

type EquipmentService interface {
	GetEquipmentData(q request.EquipmentListRequest, paginatior *helper.Paginator) (*response.EquipmentsResponse, error)
	AddEquipment(req request.EquipmentPostRequest) error
}

type equipmentService struct {
	db              *gorm.DB
	equipmentRepo   repository.EquipmentRepository
	muscleGroupRepo repository.MuscleGroupRepository
}

func NewEquipmentService(db *gorm.DB, equipmentRepo repository.EquipmentRepository, muscleGroupRepo repository.MuscleGroupRepository) EquipmentService {
	return &equipmentService{db: db, equipmentRepo: equipmentRepo, muscleGroupRepo: muscleGroupRepo}
}

func (s *equipmentService) GetEquipmentData(q request.EquipmentListRequest, paginator *helper.Paginator) (*response.EquipmentsResponse, error) {
	muscleGroup := strings.Split(q.MuscleGroup, ",")
	equipments, err := s.equipmentRepo.FindEquipmentList(q.Q, muscleGroup, paginator)

	if err != nil {
		logger.Log.WithError(err).Error("error during find all equipments")
		return nil, err
	}

	var resp response.EquipmentsResponse

	for _, equipment := range equipments {

		primaryImage := helper.FindPrimaryImageFromEquipment(equipment)
		var imagePath string
		if primaryImage == nil {
			imagePath = "https://res.cloudinary.com/drwodnunx/image/upload/v1736740947/temp/img_20250113110225.jpg.jpg"
		} else {
			imagePath = primaryImage.CloudinaryPath
		}

		price := findEquipmentMinimumPrice(equipment)

		eq := response.Equipment{
			ID:        equipment.ID,
			Name:      equipment.Name,
			Price:     price,
			ImagePath: imagePath,
		}
		resp.Equipments = append(resp.Equipments, eq)

	}
	return &resp, nil

}

func findEquipmentMinimumPrice(equipment model.Equipment) float64 {
	minimumPrice := equipment.EquipmentOptions[0].Price
	for _, option := range equipment.EquipmentOptions {
		if option.Price < minimumPrice {
			minimumPrice = option.Price
		}
	}
	return minimumPrice
}

func (s *equipmentService) AddEquipment(req request.EquipmentPostRequest) error {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	newEquipment := model.Equipment{
		Name:           req.Name,
		Brand:          req.Band,
		Model:          req.Model,
		Color:          req.Color,
		Material:       req.Material,
		SpecialFeature: req.SpecialFeature,
	}

	err := s.equipmentRepo.CreateEquipment(tx, newEquipment)
	if err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error creating equipment")
		return err
	}

	for _, option := range req.Option {
		newOption := model.EquipmentOption{
			EquipmentID:       newEquipment.ID,
			Weight:            option.Weight,
			Price:             option.Price,
			RemainingProducts: option.Available,
		}

		if err := s.equipmentRepo.AddEquipmentOption(tx, newOption); err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("error adding equipment options", newOption)
			return err
		}
	}

	for _, groupID := range req.MuscleGroupUsed {
		if _, err := s.muscleGroupRepo.FindByID(tx, groupID); err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("error finding muscleGroupID", groupID)
			return err
		}
	}

	if err := s.muscleGroupRepo.AddGroup(tx, req.MuscleGroupUsed, newEquipment.ID); err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error adding muscle Group to equipment")
		return err
	}

	//TODO: add images db (maybe archive image too)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
}