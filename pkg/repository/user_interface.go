package repository

import (
	"app-service-com/pkg/models"
	"net/url"
)

type User interface {
	Fetch(url.Values) ([]*models.User, error)
	Store(models.User) (models.User, error)
	Find(int32) (models.User, error)
	Destroy(int) error
}
