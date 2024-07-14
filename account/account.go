package account

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sinasadeghi83/SwapTask/db"
	"github.com/sinasadeghi83/SwapTask/model"
	"github.com/sinasadeghi83/SwapTask/utils"
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

	utils.EncodeResponse(w, r, http.StatusOK, accounts)
}
