package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/E-cercise/E-cercise/src/enum"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
	"time"
)

type ImageService interface {
	UploadImage(context context.Context, file multipart.File, fileHeader *multipart.FileHeader, isPrimary bool) (string, string, error)
	ArchiveImage(tx *gorm.DB, context context.Context, imgID uuid.UUID, eqpID uuid.UUID, eqOptID uuid.UUID, isPrimary bool) error
	DeleteImage(tx *gorm.DB, context context.Context, imgID uuid.UUID) error
}

type imageService struct {
	db                *gorm.DB
	imageRepo         repository.ImageRepository
	cloudinaryService CloudinaryService
}

func NewImageService(db *gorm.DB, imageRepo repository.ImageRepository, cloudinaryService CloudinaryService) ImageService {
	return &imageService{db: db, imageRepo: imageRepo, cloudinaryService: cloudinaryService}
}

func (s *imageService) UploadImage(context context.Context, file multipart.File, fileHeader *multipart.FileHeader, isPrimary bool) (string, string, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	fileName := generateFileName(enum.Temp.ToString())

	cloudinaryPath, err := s.cloudinaryService.UploadImage(context, file, fileHeader, fileName)

	if err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error uploading image to cloudinary")
		return "", "", err
	}

	newImage := model.Image{
		EquipmentOptionID: nil,
		IsPrimary:         isPrimary,
		ImgPath:           fileName,
		CloudinaryPath:    cloudinaryPath,
		State:             enum.Temp,
	}

	err = s.imageRepo.CreateImage(tx, &newImage)

	if err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error creating image")
		return "", "", err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return "", "", err
	}

	return newImage.ID.String(), cloudinaryPath, nil
}

func generateFileName(folder string) string {
	timestamp := time.Now().Format("20060102150405") // e.g., "20250112094530"
	someRandString := helper.RandomString(5)
	return fmt.Sprintf("%s/%s_%s%s", folder, "img", someRandString, timestamp)
}

func (s *imageService) ArchiveImage(tx *gorm.DB, context context.Context, imgID uuid.UUID, eqpID uuid.UUID, eqOptID uuid.UUID, isPrimary bool) error {
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
	newPublicID := strings.ReplaceAll(img.ImgPath, "temp/", fmt.Sprintf("archive/%v/%v/", eqpID, eqOptID))

	img.CloudinaryPath = strings.ReplaceAll(img.CloudinaryPath, "temp/", fmt.Sprintf("archive/%v/%v/", eqpID, eqOptID))
	img.ImgPath = newPublicID
	img.IsPrimary = isPrimary
	img.EquipmentOptionID = &eqOptID

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

func (s *imageService) DeleteImage(tx *gorm.DB, ctx context.Context, imgID uuid.UUID) error {
	logger.Log.Infof("🗑 Attempting to delete image ID: %s", imgID)

	img, err := s.imageRepo.FindByIDTransaction(tx, imgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warnf("⚠️ Image ID %s not found — skipping delete", imgID)
			return nil
		}
		return fmt.Errorf("❌ error finding image in DB: %w", err)
	}

	imgPath := img.ImgPath

	if err := s.imageRepo.DeleteImage(tx, imgID); err != nil {
		return fmt.Errorf("❌ error deleting image in DB: %w", err)
	}
	logger.Log.Infof("✅ Deleted image from DB: %s", imgID)

	if err := s.cloudinaryService.DeleteImage(ctx, imgPath); err != nil {
		return fmt.Errorf("❌ error deleting from Cloudinary (%s): %w", imgPath, err)
	}
	logger.Log.Infof("✅ Deleted image from Cloudinary: %s", imgPath)

	return nil
}
