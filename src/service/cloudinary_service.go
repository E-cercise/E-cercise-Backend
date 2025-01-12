package service

import (
	"context"
	"fmt"
	"github.com/E-cercise/E-cercise/src/config"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"mime/multipart"
)

type CloudinaryService interface {
	UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
	MoveImage(ctx context.Context, publicID, newFolder string) (string, error)
}

type cloudinaryService struct {
	cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryService() (CloudinaryService, error) {

	if config.CloudinaryCloudName == "" || config.CloudinaryApiKey == "" || config.CloudinaryApiSecret == "" {
		return nil, fmt.Errorf("missing Cloudinary environment variables")
	}

	// Initialize Cloudinary instance
	cld, err := cloudinary.NewFromParams(config.CloudinaryCloudName, config.CloudinaryApiKey, config.CloudinaryApiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	return &cloudinaryService{cloudinary: cld}, nil
}

// UploadImage uploads an image file to Cloudinary and returns the URL
func (s *cloudinaryService) UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Validate file type
	allowedTypes := []string{"image/jpeg", "image/png"}
	if err := validateFileType(fileHeader, allowedTypes); err != nil {
		return "", fmt.Errorf("file validation failed: %v", err)
	}

	// Validate file size
	const maxFileSize = 5 * 1024 * 1024 // 5 MB
	if err := validateFileSize(fileHeader, maxFileSize); err != nil {
		return "", fmt.Errorf("file validation failed: %v", err)
	}

	// Generate file name
	fileName := generateFileName(fileHeader, folder)

	// Upload parameters
	uploadParams := uploader.UploadParams{
		Folder:   folder,
		PublicID: fileName, // Use the generated file name
	}

	// Log the file name
	logger.Log.Infof("Uploading file: %s", fileHeader.Filename)

	// Upload the file
	resp, err := s.cloudinary.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}

	return resp.SecureURL, nil
}

// DeleteImage deletes an image from Cloudinary by public ID
func (s *cloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
	_, err := s.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %v", err)
	}

	return nil
}

// MoveImage moves an image from one folder to another on Cloudinary
func (s *cloudinaryService) MoveImage(ctx context.Context, publicID, newFolder string) (string, error) {
	newPublicID := fmt.Sprintf("%s/%s", newFolder, publicID)

	_, err := s.cloudinary.Upload.Rename(ctx, uploader.RenameParams{
		FromPublicID: publicID,
		ToPublicID:   newPublicID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to move image: %v", err)
	}

	return newPublicID, nil
}

func (s *cloudinaryService) SaveImage(ctx context.Context, imgPath string, equipmentID uuid.UUID) (string, error) {
	// Move the image from temp to archive folder
	newFolder := fmt.Sprintf("archive/%s", equipmentID)
	newPublicID, err := s.MoveImage(ctx, imgPath, newFolder)
	if err != nil {
		return "", err
	}

	//// Update image record in the database
	//image.EquipmentID = equipmentID
	//image.CloudinaryPath = newPublicID
	//image.State = enum.Archive
	//err = s.db.Save(&image).Error // Assumes s.db is your GORM DB instance
	//if err != nil {
	//	return "", fmt.Errorf("failed to update image record: %v", err)
	//}

	return newPublicID, nil
}

func validateFileType(fileHeader *multipart.FileHeader, allowedTypes []string) error {
	contentType := fileHeader.Header.Get("Content-Type")
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return nil
		}
	}
	return fmt.Errorf("invalid file type: %s", contentType)
}

func validateFileSize(fileHeader *multipart.FileHeader, maxSize int64) error {
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size exceeds limit: %d bytes", maxSize)
	}
	return nil
}
