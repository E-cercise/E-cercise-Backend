package repository

import (
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EquipmentRepository interface {
	FindEquipmentList(q string, muscleGroup []string, paginator *helper.Paginator, category string, minBudget int, maxBudget int) ([]model.Equipment, error)
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
	FindOptionByID(optionID uuid.UUID) (*model.EquipmentOption, error)
	DeleteEquipment(tx *gorm.DB, eqID uuid.UUID) error
	GetAllEquipmentCategories() ([]model.Equipment, error)
	FindByIDs(ids []uuid.UUID) ([]model.Equipment, error)
	FindEquipmentsByOptionIDs(optionIDs []uuid.UUID) ([]model.Equipment, error)
}

type equipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) EquipmentRepository {
	return &equipmentRepository{db: db}
}

func (r *equipmentRepository) FindEquipmentList(q string, muscleGroup []string, paginator *helper.Paginator, category string, minBudget int, maxBudget int) ([]model.Equipment, error) {
	var equipments []model.Equipment

	query := r.db.Model(&model.Equipment{})

	if q != "" {
		query = query.Where("equipment.name ILIKE ? OR equipment.description ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	if category != "" {
		query = query.Where("equipment.category = ?", category)
	}

	if minBudget != 0 || maxBudget != 0 {
		if minBudget != 0 && maxBudget != 0 {
			query = query.Joins("JOIN equipment_options eo ON eo.equipment_id = equipment.id").
				Where("eo.price BETWEEN ? and ?", minBudget, maxBudget)
		} else if minBudget != 0 {
			query = query.Joins("JOIN equipment_options eo ON eo.equipment_id = equipment.id").
				Where("eo.price >= ?", minBudget)
		} else if maxBudget != 0 {
			query = query.Joins("JOIN equipment_options eo ON eo.equipment_id = equipment.id").
				Where("eo.price <= ?", maxBudget)
		}
		query = query.Select("DISTINCT equipment.*")
	}

	if len(muscleGroup) > 0 {
		query = query.
			Joins("JOIN equipment_muscle_groups emg ON emg.equipment_id = equipment.id").
			Where("emg.muscle_group_id ILIKE ANY (?)", pq.Array(muscleGroup)).
			Group("equipment.id").
			Having("COUNT(DISTINCT emg.muscle_group_id) = ?", len(muscleGroup))
	}

	if err := query.Count(&paginator.TotalRows).Error; err != nil {
		return nil, err
	}
	paginator.CalculateTotalPages()

	err := query.Preload("EquipmentOptions.Images").Preload("MuscleGroups").Offset(paginator.Offset()).
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

func (r *equipmentRepository) DeleteEquipmentFeature(tx *gorm.DB, featIDs []uuid.UUID) error {
	return tx.Where("id IN ?", featIDs).Delete(&model.EquipmentFeature{}).Error
}

func (r *equipmentRepository) AddAttributes(tx *gorm.DB, attr []model.Attribute) error {
	return tx.Create(&attr).Error
}

func (r *equipmentRepository) FindByID(eqID uuid.UUID) (*model.Equipment, error) {
	var equipment *model.Equipment

	err := r.db.Preload("MuscleGroups").
		Preload("EquipmentOptions").
		Preload("EquipmentOptions.Images").
		Preload("EquipmentFeature").
		Preload("Attribute").
		First(&equipment, "id = ?", eqID).Error

	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (r *equipmentRepository) FindOptionByID(optionID uuid.UUID) (*model.EquipmentOption, error) {
	var opt *model.EquipmentOption
	err := r.db.Preload("Images").Find(&opt, "id = ?", optionID).Error
	return opt, err
}

func (r *equipmentRepository) FindByIDTransaction(tx *gorm.DB, eqID uuid.UUID) (*model.Equipment, error) {
	var equipment *model.Equipment

	err := tx.Preload("MuscleGroups").
		Preload("EquipmentOptions").
		Preload("EquipmentOptions.Images").
		Preload("EquipmentFeature").
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
	logger.Log.Infof("🧪 DEBUG Delete Attribute IDs: %v", attrID)

	debugTx := tx.Session(&gorm.Session{Logger: tx.Logger.LogMode(4)}) // 4 = Info level

	result := debugTx.Where("id IN ?", attrID).Delete(&model.Attribute{})
	logger.Log.Infof("🧾 Rows affected in attribute delete: %d", result.RowsAffected)

	return result.Error
}

func (r *equipmentRepository) SaveEquipment(tx *gorm.DB, equipment *model.Equipment) error {
	return tx.Save(&equipment).Error
}

func (r *equipmentRepository) CreateEquipmentFeatures(tx *gorm.DB, features []model.EquipmentFeature) error {
	return tx.Create(features).Error
}

func (r *equipmentRepository) DeleteEquipment(tx *gorm.DB, eqID uuid.UUID) error {
	return tx.Select(clause.Associations).Delete(&model.Equipment{ID: eqID}).Error
}

func (r *equipmentRepository) GetAllEquipmentCategories() ([]model.Equipment, error) {
	var equipments []model.Equipment
	err := r.db.Distinct("category").Find(&equipments).Error
	if err != nil {
		return nil, err
	}
	return equipments, nil
}

func (r *equipmentRepository) FindByIDs(ids []uuid.UUID) ([]model.Equipment, error) {
	var equipments []model.Equipment
	err := r.db.Preload("MuscleGroups").
		Preload("EquipmentOptions").
		Preload("EquipmentOptions.Images").
		Preload("EquipmentFeature").
		Preload("Attribute").
		Where("id IN ?", ids).
		Find(&equipments).Error
	if err != nil {
		return nil, err
	}

	return equipments, nil
}

func (r *equipmentRepository) FindEquipmentsByOptionIDs(optionIDs []uuid.UUID) ([]model.Equipment, error) {
	var equipments []model.Equipment

	var equipmentIDs []uuid.UUID
	if err := r.db.Model(&model.EquipmentOption{}).
		Select("DISTINCT(equipment_id)").
		Where("id IN ?", optionIDs).
		Pluck("equipment_id", &equipmentIDs).Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Preload("EquipmentOptions", "id IN ?", optionIDs).
		Preload("EquipmentOptions.Images").
		Preload("MuscleGroups").
		Find(&equipments, "id IN ?", equipmentIDs).Error; err != nil {
		return nil, err
	}

	return equipments, nil
}
