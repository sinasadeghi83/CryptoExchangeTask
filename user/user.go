package user

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sinasadeghi83/SwapTask/db"
	"github.com/sinasadeghi83/SwapTask/model"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	db, err := db.NewConn()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := model.User{}
	result := db.First(&user, userID)

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found!"))
		return
	}

	w.Write([]byte(user.String()))
}
