package account

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sinasadeghi83/SwapTask/db"
	"github.com/sinasadeghi83/SwapTask/model"
)

func HandleGetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	db, err := db.NewConn()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var accounts []model.Account

	result := db.Model(model.Account{UserID: userID}).Joins("Coin").Find(&accounts)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error"))
		return
	}

	response := ""
	for _, account := range accounts {
		response += fmt.Sprintf("%s balance is %d %s", account.Coin.Name, account.Balance, account.Coin.Unit)
	}

	w.Write([]byte(response))
}
