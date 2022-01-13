package main

import (
	"log"
	api "quote/api"
	db "quote/db"
)

const filename = "./test.sqlite"

func main() {
	database, err := db.Connect(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	api.RunServer(database)
}
