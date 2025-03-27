package service

import (
	"errors"
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/data/response"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type UserService interface {
	RegisterUser(reqBody request.RegisterRequest) (*model.User, error)
	LoginUser(reqBody request.LoginRequest) (*string, error)
	GetUserProfile(user *model.User) response.UserProfileResponse
	UpdateUserProfile(user *model.User, req request.UpdateUserProfileRequest) error
	SaveUserPreferences(userID uuid.UUID, tagIDs []uuid.UUID) error
}

type userService struct {
	db           *gorm.DB
	userRepo     repository.UserRepository
	userPrefRepo repository.UserPreferenceRepository
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository, userPrefRepo repository.UserPreferenceRepository) UserService {
	return &userService{db: db, userRepo: userRepo, userPrefRepo: userPrefRepo}
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
		Email:       reqBody.Email,
		Password:    password,
		FirstName:   reqBody.FirstName,
		LastName:    reqBody.LastName,
		Address:     reqBody.Address,
		PhoneNumber: reqBody.PhoneNumber,
		Weight:      reqBody.Weight,
		Height:      reqBody.Height,
		Experience:  reqBody.Experience,
		GoalID:      reqBody.GoalID,
	}

	err = s.userRepo.CreateUser(&newUser)
	if err != nil {
		logger.Log.WithError(err).Error("failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil

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

func (s *userService) GetUserProfile(user *model.User) response.UserProfileResponse {

	var preferences []response.PrefResponse

	for _, preference := range user.UserPreferences {

		preferences = append(preferences, response.PrefResponse{
			ID:   preference.Tag.ID,
			Name: preference.Tag.Name,
		})
	}

	res := response.UserProfileResponse{
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Weight:      user.Weight,
		Height:      user.Height,
		Experience:  user.Experience,
		Goal: &response.GoalResponse{
			ID:   user.Goal.ID,
			Name: user.Goal.Name,
		},
		Preferences: preferences,
	}

	return res
}

func (s *userService) UpdateUserProfile(user *model.User, req request.UpdateUserProfileRequest) error {

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("error ", r)
			tx.Rollback()
		}
	}()

	if req.Email != nil {
		existingUser, err := s.userRepo.FindByEmail(*req.Email)
		if existingUser != nil || err != nil {
			return errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.Address != nil {
		user.Address = *req.Address
	}

	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}

	if req.Weight != nil {
		user.Weight = req.Weight
	}
	if req.Height != nil {
		user.Height = req.Height
	}
	if req.Experience != nil {
		user.Experience = req.Experience
	}
	if req.GoalID != nil {
		user.GoalID = req.GoalID
	}

	if req.Preferences != nil {
		var newPrefs []model.UserPreference
		for _, tagID := range req.Preferences {
			newPrefs = append(newPrefs, model.UserPreference{
				UserID: user.ID,
				TagID:  tagID,
			})
		}
		if err := tx.Model(user).
			Association("UserPreferences").
			Replace(newPrefs); err != nil {
			return err
		}
	}

	err := s.userRepo.SaveUserTransaction(tx, user)
	if err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("failed to update user profile")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *userService) SaveUserPreferences(userID uuid.UUID, tagIDs []uuid.UUID) error {
	err := s.userPrefRepo.SetPreferences(userID, tagIDs)
	if err != nil {
		logger.Log.WithError(err).Error("failed to save user preferences")
		return err
	}
	return nil
}
