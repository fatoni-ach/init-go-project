package user_repository

import (
	"app-service-com/pkg/models"
	"database/sql"
	"errors"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestFetchSuccess(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(dialector, &gorm.Config{})

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "fullname"}).
		// AddRows([][]driver.Value{
		// 	{"1", "superuser", "superuser@example.com", "superuser"},
		// 	{"2", "user-1", "user1@example.com", "user 1"},
		// }...)
		AddRow("1", "superuser", "superuser@example.com", "superuser").
		AddRow("2", "user-1", "user1@example.com", "user 1")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id","username","email","fullname" FROM "users" WHERE "users"."deleted_at" IS NULL`)).
		WillReturnRows(rows)

	gormRepo := GormRepository{db: db}

	user1 := models.User{
		ID:       1,
		Email:    "superuser@example.com",
		Username: "superuser",
		Fullname: "superuser",
	}
	user2 := models.User{
		ID:       2,
		Email:    "user1@example.com",
		Username: "user-1",
		Fullname: "user 1",
	}

	expectedResult := []*models.User{&user1, &user2}
	results, err := gormRepo.Fetch(url.Values{})

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, results)
}

func TestFetchFailed(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(dialector, &gorm.Config{})

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "fullname"}).
		AddRow("", "", "", "")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "id","username","email","fullname" FROM "users" WHERE "users"."deleted_at" IS NULL`)).
		WillReturnRows(rows).
		WillReturnError(errors.New("Failed Fetch Users"))

	gormRepo := GormRepository{db: db}

	expectedResult := []*models.User(nil)
	results, err := gormRepo.Fetch(url.Values{})

	assert.Equal(t, "Failed Fetch Users", err.Error())
	assert.Equal(t, expectedResult, results)
}

func TestStoreSuccess(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(dialector, &gorm.Config{})

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "email", "username", "fullname", "password", "gender", "created_at", "updated_at"}).AddRow("1", "superuser@example.com", "superuser", "superuser", "password111", true, now, now)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(rows)
	mock.ExpectCommit()
	gormRepo := GormRepository{db: db}

	expectedResult := models.User{
		ID:       1,
		Email:    "superuser@example.com",
		Username: "superuser",
		Fullname: "superuser",
		Password: "password111",
		Gender:   true,
		CreatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}
	result, err := gormRepo.Store(models.User{})

	assert.Equal(t, expectedResult, result)
	assert.Equal(t, nil, err)
}

func TestStoreFailed(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(dialector, &gorm.Config{})

	now := time.Now()
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "email", "username", "fullname", "password", "gender", "created_at", "updated_at"}).AddRow("1", "superuser@example.com", "superuser", "superuser", "password111", true, now, now)
	mock.ExpectQuery("INSERT").WillReturnError(errors.New("Failed Store Data")).WillReturnRows(rows)
	mock.ExpectRollback()
	gormRepo := GormRepository{db: db}

	expectedError := "Failed Store Data"
	_, err := gormRepo.Store(models.User{})

	assert.Equal(t, expectedError, err.Error())
}

func TestFindSuccess(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(dialector, &gorm.Config{})

	now := time.Now()
	var id int32 = 1

	rows := sqlmock.NewRows([]string{"id", "email", "username", "fullname", "password", "gender", "created_at", "updated_at"}).AddRow("1", "superuser@example.com", "superuser", "superuser", "password111", true, now, now)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL`)).WithArgs(id).WillReturnRows(rows)

	expectedResult := models.User{
		ID:       1,
		Email:    "superuser@example.com",
		Username: "superuser",
		Fullname: "superuser",
		Password: "password111",
		Gender:   true,
		CreatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	gormRepo := GormRepository{db: db}
	result, err := gormRepo.Find(1)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, result)
}

func TestFindFailed(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	db, _ := gorm.Open(dialector, &gorm.Config{})

	now := time.Now()
	var id int32 = 1

	rows := sqlmock.NewRows([]string{"id", "email", "username", "fullname", "password", "gender", "created_at", "updated_at"}).AddRow("", "", "", "", "", false, now, now)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL`)).WithArgs(id).WillReturnRows(rows).WillReturnError(errors.New("Failed Find User"))

	expectedResult := models.User{}

	gormRepo := GormRepository{db: db}
	result, err := gormRepo.Find(1)

	assert.Equal(t, "Failed Find User", err.Error())
	assert.Equal(t, expectedResult, result)
}
