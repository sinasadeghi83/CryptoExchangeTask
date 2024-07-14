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

	users := generateFakeUsers()
	if result := db.Create(users); result.Error != nil {
		return result.Error
	}

	coins := generateInitialCoins()
	if result := db.Create(coins); result.Error != nil {
		return result.Error
	}

	return nil
}

func generateFakeUsers() []model.User {
	return []model.User{
		{Name: "sina"},
		{Name: "ali"},
		{Name: "hassan"},
	}
}

func generateInitialCoins() []model.Coin {
	// Using smallest unit of each coin
	return []model.Coin{
		{Name: "USDT", Unit: "cent"},
		{Name: "BTC", Unit: "satoshi"},
		{Name: "ETH", Unit: "wei"},
		{Name: "DOGE", Unit: "shibe"},
		{Name: "XRP", Unit: "drop"},
	}
}
