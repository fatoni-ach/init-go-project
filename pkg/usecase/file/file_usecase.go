package file_usecase

import (
	uc "app-service-com/pkg/usecase"
	"app-service-com/services"
	"mime/multipart"
)

type usecase struct {
}

func NewFileUseCase() uc.File {
	return &usecase{}
}

func (uc *usecase) Upload(file multipart.File, filename string, size int64) error {
	if err := services.UploadS3(file, filename, size); err != nil {
		return err
	}
	return nil
}
