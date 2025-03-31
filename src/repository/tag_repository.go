package repository

import (
	"github.com/E-cercise/E-cercise/src/model"
	"gorm.io/gorm"
)

type TagRepository interface {
	GetAll() ([]model.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db}
}

func (r *tagRepository) GetAll() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}
