package repository

import (
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/model"
	"gorm.io/gorm"
	"strings"
)

type EquipmentRepository interface {
	FindEquipmentList(q string, muscleGroup []string, paginator *helper.Paginator) ([]model.Equipment, error)
	CreateEquipment(tx *gorm.DB, eq model.Equipment) error
	AddEquipmentOption(tx *gorm.DB, options model.EquipmentOption) error
	AddAAttributes(tx *gorm.DB, attr []model.Attribute) error
}

type equipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) EquipmentRepository {
	return &equipmentRepository{db: db}
}

func (r *equipmentRepository) FindEquipmentList(q string, muscleGroup []string, paginator *helper.Paginator) ([]model.Equipment, error) {
	var equipments []model.Equipment

	query := r.db.Model(&model.Equipment{})

	if q != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	if len(muscleGroup) > 0 {
		nameConditions := make([]string, len(muscleGroup))
		args := make([]interface{}, len(muscleGroup))

		for i, group := range muscleGroup {
			nameConditions[i] = "name ILIKE ?"
			args[i] = "%" + group + "%"
		}

		query = query.Where(strings.Join(nameConditions, " OR "), args...)
	}

	if err := query.Count(&paginator.TotalRows).Error; err != nil {
		return nil, err
	}
	paginator.CalculateTotalPages()

	err := query.Preload("EquipmentOptions").Offset(paginator.Offset()).
		Limit(paginator.Limit).Find(&equipments).Error

	return equipments, err
}

func (r *equipmentRepository) CreateEquipment(tx *gorm.DB, eq model.Equipment) error {
	return tx.Create(&eq).Error
}

func (r *equipmentRepository) AddEquipmentOption(tx *gorm.DB, options model.EquipmentOption) error {
	return tx.Create(&options).Error
}

func (r *equipmentRepository) AddAAttributes(tx *gorm.DB, attr []model.Attribute) error {
	return tx.Create(&attr).Error
}
