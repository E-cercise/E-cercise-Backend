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

func (r *muscleGroupRepository) UpdateGroups(db *gorm.DB, groupsID []string, eqID uuid.UUID) error {
	logger.Log.Infof("🔄 Starting UpdateGroups for equipment %s", eqID)

	var currentGroupIDs []string
	if err := db.Table("equipment_muscle_groups").
		Where("equipment_id = ?", eqID).
		Pluck("muscle_group_id", &currentGroupIDs).Error; err != nil {
		logger.Log.WithError(err).Error("Failed to fetch current muscle group associations")
		return fmt.Errorf("failed to fetch current muscle group associations: %w", err)
	}

	toAdd := difference(groupsID, currentGroupIDs)
	toDelete := difference(currentGroupIDs, groupsID)

	logger.Log.Infof("🔍 To Add: %v", toAdd)
	logger.Log.Infof("🧹 To Delete: %v", toDelete)

	// Add new associations
	if len(toAdd) > 0 {
		var associations []map[string]interface{}
		for _, groupID := range toAdd {
			associations = append(associations, map[string]interface{}{
				"equipment_id":    eqID,
				"muscle_group_id": groupID,
			})
		}
		if err := db.Table("equipment_muscle_groups").Create(&associations).Error; err != nil {
			logger.Log.WithError(err).Error("Failed to add new muscle group associations")
			return fmt.Errorf("failed to add new muscle group associations: %w", err)
		}
	}

	if len(toDelete) > 0 {
		logger.Log.Infof("🗑 Executing Association Delete: equipment_id = %s, muscle_group_id IN %v", eqID, toDelete)

		equipment := model.Equipment{ID: eqID}
		var groupsToDelete []model.MuscleGroup
		for _, id := range toDelete {
			groupsToDelete = append(groupsToDelete, model.MuscleGroup{ID: id})
		}

		if err := db.Model(&equipment).Association("MuscleGroups").Delete(&groupsToDelete); err != nil {
			logger.Log.WithError(err).Error("Failed to delete muscle group associations via GORM association")
			return fmt.Errorf("failed to delete muscle group associations: %w", err)
		}

		logger.Log.Infof("✅ Deleted %d muscle group associations", len(groupsToDelete))
	}
	logger.Log.Infof("✅ UpdateGroups completed for equipment ID %s", eqID)
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
