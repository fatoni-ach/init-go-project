package http

import (
	"app-service-com/pkg/delivery/validation"
	"app-service-com/pkg/models"
	userRepo "app-service-com/pkg/repository/user"
	"app-service-com/pkg/transformer"
	"app-service-com/pkg/usecase"
	userUsecase "app-service-com/pkg/usecase/user"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type UserHandler struct {
	AUsecase usecase.User
}

func NewUserHandler(e *echo.Echo) *UserHandler {
	ur := userRepo.NewRepository()
	us := userUsecase.NewUseCase(ur)

	return &UserHandler{AUsecase: us}
}

func (a *UserHandler) Fetch(c echo.Context) error {

	filter := c.QueryParams()

	users, err := a.AUsecase.Fetch(filter)
	if err != nil {
		return err
	}

	var results []transformer.User
	for _, user := range users {
		var temp transformer.User
		transformer.TransformUser(&temp, *user)
		results = append(results, temp)
	}

	var response transformer.ResponseSuccess
	transformer.TransformResponse(&response, results)

	return c.JSON(int(response.Status), response)
}

func (a *UserHandler) Store(c echo.Context) error {
	var userValidation validation.User

	if err := c.Bind(&userValidation); err != nil {
		return err
	}

	if err := c.Validate(userValidation); err != nil {
		return err
	}

	user, err := a.AUsecase.Store(&userValidation)

	if err != nil {
		var response transformer.ResponseFailed
		transformer.TransformResponseFailed(&response, http.StatusBadRequest, "Failed Store Data", err.Error())
		return c.JSON(int(response.Status), response)
	}

	var userTransformer transformer.User
	transformer.TransformUser(&userTransformer, *user)

	var response transformer.ResponseSuccess
	transformer.TransformResponse(&response, userTransformer)

	return c.JSON(int(response.Status), response)
}

func (a *UserHandler) Find(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	user, err := a.AUsecase.Find(int32(id))
	if err != nil {
		return err
	}

	var userTransformer transformer.User
	transformer.TransformUser(&userTransformer, user)

	var response transformer.ResponseSuccess
	transformer.TransformResponse(&response, userTransformer)

	if user == (models.User{}) {
		response.Data = "User Not Found"
	}

	return c.JSON(int(response.Status), response)
}

func (a *UserHandler) Destroy(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := a.AUsecase.Destroy(id); err != nil {
		return err
	}

	var response transformer.ResponseSuccess
	transformer.TransformResponse(&response, "User Success Deleted")
	return c.JSON(response.Status, response)
}
