package main

import (
	"app-service-com/config"
	"app-service-com/helpers"
	route "app-service-com/pkg/delivery"
	"app-service-com/pkg/transformer"
	"app-service-com/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	log "github.com/sirupsen/logrus"
)

func init() {
	if config.GetBoolean(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func configureLogging() {
	lLevel := config.Get("log.level")
	fmt.Println("Setting log level to ", lLevel)
	switch strings.ToUpper(lLevel) {
	default:
		fmt.Println("Unknown level [", lLevel, "]. Log level set to ERROR")
		log.SetLevel(log.ErrorLevel)
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	}

	lType := config.Get("log.type")
	fmt.Println("Setting log type to ", lType)
	if strings.ToUpper(lType) == "FILE" {
		helpers.InitLogRotate()
		logFile := helpers.GetFileLog()
		if logFile != nil {
			log.SetOutput(logFile)
		}
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	c.JSON(code, transformer.ResponseFailed{
		Status:  code,
		Message: err.Error(),
		Error:   err,
	})
}

func main() {
	// runtime.GOMAXPROCS(2)
	configureLogging()

	dbHost := config.Get(`database.host`)
	dbPort := config.Get(`database.port`)
	dbUser := config.Get(`database.user`)
	dbPass := config.Get(`database.pass`)
	dbName := config.Get(`database.name`)

	services.OpenDBConnection(dbUser, dbPass, dbHost, dbPort, dbName)

	defer services.CloseDBConnection()
	defer services.RecoverPanic()

	address := "localhost" + config.Get("server.address")

	e := echo.New()
	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	// Initialize Validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Error Http Handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	route.InitializeRoute(e)
	e.Logger.Fatal(e.Start(address))
}
