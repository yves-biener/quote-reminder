package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	api "quote/api"
	db "quote/db"
	mail "quote/mail"
	"time"
)

const (
	dbFilename           = "./test_copy.sqlite"
	configFilename       = "./config.json"
	serverConfigFilename = "./server-config.json"
)

type ServerConfig struct {
	Address string
	Port    int
	Timeout time.Duration
}

func ApiService(database *db.Database) {
	// read api server configuration
	configJson, err := ioutil.ReadFile(serverConfigFilename)
	if err != nil {
		log.Fatal(err)
	}
	serverConfig := ServerConfig{}
	err = json.Unmarshal(configJson, &serverConfig)
	if err != nil {
		log.Fatal(err)
	}
	// start api service
	server := &http.Server{
		Handler: api.GetRouter(database),
		Addr: fmt.Sprintf("%s:%d",
			serverConfig.Address, serverConfig.Port),
		WriteTimeout: serverConfig.Timeout,
		ReadTimeout:  serverConfig.Timeout,
	}
	log.Fatal(server.ListenAndServe())
}

func MailService(database *db.Database) {
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
	// start mail service
	mail.Service(database, config)
}

func main() {
	// connect/create to local database
	database, err := db.Connect(dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// start services concurrently
	go MailService(database)
	go ApiService(database)

	fmt.Println("Services are running... Press enter to cancel...")
	fmt.Scanln()
}
