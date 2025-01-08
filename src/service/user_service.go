package service

import (
	"errors"
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"gorm.io/gorm"
)

type UserService interface {
	RegisterUser(reqBody request.RegisterRequest) (*model.User, error)
}

type userService struct {
	db       *gorm.DB
	userRepo repository.UserRepository
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository) UserService {
	return &userService{db: db, userRepo: userRepo}
}

func (s *userService) RegisterUser(reqBody request.RegisterRequest) (*model.User, error) {
	existingUser, err := s.userRepo.FindByEmail(reqBody.Email)
	if existingUser != nil || err != nil {
		return nil, errors.New("email already exists")
	}

	password, err := helper.EncryptPassword(reqBody.Password)
	if err != nil {
		return nil, errors.New("failed to encrypt password")
	}

	newUser := model.User{
		Email:    reqBody.Email,
		Password: password,
	}

	err = s.userRepo.CreateUser(&newUser)
	if err != nil {
		logger.Log.WithError(err).Error("failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil

}
