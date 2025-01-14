package service

import (
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/data/response"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/repository"
	"gorm.io/gorm"
	"strings"
)

type EquipmentService interface {
	GetEquipmentData(q request.EquipmentListRequest, paginator *helper.Paginator) (*response.EquipmentsResponse, error)
}

type equipmentService struct {
	db            *gorm.DB
	equipmentRepo repository.EquipmentRepository
}

func NewEquipmentService(db *gorm.DB, equipmentRepo repository.EquipmentRepository) EquipmentService {
	return &equipmentService{db: db, equipmentRepo: equipmentRepo}
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

		eq := response.Equipment{
			ID:        equipment.ID,
			Name:      equipment.Name,
			Price:     equipment.Price,
			ImagePath: imagePath,
		}
		resp.Equipments = append(resp.Equipments, eq)

	}
	return &resp, nil

}
