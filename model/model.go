package model

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

type Account struct {
	gorm.Model
	CoinID  int `gorm:"uniqueIndex:idx_userid_coinid"`
	Coin    Coin
	UserID  int `gorm:"uniqueIndex:idx_userid_coinid"`
	User    User
	Balance uint
}

type Coin struct {
	gorm.Model
	Name string
	Unit string
}

func (u User) String() string {
	return fmt.Sprintf("id:%v, name:%s", u.ID, u.Name)
}

func MigrateAll(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Coin{})
	db.AutoMigrate(&Account{})
}
