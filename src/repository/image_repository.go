package repository

import (
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageRepository interface {
	FindByID(imgID uuid.UUID) (*model.Image, error)
	FindByIDTransaction(tx *gorm.DB, imgID uuid.UUID) (*model.Image, error)
	FindByEquipmentID(equipmentID uuid.UUID) ([]model.Image, error)
	CreateImage(tx *gorm.DB, image *model.Image) error
	SaveImage(tx *gorm.DB, img *model.Image) error
}

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) FindByID(imgID uuid.UUID) (*model.Image, error) {
	var imgData *model.Image

	if err := r.db.Where("id = ?", imgID).First(&imgData).Error; err != nil {
		logger.Log.WithError(err).Error("cant find image ID", imgID)
		return nil, err
	}
	return imgData, nil
}

func (r *imageRepository) FindByIDTransaction(tx *gorm.DB, imgID uuid.UUID) (*model.Image, error) {
	var imgData *model.Image

	if err := tx.Where("id = ?", imgID).First(&imgData).Error; err != nil {
		logger.Log.WithError(err).Error("cant find image ID", imgID)
		return nil, err
	}

	return imgData, nil
}

func (r *imageRepository) FindByEquipmentID(equipmentID uuid.UUID) ([]model.Image, error) {
	var images []model.Image

	if err := r.db.Where("equipment_id = ?", equipmentID).Find(&images).Error; err != nil {
		logger.Log.WithError(err).Error("can't find images for equipment ID", equipmentID)
		return nil, err
	}

	return images, nil
}

func (r *imageRepository) CreateImage(tx *gorm.DB, image *model.Image) error {
	return tx.Create(image).Error
}

func (r *imageRepository) SaveImage(tx *gorm.DB, img *model.Image) error {
	return tx.Save(img).Error
}
