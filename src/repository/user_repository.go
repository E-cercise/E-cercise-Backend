package repository

import (
	"github.com/E-cercise/E-cercise/src/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(userID string) (*model.User, error)
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
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &user, nil
}
