package service

import (
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
)

type TagService interface {
	GetAllTags() ([]model.Tag, error)
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(r repository.TagRepository) TagService {
	return &tagService{repo: r}
}

func (s *tagService) GetAllTags() ([]model.Tag, error) {
	return s.repo.GetAll()
}
