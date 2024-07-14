package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sinasadeghi83/SwapTask/user"
)

func addRoutes(router *mux.Router) {
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
		fmt.Printf("Request has been made through host %s for Hello world!\n", r.Host)
	})

	router.HandleFunc("/user/{id:[-]?[0-9]+}", user.HandleGetUser).Methods("GET")
}

func NewHandler() http.Handler {
	router := mux.NewRouter()
	addRoutes(router)

	return router
}

func NewServer() *http.Server {
	server := http.Server{Addr: ":8080", Handler: NewHandler()}
	fmt.Printf("Server created on %s\n", server.Addr)
	return &server
}
