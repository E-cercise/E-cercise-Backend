package repository

import (
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/model"
	"gorm.io/gorm"
	"strings"
)

type EquipmentRepository interface {
	FindEquipmentList(q string, muscleGroup []string, paginator *helper.Paginator) ([]model.Equipment, error)
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
		nameConditions := make([]string, 0)
		args := make([]interface{}, 0)
		for _, group := range muscleGroup {
			cond := "ILIKE ?"
			nameConditions = append(nameConditions, cond)
			likePattern := "%" + group + "%"
			args = append(args, likePattern, likePattern)
		}
		query = query.Where(strings.Join(nameConditions, " OR "), args...)
	}

	if err := query.Count(&paginator.TotalRows).Error; err != nil {
		return nil, err
	}
	paginator.CalculateTotalPages()

	err := r.db.Offset(paginator.Offset()).Limit(paginator.Limit).Find(&equipments).Error

	return equipments, err
}
