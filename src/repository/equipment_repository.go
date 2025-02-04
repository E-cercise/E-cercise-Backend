package repository

import (
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type EquipmentRepository interface {
	FindEquipmentList(q string, muscleGroup []string, paginator *helper.Paginator) ([]model.Equipment, error)
	CreateEquipment(tx *gorm.DB, eq model.Equipment) error
	AddAttributes(tx *gorm.DB, attr []model.Attribute) error
	FindByID(eqID uuid.UUID) (*model.Equipment, error)
	FindByIDTransaction(tx *gorm.DB, eqID uuid.UUID) (*model.Equipment, error)
	AddEquipmentOption(tx *gorm.DB, options model.EquipmentOption) error
	SaveEquipmentOption(tx *gorm.DB, option model.EquipmentOption) error
	SaveEquipmentFeature(tx *gorm.DB, option model.EquipmentFeature) error
	DeleteEquipmentOption(tx *gorm.DB, optID []uuid.UUID) error
	DeleteEquipmentFeature(tx *gorm.DB, optID []uuid.UUID) error
	SaveAttributes(tx *gorm.DB, attr *model.Attribute) error
	DeletesAttributes(tx *gorm.DB, attrID []uuid.UUID) error
	SaveEquipment(tx *gorm.DB, equipment *model.Equipment) error
	CreateEquipmentFeatures(tx *gorm.DB, features []model.EquipmentFeature) error
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

	err := query.Preload("EquipmentOptions.Images").Offset(paginator.Offset()).
		Limit(paginator.Limit).Find(&equipments).Error

	return equipments, err
}

func (r *equipmentRepository) CreateEquipment(tx *gorm.DB, eq model.Equipment) error {
	return tx.Create(&eq).Error
}

func (r *equipmentRepository) AddEquipmentOption(tx *gorm.DB, options model.EquipmentOption) error {
	return tx.Create(&options).Error
}

func (r *equipmentRepository) DeleteEquipmentOption(tx *gorm.DB, optID []uuid.UUID) error {
	return tx.Where("id IN ?", optID).Delete(&model.EquipmentOption{}).Error
}

func (r *equipmentRepository) DeleteEquipmentFeature(tx *gorm.DB, optID []uuid.UUID) error {
	return tx.Where("id IN ?", optID).Delete(&model.EquipmentFeature{}).Error
}

func (r *equipmentRepository) AddAttributes(tx *gorm.DB, attr []model.Attribute) error {
	return tx.Create(&attr).Error
}

func (r *equipmentRepository) FindByID(eqID uuid.UUID) (*model.Equipment, error) {
	var equipment *model.Equipment

	err := r.db.Preload("Images").
		Preload("MuscleGroups").
		Preload("EquipmentOptions").
		Preload("EquipmentFeatures").
		Preload("Attribute").
		First(&equipment, "id = ?", eqID).Error

	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (r *equipmentRepository) FindByIDTransaction(tx *gorm.DB, eqID uuid.UUID) (*model.Equipment, error) {
	var equipment *model.Equipment

	err := tx.Preload("Images").
		Preload("MuscleGroups").
		Preload("EquipmentOptions").
		Preload("Attribute").
		First(&equipment, "id = ?", eqID).Error

	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (r *equipmentRepository) SaveEquipmentOption(tx *gorm.DB, option model.EquipmentOption) error {
	return tx.Save(&option).Error
}

func (r *equipmentRepository) SaveEquipmentFeature(tx *gorm.DB, feature model.EquipmentFeature) error {
	return tx.Save(&feature).Error
}

func (r *equipmentRepository) SaveAttributes(tx *gorm.DB, attr *model.Attribute) error {
	return tx.Save(&attr).Error
}

func (r *equipmentRepository) DeletesAttributes(tx *gorm.DB, attrID []uuid.UUID) error {
	return tx.Delete(model.Attribute{}, attrID).Error
}

func (r *equipmentRepository) SaveEquipment(tx *gorm.DB, equipment *model.Equipment) error {
	return tx.Save(&equipment).Error
}

func (r *equipmentRepository) CreateEquipmentFeatures(tx *gorm.DB, features []model.EquipmentFeature) error {
	return tx.Create(features).Error
}
