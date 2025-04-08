package repository

import (
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(userID string) (*model.User, error)
	SaveUser(user *model.User) error
	UpdateUserTransaction(tx *gorm.DB, userID *uuid.UUID, updateFields map[string]interface{}) error
	UpdateUserPreferences(tx *gorm.DB, user *model.User, pref []model.UserPreference) error
	FindByEmailNotPreloaded(email string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("LOWER(email) = LOWER(?)", email).
		First(&user)

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, result.Error
}

func (r *userRepository) FindByID(userID string) (*model.User, error) {
	var user model.User
	result := r.db.Where("id = ?", userID).
		Preload("Goal").
		Preload("UserPreferences.Tag").
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}

func (r *userRepository) SaveUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateUserTransaction(tx *gorm.DB, userID *uuid.UUID, updateFields map[string]interface{}) error {
	return tx.Model(model.User{}).Where("id = ?", userID).Updates(updateFields).Error
}

func (r *userRepository) UpdateUserPreferences(tx *gorm.DB, user *model.User, pref []model.UserPreference) error {
	if err := tx.Where("user_id = ?", user.ID).Delete(&model.UserPreference{}).Error; err != nil {
		return err
	}

	if len(pref) > 0 {
		if err := tx.Create(&pref).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) FindByEmailNotPreloaded(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("LOWER(email) = LOWER(?)", email).
		First(&user)
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, result.Error
}
