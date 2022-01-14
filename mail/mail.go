package quote

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	db "quote/db"
	"strings"
	"time"
)

type Config struct {
	Sender   string   `json: "sender"`
	Password string   `json: "password"`
	Receiver []string `json: "receiver"`
	SmtpHost string   `json: "smtpHost"`
	SmtpPort int      `json: "smtpPort"`
}

func (c Config) sendMail(quotes []db.Quote) (err error) {
	auth := smtp.PlainAuth("", c.Sender, c.Password, c.SmtpHost)
	err = smtp.SendMail(fmt.Sprintf("%s:%d", c.SmtpHost, c.SmtpPort),
		auth, c.Sender, c.Receiver, []byte(c.message(quotes)))
	return
}

func (c Config) message(quotes []db.Quote) (message string) {
	// Set header
	message += fmt.Sprintf("From: %s\r\n", c.Sender)
	message += fmt.Sprintf("To: %s\r\n", strings.Join(c.Receiver, " "))
	message += "Subject: Quote-reminder\r\n\r\n"
	// Set body
	for _, quote := range quotes {
		message += fmt.Sprintf("'%s' from '%s' by %s\n",
			quote.Quote, quote.Book.Title, quote.Book.Author.Name)
	}
	return
}

func selectQuotes(database *db.Database) (selection []db.Quote) {
	quotes, err := database.GetQuotes()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 5; i += 1 {
		selection = append(selection, quotes[rand.Intn(len(quotes))])
	}
	return
}

func Service(database *db.Database, config Config) {
	for _ = range time.Tick(time.Hour * 24) {
		err := config.sendMail(selectQuotes(database))
		if err != nil {
			log.Fatal(err)
		}
	}
}
