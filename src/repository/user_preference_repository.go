package repository

import (
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPreferenceRepository interface {
	SetPreferences(userID uuid.UUID, tagIDs []uuid.UUID) error
	GetPreferences(userID uuid.UUID) ([]model.Tag, error)
}

type userPrefRepo struct {
	db *gorm.DB
}

func NewUserPreferenceRepository(db *gorm.DB) UserPreferenceRepository {
	return &userPrefRepo{db}
}

func (r *userPrefRepo) SetPreferences(userID uuid.UUID, tagIDs []uuid.UUID) error {

	r.db.Where("user_id = ?", userID).Delete(&model.UserPreference{})

	for _, tagID := range tagIDs {
		r.db.Create(&model.UserPreference{
			UserID: userID,
			TagID:  tagID,
		})
	}
	return nil
}

func (r *userPrefRepo) GetPreferences(userID uuid.UUID) ([]model.Tag, error) {
	var prefs []model.UserPreference
	err := r.db.Preload("Tag").Where("user_id = ?", userID).Find(&prefs).Error
	if err != nil {
		return nil, err
	}

	var tags []model.Tag
	for _, p := range prefs {
		tags = append(tags, p.Tag)
	}
	return tags, nil
}
