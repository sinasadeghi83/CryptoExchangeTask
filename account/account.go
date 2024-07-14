package account

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sinasadeghi83/SwapTask/db"
	"github.com/sinasadeghi83/SwapTask/model"
	"github.com/sinasadeghi83/SwapTask/utils"
)

type accountResponse struct {
	model.Account
	Uscent uint64 `json:"uscent"`
}

type ConversionForm struct {
	UserID       int    `json:"user_id"`
	SourceCoinID int    `json:"source_coin_id"`
	DestCoinID   int    `json:"dest_coin_id"`
	Amount       uint64 `json:"amount"`
}

func HandleConversion(w http.ResponseWriter, r *http.Request) {
	db, err := db.NewConn()
	if err != nil {
		utils.EncodeResponse(w, r, http.StatusInternalServerError, "Internal Error")
		return
	}

	reqData, err := utils.DecodeJson[ConversionForm](r)

	if err != nil {
		utils.EncodeResponse(w, r, http.StatusBadRequest, "Conversion form is not valid!")
		return
	}

	if db.First(&model.User{}, reqData.UserID).RowsAffected == 0 {
		utils.EncodeResponse(w, r, http.StatusNotFound, "User not found!")
		return
	}

	conversion := model.Conversion{
		SourceCoinID: reqData.SourceCoinID,
		DestCoinID:   reqData.DestCoinID,
		SourceAmount: reqData.Amount,
		UserID:       reqData.UserID,
		State:        0,
	}

	if err := conversion.CheckBalance(db); err != nil {
		utils.EncodeResponse(w, r, http.StatusOK, "Insufficient Balance!")
		return
	}

	conversion.LoadAssociates(db)
	if err := conversion.CalculateExpected(); err != nil {
		utils.EncodeResponse(w, r, http.StatusBadRequest, "Conversion is not valid!")
		fmt.Println(err)
		return
	}

	result := db.Create(&conversion)

	if result.Error != nil {
		utils.EncodeResponse(w, r, http.StatusBadRequest, "Conversion is not valid!")
		fmt.Println(result.Error)
		return
	}
	db.Save(&conversion)

	utils.EncodeResponse(w, r, http.StatusOK, conversion)
}

func HandleGetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	db, err := db.NewConn()
	if err != nil {
		utils.EncodeResponse(w, r, http.StatusInternalServerError, "Internal Error")
		return
	}

	var accounts []model.Account

	result := db.Where("user_id = ?", userID).Preload("Coin").Preload("User").Find(&accounts)

	if result.Error != nil {
		utils.EncodeResponse(w, r, http.StatusInternalServerError, "Internal Error")
		return
	}
	data, err := retrieveAccountWithCent(accounts)

	if err != nil {
		utils.EncodeResponse(w, r, http.StatusInternalServerError, "Could not retrieve the prices")
		return
	}
	utils.EncodeResponse(w, r, http.StatusOK, data)
}

func retrieveAccountWithCent(accounts []model.Account) ([]accountResponse, error) {
	data := make([]accountResponse, 0)
	for _, account := range accounts {
		uscent, err := account.CalculateUSCent()
		if err != nil {
			return nil, err
		}
		data = append(data, accountResponse{
			Account: account,
			Uscent:  uscent,
		})
	}
	return data, nil
}
