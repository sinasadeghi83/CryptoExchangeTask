package account

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sinasadeghi83/SwapTask/db"
	"github.com/sinasadeghi83/SwapTask/model"
	"github.com/sinasadeghi83/SwapTask/utils"
)

type accountResponse struct {
	model.Account
	Uscent uint `json:"uscent"`
}

func HandleGetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	db, err := db.NewConn()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
