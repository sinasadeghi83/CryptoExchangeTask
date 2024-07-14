package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	BaseModel
	Name string `json:"name"`
}

type Account struct {
	BaseModel
	CoinID  int  `gorm:"uniqueIndex:idx_userid_coinid" json:"coin_id"`
	Coin    Coin `json:"coin"`
	UserID  int  `gorm:"uniqueIndex:idx_userid_coinid" json:"user_id"`
	User    User `json:"user"`
	Balance uint `json:"balance"`
}

type Coin struct {
	BaseModel
	Name string `json:"name"`
	Unit string `json:"unit"`
}

func (account Account) CalculateUSCent() (uint, error) {
	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/generateAvg?fsym=%s&tsym=USD&e=coinbase", account.Coin.Name)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	data := struct {
		RAW struct {
			PRICE float64
		}
	}{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return 0, err
	}

	centPrice := uint(data.RAW.PRICE * 100)

	return centPrice * account.Balance, nil
}

func (u User) String() string {
	return fmt.Sprintf("id:%v, name:%s", u.ID, u.Name)
}

func MigrateAll(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Coin{})
	db.AutoMigrate(&Account{})
}
