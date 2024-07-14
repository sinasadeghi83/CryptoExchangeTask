package main

import (
	"log"

	"github.com/sinasadeghi83/SwapTask/db"
	"github.com/sinasadeghi83/SwapTask/server"
)

func main() {
	if _, err := db.SetupDB(); err != nil {
		log.Fatal("Unable to setup db: ", err)
		return
	}

	if err := db.CreateDummyData(); err != nil {
		log.Fatal("Unable to create dummy data: ", err)
		return
	}

	if err := server.NewServer().ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
