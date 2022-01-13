package quote

import (
	"fmt"
	"log"
	"net/smtp"
	db "quote/db"
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
	err = smtp.SendMail(c.SmtpHost+":"+string(c.SmtpHost),
		auth, c.Sender, c.Receiver, []byte(message(quotes)))
	return
}

func selectQuotes(database *db.Database) (quotes []db.Quote) {

	return
}

func message(quotes []db.Quote) (message string) {
	for _, quote := range quotes {
		message += fmt.Sprintf(`"%s" from "%s" by %s\n`,
			quote.Quote, quote.Book.Title, quote.Book.Author.Name)
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
