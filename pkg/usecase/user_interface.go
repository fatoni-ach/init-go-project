package usecase

import (
	"app-service-com/pkg/delivery/validation"
	"app-service-com/pkg/models"
	"net/url"
)

type User interface {
	Fetch(url.Values) ([]*models.User, error)
	Store(*validation.User) (*models.User, error)
	Find(int32) (models.User, error)
	Destroy(int) error
}
