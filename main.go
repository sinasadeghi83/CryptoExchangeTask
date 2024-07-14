package main

import (
	"log"
	"server/server"
)

func main() {
	if err := server.NewServer().ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
