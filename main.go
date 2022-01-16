package main

import (
	"encoding/json"
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
	// connect to local database
	database, err := db.Connect(dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// read mail service configuration
	configJson, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatal(err)
	}
	config := mail.Config{}
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		log.Fatal(err)
	}

	// start both services
	go mail.Service(database, config)
	api.RunServer(database)
}
