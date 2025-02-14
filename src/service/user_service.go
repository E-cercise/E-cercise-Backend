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
	"strings"
)

type UserService interface {
	RegisterUser(reqBody request.RegisterRequest) error
	LoginUser(reqBody request.LoginRequest) (*string, error)
}

type userService struct {
	db       *gorm.DB
	userRepo repository.UserRepository
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository) UserService {
	return &userService{db: db, userRepo: userRepo}
}

func (s *userService) RegisterUser(reqBody request.RegisterRequest) error {
	existingUser, err := s.userRepo.FindByEmail(reqBody.Email)
	if existingUser != nil || err != nil {
		return errors.New("email already exists")
	}

	password, err := helper.EncryptPassword(reqBody.Password)
	if err != nil {
		return errors.New("failed to encrypt password")
	}

	newUser := model.User{
		Email:       reqBody.Email,
		Password:    password,
		FirstName:   reqBody.FirstName,
		LastName:    reqBody.LastName,
		Address:     reqBody.Address,
		PhoneNumber: reqBody.PhoneNumber,
	}

	err = s.userRepo.CreateUser(&newUser)
	if err != nil {
		logger.Log.WithError(err).Error("failed to create user")
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil

}

func (s *userService) LoginUser(reqBody request.LoginRequest) (*string, error) {

	user, err := s.userRepo.FindByEmail(strings.ToLower(reqBody.Email))

	if user == nil && err == nil {
		return nil, errors.New(fmt.Sprintf("Email %v does not exist", reqBody.Email))
	}

	valid := helper.ComparePassword(reqBody.Password, user.Password)

	if valid != true {
		return nil, errors.New("invalid password")
	}

	token, err := helper.CreateToken(user.ID, user.FirstName, user.LastName, user.Role)

	if err != nil {
		return nil, errors.New("failed to create token, JWT Error")
	}

	return &token, nil

}
