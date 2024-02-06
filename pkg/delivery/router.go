package route

import (
	"app-service-com/pkg/delivery/http"

	"github.com/labstack/echo"
)

func InitializeRoute(e *echo.Echo) {

	userHandler := http.NewUserHandler(e)
	fileHandler := http.NewFileHandler(e)

	// handler
	v1 := e.Group("/v1")
	v1.GET("/users", userHandler.Fetch).Name = "v1-get-user"
	v1.POST("/user", userHandler.Store)
	v1.GET("/user/:id", userHandler.Find)
	v1.DELETE("/user/:id", userHandler.Destroy)

	v1.POST("/upload", fileHandler.Upload)
}
