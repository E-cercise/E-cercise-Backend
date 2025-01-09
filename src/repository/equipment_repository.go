package repository

import (
	"github.com/E-cercise/E-cercise/src/model"
	"gorm.io/gorm"
)

type EquipmentRepository interface {
	FindAll() ([]model.Equipment, error)
}

type equipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *equipmentRepository) FindAll() ([]model.Equipment, error) {
	var equipments []model.Equipment
	err := r.db.Find(&equipments).Error

	return equipments, err
}
