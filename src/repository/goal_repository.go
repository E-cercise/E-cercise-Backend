package repository

import (
	"github.com/E-cercise/E-cercise/src/model"
	"gorm.io/gorm"
)

type GoalRepository interface {
	FindAll() ([]model.Goal, error)
}

type goalRepository struct {
	db *gorm.DB
}

func NewGoalRepository(db *gorm.DB) GoalRepository {
	return &goalRepository{db}
}

func (r *goalRepository) FindAll() ([]model.Goal, error) {
	var goals []model.Goal
	err := r.db.Find(&goals).Error
	return goals, err
}
