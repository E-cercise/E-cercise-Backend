package repository

import (
	"fmt"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MuscleGroupRepository interface {
	FindByID(tx *gorm.DB, mID string) (*model.MuscleGroup, error)
	AddGroup(tx *gorm.DB, groupsID []string, eqID uuid.UUID) error
}

type muscleGroupRepository struct {
	db *gorm.DB
}

func NewMuscleGroupRepository(db *gorm.DB) MuscleGroupRepository {
	return &muscleGroupRepository{db: db}
}

func (r *muscleGroupRepository) FindByID(tx *gorm.DB, mID string) (*model.MuscleGroup, error) {
	var muscleGroup *model.MuscleGroup

	if err := tx.Where("id = ?", mID).First(&muscleGroup).Error; err != nil {
		logger.Log.WithError(err).Error("cant find muscleGroup ID", mID)
		return nil, err
	}

	return muscleGroup, nil
}

func (r *muscleGroupRepository) AddGroup(tx *gorm.DB, groupsID []string, eqID uuid.UUID) error {
	if len(groupsID) == 0 {
		return fmt.Errorf("groupIDs cannot be empty")
	}

	var associations []map[string]interface{}
	for _, groupID := range groupsID {
		associations = append(associations, map[string]interface{}{
			"equipment_id":    eqID,
			"muscle_group_id": groupID,
		})
	}

	if err := tx.Table("equipment_muscle_groups").Create(associations).Error; err != nil {
		return fmt.Errorf("failed to associate muscle groups with equipment: %w", err)
	}
	return nil
}
