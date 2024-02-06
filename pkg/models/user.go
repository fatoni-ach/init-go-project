package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"column:username"`
	Email     string
	Password  string
	Fullname  string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	DeletedAt gorm.DeletedAt
	Gender    bool
}

// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	services.WriteLog("Log Before create")
// 	services.WriteLog(u)
// 	return nil
// }

// func (u *User) AfterCreate(tx *gorm.DB) (err error) {
// 	services.WriteLog("Log Aafter create")
// 	services.WriteLog(u)
// 	return nil
// }
