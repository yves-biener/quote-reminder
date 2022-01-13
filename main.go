package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	api "quote/api"
	db "quote/db"
	mail "quote/mail"
)

const (
	dbFilename     = "./test.sqlite"
	configFilename = "./config.json"
)

func main() {
	database, err := db.Connect(dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	configJson, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatal(err)
	}
	config := mail.Config{}
	err = json.Unmarshal([]byte(configJson), &config)
	fmt.Println(config)
	go mail.Service(database, config)

	api.RunServer(database)
}
