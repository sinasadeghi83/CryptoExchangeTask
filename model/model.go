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
	CoinID  int    `gorm:"uniqueIndex:idx_userid_coinid" json:"coin_id"`
	Coin    Coin   `json:"coin"`
	UserID  int    `gorm:"uniqueIndex:idx_userid_coinid" json:"user_id"`
	User    User   `json:"user"`
	Balance uint64 `json:"balance"`
}

type Coin struct {
	BaseModel
	Name       string `gorm:"index:idx_coin_name" json:"name"`
	Unit       string `json:"unit"`
	UnitFactor uint64 `json:"unit_factor"`
}

type Conversion struct {
	BaseModel
	UserID         int    `json:"user_id"`
	User           User   `gorm:"foreignKey:UserID" json:"user"`
	SourceCoinID   int    `json:"source_coin_id"`
	SourceCoin     Coin   `gorm:"foreignKey:SourceCoinID" json:"source_coin"`
	DestCoinID     int    `json:"dest_coin_id"`
	DestCoin       Coin   `gorm:"foreignKey:DestCoinID" json:"dest_coin"`
	SourceAmount   uint64 `json:"source_amount"`
	ExpectedAmount uint64 `json:"expected_amount"`
	State          uint
}

func (conversion *Conversion) CheckBalance(db *gorm.DB) error {
	var account Account

	result := db.Where("user_id=? and coin_id=?", conversion.UserID, conversion.SourceCoinID).First(&account)
	if result.RowsAffected == 0 {
		return fmt.Errorf("Account not found")
	}

	if account.Balance < conversion.SourceAmount {
		return fmt.Errorf("insufficient balance")
	}
	return nil
}

func (conversion *Conversion) Convert(db *gorm.DB) error {
	var sourceAccount, destAccount Account
	db.Where(Account{UserID: conversion.UserID, CoinID: conversion.SourceCoinID}).First(&sourceAccount)
	db.Where(Account{UserID: conversion.UserID, CoinID: conversion.DestCoinID}).FirstOrCreate(&destAccount)

	sourceAccount.Balance -= conversion.SourceAmount
	destAccount.Balance += conversion.ExpectedAmount

	if err := db.Save(&sourceAccount).Error; err != nil {
		return err
	}

	if err := db.Save(&destAccount).Error; err != nil {
		return err
	}

	db.Delete(conversion)

	return nil
}

func (conversion *Conversion) Validate() error {
	if time.Since(conversion.UpdatedAt) > time.Minute {
		return fmt.Errorf("conversion has expired")
	}
	return nil
}

func (conversion *Conversion) LoadAssociates(db *gorm.DB) {
	db.Find(&conversion.SourceCoin, conversion.SourceCoinID)
	db.Find(&conversion.DestCoin, conversion.DestCoinID)
}

func (conversion *Conversion) CalculateExpected() error {
	sourcePrice, err := conversion.SourceCoin.RetrievePrice("USD")
	if err != nil {
		return err
	}
	destPrice, err := conversion.DestCoin.RetrievePrice("USD")
	if err != nil {
		return err
	}

	convertFactor := destPrice / sourcePrice
	conversion.ExpectedAmount = uint64(float64(conversion.SourceAmount) * convertFactor * float64(conversion.DestCoin.UnitFactor/conversion.SourceCoin.UnitFactor))
	return nil
}

func (account Account) CalculateUSCent() (uint64, error) {
	price, err := account.Coin.RetrievePrice("USD")
	if err != nil {
		return 0, err
	}

	stdBalance := float64(account.Balance) / float64(account.Coin.UnitFactor)
	return uint64(price*stdBalance) * 100, nil
}

func (coin Coin) RetrievePrice(destCoinName string) (float64, error) {
	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/generateAvg?fsym=%s&tsym=%s&e=coinbase", coin.Name, destCoinName)
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

	return data.RAW.PRICE, nil
}

func (u User) String() string {
	return fmt.Sprintf("id:%v, name:%s", u.ID, u.Name)
}

func MigrateAll(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Coin{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Conversion{})
}
