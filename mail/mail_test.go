package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	db "quote/db"
	"strings"
	"testing"
)

const (
	testSource   = "./../test.sqlite"
	testDatabase = "./cur_test.sqlite"
	testConfig   = "./../test-config.json"
	lenError     = "The number of elements does not match\nexpected: %d\nactual: %d\n"
	headerError  = "Header of message does not match\nexpected: %s\nactual: %s\n"
)

func initDatabase(t *testing.T) {
	source, err := os.Open(testSource)
	if err != nil {
		t.Fatal(err)
	}
	defer source.Close()

	destination, err := os.Create(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer destination.Close()
	io.Copy(destination, source)
}

func initConfig() (config Config, err error) {
	var configJson []byte
	configJson, err = ioutil.ReadFile(testConfig)
	if err != nil {
		return
	}
	err = json.Unmarshal(configJson, &config)
	return
}

func TestSelectQuotes(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	selectedQuotes := selectQuotes(database)
	// Assert
	expectedLen := 5
	if actualLen := len(selectedQuotes); actualLen != expectedLen {
		t.Errorf(lenError, expectedLen, actualLen)
	}
}

func TestConfigMessage(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	config, err := initConfig()
	if err != nil {
		t.Fatal(err)
	}
	quotes, err := database.GetQuotes()
	if err != nil {
		t.Fatal(err)
	}
	// Act
	actualMessage := config.message(quotes)
	// Assert
	expectedMessage := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Quote-reminder\r\n\r\n",
		config.Sender, strings.Join(config.Receiver, " "))
	if !strings.HasPrefix(actualMessage, expectedMessage) {
		t.Errorf(headerError, expectedMessage, actualMessage)
	}
}

func TestSendMail(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	config, err := initConfig()
	if err != nil {
		t.Fatal(err)
	}
	quotes, err := database.GetQuotes()
	if err != nil {
		t.Fatal(err)
	}
	// Act
	err = config.sendMail(quotes)
	// Assert
	if err == nil {
		t.Error("Expected an error but got nil")
	}
}
