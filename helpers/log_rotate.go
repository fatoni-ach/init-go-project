package helpers

import (
	"app-service-com/config"
	"fmt"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// LOG ROTATE ALGORITHM
// 1. Check if ENV for log file active
// 2. Create directory logs
// 3. Create log file "hansip.log" inside .logs directory

var (
	fileLogs      *rotatelogs.RotateLogs
	fileLogFormat = "/go_log_.%Y-%m-%d.log"
	fileLogMaxAge = config.GetInt("log.max.age")
	logsPath      string
	err           error
)

// InitLogRotate ...
func InitLogRotate() {

	fmt.Println("Starting LOG Rotate")

	//Get the base file dir
	baseDir, _ := os.Getwd()
	fmt.Println("Base directory: ", baseDir)

	//Creating logs directory
	logsPath = filepath.Join(baseDir, config.Get("log.path"))
	os.MkdirAll(logsPath, 0755)
	fmt.Println("Log directory: ", logsPath)

	//Creating log file "hansip.log"
	LogFilePath := logsPath + fileLogFormat
	fileLogs, err = rotatelogs.New(LogFilePath,
		rotatelogs.WithMaxAge(time.Duration(fileLogMaxAge)*24*time.Hour))
	if err != nil {
		fmt.Println(err)
	}

}

// GetFileLog ...
func GetFileLog() *rotatelogs.RotateLogs {
	return fileLogs
}
