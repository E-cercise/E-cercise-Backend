package service

import (
	"context"
	"fmt"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type ImageService interface {
	UploadImage(context context.Context, file multipart.File, fileHeader *multipart.FileHeader, isPrimary bool) (string, error)
	//GetAllEquipmentData() (*response.EquipmentsResponse, error)
	ArchiveImage(tx *gorm.DB, context context.Context, imgID uuid.UUID, eqpID uuid.UUID) error
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

	fileName := generateFileName(enum.Temp.ToString(), fileHeader)

	cloudinaryPath, err := s.cloudinaryService.UploadImage(context, file, fileHeader, fileName)

	if err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error uploading image to cloudinary")
		return "", err
	}

	newImage := model.Image{
		EquipmentID:    nil,
		IsPrimary:      isPrimary,
		ImgPath:        fileName,
		CloudinaryPath: cloudinaryPath,
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

func generateFileName(folder string, fileheader *multipart.FileHeader) string {
	timestamp := time.Now().Format("20060102150405") // e.g., "20250112094530"
	ext := filepath.Ext(fileheader.Filename)
	return fmt.Sprintf("%s/%s_%s%s", folder, "img", timestamp, ext)
}

func (s *imageService) ArchiveImage(tx *gorm.DB, context context.Context, imgID uuid.UUID, eqpID uuid.UUID) error {
	img, err := s.imageRepo.FindByIDTransaction(tx, imgID)

	if err != nil {
		logger.Log.WithError(err).Error("error finding image ID", imgID)
		return err
	}

	if img.State == enum.Archive {
		logger.Log.Warnf("img ID %v has already archive", imgID)
		return nil
	}

	oldPublicID := img.ImgPath
	newPublicID := strings.ReplaceAll(img.ImgPath, "temp/", fmt.Sprintf("archive/%v/", eqpID))

	img.CloudinaryPath = strings.ReplaceAll(img.CloudinaryPath, "temp/", fmt.Sprintf("archive/%v/", eqpID))
	img.ImgPath = newPublicID

	if err = s.imageRepo.SaveImage(tx, img); err != nil {
		logger.Log.WithError(err).Error("cannot save image in repo", img)
		return err
	}

	err = s.cloudinaryService.MoveImage(context, oldPublicID, newPublicID)
	if err != nil {
		logger.Log.WithError(err).Error("error move image in cloudinary", "err", err.Error())
		return err
	}

	return nil
}
