package main

import (
	"log"
	api "quote/api"
	db "quote/db"
)

const filename = "./test.sqlite"

func main() {
	database, err := db.Connect(filename)
	defer database.Close()
	if err != nil {
		log.Fatal(err)
	}
	api.RunServer(database)
	// TODO move these into an test for the db implementation
	// database, err := db.Connect("./test.sqlite")
	// defer database.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// topic := database.NewTopic()
	// topic.Topic = "Schlechte Witze"

	// author := database.NewAuthor()
	// author.Name = "Hans Peter"

	// language := database.NewLanguage()
	// language.Language = "Deutsch"

	// book := database.NewBook(author, topic, language)
	// book.Title = "Deine Mudda"
	// book.ReleaseDate = time.Now()
	// book.ISBN = "ISBN-1337-69-420"

	// quote := database.NewQuote(book)
	// quote.Quote = "Du kannst mich mal am Arsch lecken!"
	// quote.Page = 69

	// _, err = quote.Commit()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// books, err := database.GetBooks()
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println(books)
	// }

	// quotes, err := database.GetQuotes()
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println(quotes)
	// }
}
