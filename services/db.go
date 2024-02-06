package services

import (
	"app-service-com/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// Mysql gorm
	// _ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB variable global db
var DB *gorm.DB

// OpenDBConnection Open connection Database
func OpenDBConnection(dbUser string, dbPass string, dbHost string, dbPort string, dbName string) {

	// connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	// var err error
	// DB, err = gorm.Open("mysql", connectionString)
	// if err != nil {
	// 	RecoverPanic()
	// }

	// DB.DB().SetMaxIdleConns(0)

	// if config.GetBoolean(`debug`) {
	// DB.LogMode(true)
	// }

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPass, dbName, dbPort)
	// dsn := "host=localhost user=root password=example dbname=asset_management port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		RecoverPanic()
	}
	sqlDB, err := DB.DB()

	if err != nil {
		RecoverPanic()
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if config.GetBoolean(`debug`) {
		DB.Debug()
	}
}

// CloseDBConnection Closing DB Connection
func CloseDBConnection() error {
	// var err error
	// var sqlDB sql.
	sqlDB, err := DB.DB()
	if err != nil {
		RecoverPanic()
	}
	return sqlDB.Close()
}
