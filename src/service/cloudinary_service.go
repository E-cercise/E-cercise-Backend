package service

import (
	"context"
	"fmt"
	"github.com/E-cercise/E-cercise/src/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"mime/multipart"
)

type CloudinaryService interface {
	UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
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
	// Upload parameters
	uploadParams := uploader.UploadParams{
		Folder: folder,
	}

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
