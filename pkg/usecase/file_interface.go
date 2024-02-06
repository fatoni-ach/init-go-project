package usecase

import (
	"mime/multipart"
)

type File interface {
	Upload(multipart.File, string, int64) error
}
