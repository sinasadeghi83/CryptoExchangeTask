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

	result := db.Where("user_id = ?", userID).Preload("Coin").Preload("User").Find(&accounts)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Error"))
		return
	}

	response := ""
	for _, account := range accounts {
		response += fmt.Sprintf("%s of %s balance is %d %s", account.User.Name, account.Coin.Name, account.Balance, account.Coin.Unit)
	}

	w.Write([]byte(response))
}
