package service

import (
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
)

type GoalService interface {
	GetAllGoal() ([]model.Goal, error)
}

type goalService struct {
	repo repository.GoalRepository
}

func NewGoalService(r repository.GoalRepository) GoalService {
	return &goalService{repo: r}
}

func (s *goalService) GetAllGoal() ([]model.Goal, error) {
	return s.repo.FindAll()
}
