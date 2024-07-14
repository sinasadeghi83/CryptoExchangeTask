package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	id    int
	name  string
	coins []map[Coin]int
}

type Coin struct {
	id   int
	name string
}

func (u User) String() string {
	return fmt.Sprintf("id:%v, name:%s", u.id, u.name)
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	users := map[int]User{
		1: {id: 1, name: "sina"},
		2: {id: 2, name: "ali"},
		3: {id: 3, name: "hassan"},
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.Write([]byte("Invalid user id!"))
		return
	}
	w.Write([]byte(users[userID].String()))
}
