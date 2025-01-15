package service

import (
	"context"
	"fmt"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
	"time"
)

type ImageService interface {
	UploadImage(context context.Context, file multipart.File, fileHeader *multipart.FileHeader, isPrimary bool) (string, error)
	//GetAllEquipmentData() (*response.EquipmentsResponse, error)

}

type imageService struct {
	db                *gorm.DB
	imageRepo         repository.ImageRepository
	cloudinaryService CloudinaryService
}

func NewImageService(db *gorm.DB, imageRepo repository.ImageRepository, cloudinaryService CloudinaryService) ImageService {
	return &imageService{db: db, imageRepo: imageRepo, cloudinaryService: cloudinaryService}
}

func (s *imageService) UploadImage(context context.Context, file multipart.File, fileHeader *multipart.FileHeader, isPrimary bool) (string, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	filePath, err := s.cloudinaryService.UploadImage(context, file, fileHeader, enum.Temp.ToString())

	if err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error uploading image to cloudinary")
		return "", err
	}

	newImage := model.Image{
		EquipmentID:    nil,
		IsPrimary:      isPrimary,
		CloudinaryPath: filePath,
		State:          enum.Temp,
	}

	err = s.imageRepo.CreateImage(tx, &newImage)

	if err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error creating image")
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return "", err
	}

	return newImage.ID.String(), nil
}

func generateFileName(fileHeader *multipart.FileHeader, folder string) string {
	timestamp := time.Now().Format("20060102150405") // e.g., "20250112094530"
	extension := filepath.Ext(fileHeader.Filename)   // Get the file extension
	return fmt.Sprintf("%s/%s_%s%s", folder, "img", timestamp, extension)
}
