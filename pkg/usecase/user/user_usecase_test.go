package user_usecase

import (
	"app-service-com/pkg/delivery/validation"
	"app-service-com/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tkuchiki/faketime"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) Store(user models.User) (models.User, error) {
	arguments := repository.Mock.Called(user)

	result := arguments.Get(0)

	return result.(models.User), arguments.Error(1)
}

func (repository *UserRepositoryMock) Fetch(filter url.Values) ([]*models.User, error) {

	args := repository.Mock.Called(filter)

	result := args.Get(0).([]*models.User)
	err := args.Error(1)

	return result, err
}

func (repository *UserRepositoryMock) Find(id int32) (models.User, error) {

	// fmt.Println(id)
	return models.User{}, nil
}

func (repository *UserRepositoryMock) Destroy(id int) error {

	fmt.Println(id)
	return nil
}

func TestUseUsecase_Store(t *testing.T) {

	t.Run("test-store-success", func(t *testing.T) {
		f := faketime.NewFaketime(2010, time.February, 10, 23, 0, 0, 0, time.Local)
		defer f.Undo()
		f.Do()

		userValidation := validation.User{
			Username: "superuser",
			Email:    "superuser@example.com",
			Password: "Password1",
			Fullname: "superuser",
			Gender:   true,
		}

		now := time.Now()
		userModel := models.User{
			Fullname: "superuser",
			Gender:   true,
			Username: "superuser",
			Email:    "superuser@example.com",
			Password: "Password1",
			CreatedAt: sql.NullTime{
				Valid: true,
				Time:  now,
			},
			UpdatedAt: sql.NullTime{
				Valid: true,
				Time:  now,
			},
		}

		mockRepo := new(UserRepositoryMock)

		mockRepo.Mock.On("Store", userModel).Return(userModel, nil)
		testUsecase := NewUseCase(mockRepo)

		user, err := testUsecase.Store(&userValidation)

		mockRepo.Mock.AssertExpectations(t)
		assert.Nil(t, err)
		assert.NotNil(t, user)
	})

	t.Run("test-store-failed", func(t *testing.T) {
		f := faketime.NewFaketime(2010, time.February, 10, 23, 0, 0, 0, time.Local)
		defer f.Undo()
		f.Do()

		userValidation := validation.User{}

		now := time.Now()
		userModel := models.User{
			CreatedAt: sql.NullTime{
				Valid: true,
				Time:  now,
			},
			UpdatedAt: sql.NullTime{
				Valid: true,
				Time:  now,
			},
		}

		mockRepo := new(UserRepositoryMock)

		mockRepo.Mock.On("Store", userModel).Return(models.User{}, errors.New("Failed Store User"))
		testUsecase := NewUseCase(mockRepo)

		result, err := testUsecase.Store(&userValidation)

		mockRepo.Mock.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Equal(t, "Failed Store User", err.Error())
		assert.Equal(t, &models.User{}, result)
	})

}

func TestUserUsecase_Fetch(t *testing.T) {
	t.Run("test-fetch-success", func(t *testing.T) {
		f := faketime.NewFaketime(2010, time.February, 10, 0, 0, 0, 0, time.Local)
		defer f.Undo()
		f.Do()

		now := time.Now()
		expectedResult := []*models.User{
			{
				ID:       1,
				Username: "superuser",
				Email:    "superuser@example.com",
				Fullname: "Super User",
				CreatedAt: sql.NullTime{
					Valid: true,
					Time:  now,
				},
				UpdatedAt: sql.NullTime{
					Valid: true,
					Time:  now,
				},
				Password: "password111",
			},
			{
				ID:       2,
				Username: "superuser-2",
				Email:    "superuser2@example.com",
				Fullname: "Super User 2",
				CreatedAt: sql.NullTime{
					Valid: true,
					Time:  now,
				},
				UpdatedAt: sql.NullTime{
					Valid: true,
					Time:  now,
				},
				Password: "password1111",
			},
		}

		userRepoMock := new(UserRepositoryMock)
		userUsecase := NewUseCase(userRepoMock)

		userRepoMock.Mock.On("Fetch", url.Values{}).Return(expectedResult, nil)

		result, err := userUsecase.Fetch(url.Values{})

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("test-fetch-failed", func(t *testing.T) {
		expectedResult := []*models.User{}

		userRepoMock := new(UserRepositoryMock)
		userUsecase := NewUseCase(userRepoMock)

		userRepoMock.Mock.On("Fetch", url.Values{}).Return(expectedResult, errors.New("Failed Fetch User"))

		result, err := userUsecase.Fetch(url.Values{})

		assert.NotNil(t, err)
		assert.Equal(t, "Failed Fetch User", err.Error())
		assert.Equal(t, expectedResult, result)
	})
}
