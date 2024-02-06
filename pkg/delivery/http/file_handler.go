package http

import (
	"app-service-com/pkg/transformer"
	"app-service-com/pkg/usecase"
	fileUsecase "app-service-com/pkg/usecase/file"

	"github.com/labstack/echo"
)

type FileHandler struct {
	FUsecase usecase.File
}

func NewFileHandler(e *echo.Echo) *FileHandler {
	us := fileUsecase.NewFileUseCase()

	return &FileHandler{FUsecase: us}
}

func (uf *FileHandler) Upload(c echo.Context) error {
	file, _ := c.FormFile("file")

	filename := file.Filename
	size := file.Size

	src, _ := file.Open()
	defer src.Close()

	if err := uf.FUsecase.Upload(src, filename, size); err != nil {
		return err
	}

	var response transformer.ResponseSuccess
	transformer.TransformResponse(&response, "Success Upload File")

	return c.JSON(response.Status, response)
}
