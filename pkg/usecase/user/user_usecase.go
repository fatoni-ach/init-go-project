package user_usecase

import (
	"app-service-com/pkg/delivery/validation"
	"app-service-com/pkg/models"
	"app-service-com/pkg/repository"
	uc "app-service-com/pkg/usecase"
	"app-service-com/services"
	"net/url"
	"time"
)

type usecase struct {
	repo repository.User
}

func NewUseCase(repo repository.User) uc.User {
	return &usecase{
		repo: repo,
	}
}

func (uc *usecase) Fetch(filter url.Values) ([]*models.User, error) {
	users, err := uc.repo.Fetch(filter)
	if err != nil {
		services.WriteLog(err)
		return []*models.User{}, err
	}

	return users, nil
}

func (uc *usecase) Store(userValidation *validation.User) (*models.User, error) {
	var user models.User

	user.Email = userValidation.Email
	user.Username = userValidation.Username
	user.Fullname = userValidation.Fullname
	user.Password = userValidation.Password
	user.Gender = userValidation.Gender
	user.CreatedAt.Time = time.Now()
	user.CreatedAt.Valid = true
	user.UpdatedAt.Time = time.Now()
	user.UpdatedAt.Valid = true

	user, err := uc.repo.Store(user)
	if err != nil {
		services.WriteLog(err)
		return &models.User{}, err
	}
	return &user, nil
}

func (uc *usecase) Find(id int32) (models.User, error) {
	var user models.User

	user, err := uc.repo.Find(id)
	if err != nil {
		services.WriteLog(err)
		return user, err
	}
	return user, nil
}

func (uc *usecase) Destroy(id int) error {
	if err := uc.repo.Destroy(id); err != nil {
		return err
	}
	return nil
}
