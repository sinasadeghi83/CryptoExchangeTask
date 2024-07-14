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
		{Name: "USDT", Unit: "cent-tether", UnitFactor: 100},
		{Name: "BTC", Unit: "satoshi", UnitFactor: 100000},
		{Name: "ETH", Unit: "wei", UnitFactor: 1000000000000000000},
		{Name: "DOGE", Unit: "shibe", UnitFactor: 100000},
		{Name: "XRP", Unit: "drop", UnitFactor: 1000000},
		{Name: "USD", Unit: "cent", UnitFactor: 100},
	}
}
