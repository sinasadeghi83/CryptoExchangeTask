package db

import (
	"github.com/sinasadeghi83/SwapTask/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewConn() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("task.db"), &gorm.Config{})
}

func SetupDB() (*gorm.DB, error) {
	db, err := NewConn()
	if err != nil {
		return db, err
	}

	model.MigrateAll(db)
	return db, err
}

func CreateDummyData() error {
	db, err := NewConn()
	if err != nil {
		return err
	}

	users := []model.User{
		{Name: "sina"},
		{Name: "ali"},
		{Name: "hassan"},
	}

	result := db.Create(users)

	return result.Error
}
