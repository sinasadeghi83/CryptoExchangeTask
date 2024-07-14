package server

import (
	"fmt"
	"net/http"
)

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
		fmt.Printf("Request has been made through host %s for Hello world!\n", r.Host)
	})
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()

	addRoutes(mux)

	return mux
}

func NewServer() *http.Server {
	server := http.Server{Addr: ":8080", Handler: NewHandler()}
	fmt.Printf("Server created on %s\n", server.Addr)
	return &server
}
