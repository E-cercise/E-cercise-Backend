package service

import (
	"github.com/E-cercise/E-cercise/src/data/response"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/repository"
	"gorm.io/gorm"
)

type EquipmentService interface {
	GetAllEquipmentData() (*response.EquipmentsResponse, error)
}

type equipmentService struct {
	db            *gorm.DB
	equipmentRepo repository.EquipmentRepository
}

func NewEquipmentService(db *gorm.DB, equipmentRepo repository.EquipmentRepository) EquipmentService {
	return &equipmentService{db: db, equipmentRepo: equipmentRepo}
}

func (s *equipmentService) GetAllEquipmentData() (*response.EquipmentsResponse, error) {
	equipments, err := s.equipmentRepo.FindAll()

	if err != nil {
		logger.Log.WithError(err).Error("error during find all equipments")
		return nil, err
	}

	var resp response.EquipmentsResponse
	for _, equipment := range equipments {
		eq := response.Equipment{
			ID:    equipment.ID,
			Name:  equipment.Name,
			Price: equipment.Price,
		}
		resp.Equipments = append(resp.Equipments, eq)

	}
	return &resp, nil

}
