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
	UpdateGroups(tx *gorm.DB, groupsID []string, eqID uuid.UUID) error
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

func (r *muscleGroupRepository) UpdateGroups(tx *gorm.DB, groupsID []string, eqID uuid.UUID) error {
	// Fetch the current muscle group associations for the equipment
	var currentGroupIDs []string
	if err := tx.Table("equipment_muscle_groups").
		Where("equipment_id = ?", eqID).
		Pluck("muscle_group_id", &currentGroupIDs).Error; err != nil {
		logger.Log.WithError(err).Error("failed to fetch current muscle group associations")
		return fmt.Errorf("failed to fetch current muscle group associations: %w", err)
	}

	// Determine the muscle groups to add and delete
	toAdd := difference(groupsID, currentGroupIDs)
	toDelete := difference(currentGroupIDs, groupsID)

	// Add new associations
	var associations []map[string]interface{}
	for _, groupID := range toAdd {
		associations = append(associations, map[string]interface{}{
			"equipment_id":    eqID,
			"muscle_group_id": groupID,
		})
	}

	if len(associations) > 0 {
		if err := tx.Table("equipment_muscle_groups").Create(associations).Error; err != nil {
			logger.Log.WithError(err).Error("failed to add new muscle group associations")
			return fmt.Errorf("failed to add new muscle group associations: %w", err)
		}
	}

	// Delete obsolete associations
	if len(toDelete) > 0 {
		if err := tx.Table("equipment_muscle_groups").
			Where("equipment_id = ? AND muscle_group_id IN ?", eqID, toDelete).
			Delete(nil).Error; err != nil {
			logger.Log.WithError(err).Error("failed to delete obsolete muscle group associations")
			return fmt.Errorf("failed to delete obsolete muscle group associations: %w", err)
		}
	}

	return nil
}

// Helper function to calculate the difference between two slices
func difference(slice1, slice2 []string) []string {
	set := make(map[string]struct{}, len(slice2))
	for _, s := range slice2 {
		set[s] = struct{}{}
	}

	var diff []string
	for _, s := range slice1 {
		if _, found := set[s]; !found {
			diff = append(diff, s)
		}
	}
	return diff
}
