package main

import (
	"log"

	"github.com/sinasadeghi83/SwapTask/server"
)

func main() {
	if err := server.NewServer().ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
