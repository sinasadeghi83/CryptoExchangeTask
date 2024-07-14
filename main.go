package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
		fmt.Printf("Request has been made through host %s for Hello world!\n", r.Host)
	})

	server := http.Server{Addr: ":8080", Handler: mux}

	fmt.Printf("Server is running on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
