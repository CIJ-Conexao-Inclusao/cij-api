package service

import (
	"cij_api/src/config"
	"cij_api/src/integration"
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type filesService struct {
	cloudinaryIntegration *cloudinary.Cloudinary
}

func NewFilesService() *filesService {
	config, err := config.LoadCloudinaryConfig(".")
	if err != nil {
		panic("failed to load cloudinary config")
	}

	cloudinaryIntegration := integration.CloudinaryConnect(config.CloudinaryUrl)

	return &filesService{
		cloudinaryIntegration: cloudinaryIntegration,
	}
}

func (f *filesService) UploadFile(file multipart.File, filePath string) (string, error) {
	ctx := context.Background()

	uploadResult, err := f.cloudinaryIntegration.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			PublicID: filePath,
		},
	)

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
