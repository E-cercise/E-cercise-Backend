package service

import (
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
)

type UserPreferenceService interface {
	SetUserPreferences(userID string, tagIDs []string) error
	GetUserPreferences(userID string) ([]model.Tag, error)
}

type userPrefService struct {
	repo repository.UserPreferenceRepository
}

func NewUserPreferenceService(r repository.UserPreferenceRepository) UserPreferenceService {
	return &userPrefService{repo: r}
}

func (s *userPrefService) SetUserPreferences(userID string, tagIDs []string) error {
	return s.repo.SetPreferences(userID, tagIDs)
}

func (s *userPrefService) GetUserPreferences(userID string) ([]model.Tag, error) {
	return s.repo.GetPreferences(userID)
}
