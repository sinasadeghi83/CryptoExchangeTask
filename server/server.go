package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sinasadeghi83/SwapTask/account"
	"github.com/sinasadeghi83/SwapTask/user"
)

var API_ROUTES = map[string]http.HandlerFunc{
	"GET /user/{id:[-]?[0-9]+}":     user.HandleGetUser,
	"GET /account/{id:[-]?[0-9]+}":  account.HandleGetAccount,
	"POST /convert":                 account.HandleConversion,
	"POST /convert/{id:[-]?[0-9]+}": account.HandleFinalConversion,
}

func addRoutes(router *mux.Router) {
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
		fmt.Printf("Request has been made through host %s for Hello world!\n", r.Host)
	})

	for path, handlerFunc := range API_ROUTES {
		splitedPath := strings.Split(path, " ")
		method, route := splitedPath[0], splitedPath[1]
		router.HandleFunc(route, handlerFunc).Methods(method)
	}
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
