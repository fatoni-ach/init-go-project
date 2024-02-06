package user_repository

import (
	"app-service-com/pkg/models"
	"app-service-com/pkg/repository"
	"app-service-com/services"
	"net/url"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewRepository() repository.User {
	return &gormRepository{
		services.DB,
	}
}

func (repo *gormRepository) Fetch(filter url.Values) ([]*models.User, error) {
	var users []*models.User

	query := repo.db.Select("username", "email", "fullname")
	if filter.Has("email") {
		query.Where("email LIKE ?", "%"+filter.Get("email")+"%")
	}

	if filter.Has("username") {
		query.Where("username LIKE ?", "%"+filter.Get("username")+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		services.WriteLogWarn(err)
		return nil, err
	}

	return users, nil
}

func (repo *gormRepository) Store(user models.User) (models.User, error) {
	err := repo.db.Create(&user).Error

	if err != nil {
		services.RecoverPanic()
		return user, err
	}

	return user, nil
}

func (repo *gormRepository) Find(id int32) (models.User, error) {
	var user models.User
	if err := repo.db.Find(&user, id).Error; err != nil {
		services.RecoverPanic()
		return user, err
	}
	return user, nil
}

func (repo *gormRepository) Destroy(id int) error {
	if err := repo.db.Delete(&models.User{}, id).Error; err != nil {
		services.RecoverPanic()
		return err
	}

	return nil
}
