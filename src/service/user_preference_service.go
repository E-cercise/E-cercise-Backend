package service

import (
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/google/uuid"
)

type UserPreferenceService interface {
	SetUserPreferences(userID uuid.UUID, tagIDs []uuid.UUID) error
	GetUserPreferences(userID uuid.UUID) ([]model.Tag, error)
}

type userPrefService struct {
	repo repository.UserPreferenceRepository
}

func NewUserPreferenceService(r repository.UserPreferenceRepository) UserPreferenceService {
	return &userPrefService{repo: r}
}

func (s *userPrefService) SetUserPreferences(userID uuid.UUID, tagIDs []uuid.UUID) error {
	return s.repo.SetPreferences(userID, tagIDs)
}

func (s *userPrefService) GetUserPreferences(userID uuid.UUID) ([]model.Tag, error) {
	return s.repo.GetPreferences(userID)
}
