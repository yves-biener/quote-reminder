package quote

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	db "quote/db"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var database *db.Database

func help(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// TODO: create help message for all the available functions maybe this
	// can be automated?
	w.Write([]byte(`{"message": "not found"}`))
}

func fail(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
}

// TODO: move filtering to db.go maybe create custom type for slice of elements
// as this contains information about what is searched for in the search
// functions
func filterBooks(books []db.Book, filters []string) (filtered []db.Book) {
	for _, filter := range filters {
		for _, book := range books {
			if strings.Contains(book.Title, filter) ||
				strings.Contains(book.ISBN.String, filter) {
				filtered = append(filtered, book)
			}
		}
	}
	return
}

func filterQuotes(quotes []db.Quote, filters []string) (filtered []db.Quote) {
	for _, filter := range filters {
		for _, quote := range quotes {
			if strings.Contains(quote.Quote, filter) {
				filtered = append(filtered, quote)
			}
		}
	}
	return
}

func getTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := database.GetTopics()
	if err != nil {
		fail(w, err)
		return
	}
	var jsonTopics []string
	for _, topic := range topics {
		jsonTopic, err := json.Marshal(topic)
		if err != nil {
			fail(w, err)
			return
		}
		jsonTopics = append(jsonTopics, string(jsonTopic))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonTopics, ","))))
}

func searchTopics(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	var topics []db.Topic
	if val, ok := pathParams["search"]; ok {
		search := strings.Split(val, " ")
		for _, q := range search {
			searchResult, err := database.SearchTopics(q)
			if err != nil {
				fail(w, err)
				return
			}
			for _, topic := range searchResult {
				// there can be search Results more than one in
				// the resulting topics slice
				topics = append(topics, topic)
			}
		}
	}
	var jsonTopics []string
	for _, topic := range topics {
		jsonTopic, err := json.Marshal(topic)
		if err != nil {
			fail(w, err)
			return
		}
		jsonTopics = append(jsonTopics, string(jsonTopic))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonTopics, ","))))
}

func getTopic(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	topic, err := database.GetTopic(id)
	if err != nil {
		fail(w, err)
		return
	}
	defaultTopic := db.Topic{}
	if topic == defaultTopic {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonTopic, err := json.Marshal(topic)
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%s`, string(jsonTopic))))
}

func getRelatedBooksOfTopic(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	books, err := database.RelatedBooksOfTopic(id)
	if err != nil {
		fail(w, err)
		return
	}
	if val, ok := pathParams["filter"]; ok {
		filters := strings.Split(val, " ")
		books = filterBooks(books, filters)
	}
	var jsonBooks []string
	for _, book := range books {
		jsonBook, err := json.Marshal(book)
		if err != nil {
			fail(w, err)
			return
		}
		jsonBooks = append(jsonBooks, string(jsonBook))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonBooks, ","))))
}

func getRelatedQuotesOfTopic(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	quotes, err := database.RelatedQuotesOfTopic(id)
	if err != nil {
		fail(w, err)
		return
	}
	if val, ok := pathParams["filter"]; ok {
		filters := strings.Split(val, " ")
		quotes = filterQuotes(quotes, filters)
	}
	var jsonQuotes []string
	for _, quote := range quotes {
		jsonQuote, err := json.Marshal(quote)
		if err != nil {
			fail(w, err)
			return
		}
		jsonQuotes = append(jsonQuotes, string(jsonQuote))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonQuotes, ","))))
}

func postTopic(w http.ResponseWriter, r *http.Request) {
	topic := database.NewTopic()
	topic.Topic = r.PostFormValue("Topic")
	id, err := topic.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, id)))
}

func patchTopic(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("Id")
	topicId, err := strconv.Atoi(id)
	if err != nil {
		fail(w, err)
		return
	}
	topic, err := database.GetTopic(topicId)
	if err != nil {
		fail(w, err)
		return
	}
	defaultTopic := db.Topic{}
	if defaultTopic == topic {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	topic.Topic = r.PostFormValue("Topic")
	_, err = topic.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, topicId)))
}

func getAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := database.GetAuthors()
	if err != nil {
		fail(w, err)
		return
	}
	var jsonAuthors []string
	for _, author := range authors {
		jsonAuthor, err := json.Marshal(author)
		if err != nil {
			fail(w, err)
			return
		}
		jsonAuthors = append(jsonAuthors, string(jsonAuthor))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonAuthors, ","))))
}

func searchAuthors(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	var authors []db.Author
	if val, ok := pathParams["search"]; ok {
		search := strings.Split(val, " ")
		for _, q := range search {
			searchResult, err := database.SearchAuthors(q)
			if err != nil {
				fail(w, err)
				return
			}
			for _, author := range searchResult {
				// there can be search Results more than one in
				// the resulting topics slice
				authors = append(authors, author)
			}
		}
	}
	var jsonAuthors []string
	for _, author := range authors {
		jsonAuthor, err := json.Marshal(author)
		if err != nil {
			fail(w, err)
			return
		}
		jsonAuthors = append(jsonAuthors, string(jsonAuthor))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonAuthors, ","))))
}

func getAuthor(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	author, err := database.GetAuthor(id)
	if err != nil {
		fail(w, err)
		return
	}
	defaultAuthor := db.Author{}
	if author == defaultAuthor {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonAuthor, err := json.Marshal(author)
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%s`, string(jsonAuthor))))
}

func getRelatedBooksOfAuthor(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	books, err := database.RelatedBooksOfAuthor(id)
	if err != nil {
		fail(w, err)
		return
	}
	if val, ok := pathParams["filter"]; ok {
		filters := strings.Split(val, " ")
		books = filterBooks(books, filters)
	}
	var jsonBooks []string
	for _, book := range books {
		jsonBook, err := json.Marshal(book)
		if err != nil {
			fail(w, err)
			return
		}
		jsonBooks = append(jsonBooks, string(jsonBook))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonBooks, ","))))
}

func getRelatedQuotesOfAuthor(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	quotes, err := database.RelatedQuotesOfAuthor(id)
	if err != nil {
		fail(w, err)
		return
	}
	if val, ok := pathParams["filter"]; ok {
		filters := strings.Split(val, " ")
		quotes = filterQuotes(quotes, filters)
	}
	var jsonQuotes []string
	for _, quote := range quotes {
		jsonQuote, err := json.Marshal(quote)
		if err != nil {
			fail(w, err)
			return
		}
		jsonQuotes = append(jsonQuotes, string(jsonQuote))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonQuotes, ","))))
}

func postAuthor(w http.ResponseWriter, r *http.Request) {
	author := database.NewAuthor()
	author.Name = r.PostFormValue("Name")
	id, err := author.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, id)))
}

func patchAuthor(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("Id")
	authorId, err := strconv.Atoi(id)
	if err != nil {
		fail(w, err)
		return
	}
	author, err := database.GetAuthor(authorId)
	if err != nil {
		fail(w, err)
		return
	}
	defaultAuthor := db.Author{}
	if defaultAuthor == author {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	author.Name = r.PostFormValue("Name")
	_, err = author.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, authorId)))
}

func getLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := database.GetLanguages()
	if err != nil {
		fail(w, err)
		return
	}
	var jsonLanguages []string
	for _, language := range languages {
		jsonLanguage, err := json.Marshal(language)
		if err != nil {
			fail(w, err)
			return
		}
		jsonLanguages = append(jsonLanguages, string(jsonLanguage))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonLanguages, ","))))
}

func searchLanguages(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	var languages []db.Language
	if val, ok := pathParams["search"]; ok {
		search := strings.Split(val, " ")
		for _, q := range search {
			searchResult, err := database.SearchLanguages(q)
			if err != nil {
				fail(w, err)
				return
			}
			for _, language := range searchResult {
				// there can be search Results more than one in
				// the resulting topics slice
				languages = append(languages, language)
			}
		}
	}
	var jsonLanguages []string
	for _, language := range languages {
		jsonLanguage, err := json.Marshal(language)
		if err != nil {
			fail(w, err)
			return
		}
		jsonLanguages = append(jsonLanguages, string(jsonLanguage))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonLanguages, ","))))
}

func getLanguage(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	language, err := database.GetLanguage(id)
	if err != nil {
		fail(w, err)
		return
	}
	defaultLanguage := db.Language{}
	if language == defaultLanguage {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonLanguage, err := json.Marshal(language)
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%s`, string(jsonLanguage))))
}

func getRelatedBooksOfLanguage(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	books, err := database.RelatedBooksOfLanguage(id)
	if err != nil {
		fail(w, err)
		return
	}
	if val, ok := pathParams["filter"]; ok {
		filters := strings.Split(val, " ")
		books = filterBooks(books, filters)
	}
	var jsonBooks []string
	for _, book := range books {
		jsonBook, err := json.Marshal(book)
		if err != nil {
			fail(w, err)
			return
		}
		jsonBooks = append(jsonBooks, string(jsonBook))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonBooks, ","))))
}

func getRelatedQuotesOfLanguage(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	quotes, err := database.RelatedQuotesOfLanguage(id)
	if err != nil {
		fail(w, err)
		return
	}
	if val, ok := pathParams["filter"]; ok {
		filters := strings.Split(val, " ")
		quotes = filterQuotes(quotes, filters)
	}
	var jsonQuotes []string
	for _, quote := range quotes {
		jsonQuote, err := json.Marshal(quote)
		if err != nil {
			fail(w, err)
			return
		}
		jsonQuotes = append(jsonQuotes, string(jsonQuote))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonQuotes, ","))))
}

func postLanguage(w http.ResponseWriter, r *http.Request) {
	language := database.NewLanguage()
	language.Language = r.PostFormValue("Language")
	id, err := language.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, id)))
}

func patchLanguage(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("Id")
	languageId, err := strconv.Atoi(id)
	if err != nil {
		fail(w, err)
		return
	}
	language, err := database.GetLanguage(languageId)
	if err != nil {
		fail(w, err)
		return
	}
	defaultLanguage := db.Language{}
	if defaultLanguage == language {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	language.Language = r.PostFormValue("Language")
	_, err = language.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, languageId)))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := database.GetBooks()
	if err != nil {
		fail(w, err)
		return
	}
	var jsonBooks []string
	for _, book := range books {
		jsonBook, err := json.Marshal(book)
		if err != nil {
			fail(w, err)
			return
		}
		jsonBooks = append(jsonBooks, string(jsonBook))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonBooks, ","))))
}

func searchBooks(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	var books []db.Book
	if val, ok := pathParams["search"]; ok {
		search := strings.Split(val, " ")
		for _, q := range search {
			searchResult, err := database.SearchBooks(q)
			if err != nil {
				fail(w, err)
				return
			}
			for _, book := range searchResult {
				books = append(books, book)
			}
		}
	}
	var jsonBooks []string
	for _, book := range books {
		jsonBook, err := json.Marshal(book)
		if err != nil {
			fail(w, err)
			return
		}
		jsonBooks = append(jsonBooks, string(jsonBook))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonBooks, ","))))
}

func getBook(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	book, err := database.GetBook(id)
	if err != nil {
		fail(w, err)
		return
	}
	defaultBook := db.Book{}
	if book == defaultBook {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonBook, err := json.Marshal(book)
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%s`, string(jsonBook))))
}

func getRelatedQuotesOfBook(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	quotes, err := database.RelatedQuotesOfBook(id)
	if err != nil {
		fail(w, err)
		return
	}
	if val, ok := pathParams["filter"]; ok {
		filters := strings.Split(val, " ")
		quotes = filterQuotes(quotes, filters)
	}
	var jsonQuotes []string
	for _, quote := range quotes {
		jsonQuote, err := json.Marshal(quote)
		if err != nil {
			fail(w, err)
			return
		}
		jsonQuotes = append(jsonQuotes, string(jsonQuote))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonQuotes, ","))))
}

func postBook(w http.ResponseWriter, r *http.Request) {
	var err error
	id := r.PostFormValue("AuthorId")
	authorId, err := strconv.Atoi(id)
	if err != nil {
		fail(w, err)
		return
	}
	author, err := database.GetAuthor(authorId)
	if err != nil {
		fail(w, err)
		return
	}
	defaultAuthor := db.Author{}
	if defaultAuthor == author {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id = r.PostFormValue("TopicId")
	topicId, err := strconv.Atoi(id)
	if err != nil {
		fail(w, err)
		return
	}
	topic, err := database.GetTopic(topicId)
	if err != nil {
		fail(w, err)
		return
	}
	defaultTopic := db.Topic{}
	if defaultTopic == topic {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id = r.PostFormValue("LanguageId")
	languageId, err := strconv.Atoi(id)
	if err != nil {
		fail(w, err)
		return
	}
	language, err := database.GetLanguage(languageId)
	if err != nil {
		fail(w, err)
		return
	}
	defaultLanguage := db.Language{}
	if defaultLanguage == language {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	book := database.NewBook(author, topic, language)
	book.Title = r.PostFormValue("Title")
	book.ISBN.Scan(r.PostFormValue("ISBN"))
	releaseDate := r.PostFormValue("ReleaseDate")
	if releaseDate != "" {
		book.ReleaseDate, err = time.Parse(time.ANSIC, releaseDate)
		if err != nil {
			fail(w, err)
			return
		}
	} else {
		book.ReleaseDate = time.Now()
	}
	bookId, err := book.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, bookId)))
}

func patchBook(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("Id")
	bookId, err := strconv.Atoi(id)
	if err != nil {
		fail(w, err)
		return
	}
	book, err := database.GetBook(bookId)
	if err != nil {
		fail(w, err)
		return
	}
	defaultBook := db.Book{}
	if defaultBook == book {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if title := r.PostFormValue("Title"); title != "" {
		book.Title = title
	}
	if isbn := r.PostFormValue("ISBN"); isbn != "" {
		book.ISBN.String = isbn
	}
	if releaseDate := r.PostFormValue("ReleaseDate"); releaseDate != "" {
		book.ReleaseDate, err = time.Parse(time.ANSIC, releaseDate)
		if err != nil {
			fail(w, err)
			return
		}
	}
	_, err = book.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, bookId)))
}

func getQuotes(w http.ResponseWriter, r *http.Request) {
	quotes, err := database.GetQuotes()
	if err != nil {
		fail(w, err)
		return
	}
	var jsonQuotes []string
	for _, quote := range quotes {
		jsonQuote, err := json.Marshal(quote)
		if err != nil {
			fail(w, err)
			return
		}
		jsonQuotes = append(jsonQuotes, string(jsonQuote))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonQuotes, ","))))
}

func searchQuotes(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	var quotes []db.Quote
	if val, ok := pathParams["search"]; ok {
		search := strings.Split(val, " ")
		for _, q := range search {
			searchResult, err := database.SearchQuotes(q)
			if err != nil {
				fail(w, err)
				return
			}
			for _, quote := range searchResult {
				// there can be search Results more than one in
				// the resulting topics slice
				quotes = append(quotes, quote)
			}
		}
	}
	var jsonQuotes []string
	for _, quote := range quotes {
		jsonQuote, err := json.Marshal(quote)
		if err != nil {
			fail(w, err)
			return
		}
		jsonQuotes = append(jsonQuotes, string(jsonQuote))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`[%s]`, strings.Join(jsonQuotes, ","))))
}

func getQuote(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	quote, err := database.GetQuote(id)
	if err != nil {
		fail(w, err)
		return
	}
	defaultQuote := db.Quote{}
	if quote == defaultQuote {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonQuote, err := json.Marshal(quote)
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%s`, string(jsonQuote))))
}

func postQuote(w http.ResponseWriter, r *http.Request) {
	var err error
	id := r.PostFormValue("BookId")
	bookId, err := strconv.Atoi(id)
	if err != nil {
		fail(w, err)
		return
	}
	book, err := database.GetBook(bookId)
	if err != nil {
		fail(w, err)
		return
	}
	defaultBook := db.Book{}
	if book == defaultBook {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	quote := database.NewQuote(book)
	quote.Quote = r.PostFormValue("Quote")
	quote.RecordDate = time.Now()
	page := r.PostFormValue("Page")
	if page != "" {
		quote.Page, err = strconv.Atoi(page)
		if err != nil {
			fail(w, err)
			return
		}
	} else {
		quote.Page = 0
	}
	quoteId, err := quote.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, quoteId)))
}

func patchQuote(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("Id")
	quoteId, err := strconv.Atoi(id)
	if err != nil {
		fail(w, err)
		return
	}
	quote, err := database.GetQuote(quoteId)
	if err != nil {
		fail(w, err)
		return
	}
	defaultQuote := db.Quote{}
	if defaultQuote == quote {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if text := r.PostFormValue("Quote"); text != "" {
		quote.Quote = text
	}
	if page := r.PostFormValue("Page"); page != "" {
		quote.Page, err = strconv.Atoi(page)
		if err != nil {
			fail(w, err)
			return
		}
	}
	_, err = quote.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"Id": %d}`, quoteId)))
}

func jsonContentWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func RunServer(db *db.Database) {
	server := &http.Server{
		Handler:      GetRouter(db),
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

// HTTP Methods
const (
	Get    = "GET"    // -> database select
	Post   = "POST"   // -> database insert
	Patch  = "PATCH"  // -> database patch
	Delete = "DELETE" // -> database drop
)

func GetRouter(db *db.Database) (router *mux.Router) {
	database = db
	router = mux.NewRouter()

	root := router.PathPrefix("/api").Subrouter()
	root.Use(jsonContentWrapper)
	root.Path("").HandlerFunc(help)

	topicsRouter := root.PathPrefix("/topics").Subrouter()
	// Get Methods
	topicsRouter.
		Path("").
		Queries("q", "{search}").
		HandlerFunc(searchTopics).
		Methods(Get)
	topicsRouter.
		Path("").
		HandlerFunc(getTopics).
		Methods(Get)
	topicsRouter.
		Path("/{id:[0-9]+}").
		HandlerFunc(getTopic).
		Methods(Get)
	topicsRouter.
		Path("/{id:[0-9]+}/books").
		Queries("q", "{filter}").
		HandlerFunc(getRelatedBooksOfTopic).
		Methods(Get)
	topicsRouter.
		Path("/{id:[0-9]+}/books").
		HandlerFunc(getRelatedBooksOfTopic).
		Methods(Get)
	topicsRouter.
		Path("/{id:[0-9]+}/quotes").
		Queries("q", "{filter}").
		HandlerFunc(getRelatedQuotesOfTopic).
		Methods(Get)
	topicsRouter.
		Path("/{id:[0-9]+}/quotes").
		HandlerFunc(getRelatedQuotesOfTopic).
		Methods(Get)
	// Post Methods
	topicsRouter.
		Path("").
		HandlerFunc(postTopic).
		Methods(Post)
	// Patch Methods
	topicsRouter.
		Path("").
		HandlerFunc(patchTopic).
		Methods(Patch)

	authorsRouter := root.PathPrefix("/authors").Subrouter()
	// Get Methods
	authorsRouter.
		Path("").
		Queries("q", "{search}").
		HandlerFunc(searchAuthors).
		Methods(Get)
	authorsRouter.
		Path("").
		HandlerFunc(getAuthors).
		Methods(Get)
	authorsRouter.
		Path("/{id:[0-9]+}").
		HandlerFunc(getAuthor).
		Methods(Get)
	authorsRouter.
		Path("/{id:[0-9]+}/books").
		Queries("q", "{filter}").
		HandlerFunc(getRelatedBooksOfAuthor).
		Methods(Get)
	authorsRouter.
		Path("/{id:[0-9]+}/books").
		HandlerFunc(getRelatedBooksOfAuthor).
		Methods(Get)
	authorsRouter.
		Path("/{id:[0-9]+}/quotes").
		Queries("q", "{filter}").
		HandlerFunc(getRelatedQuotesOfAuthor).
		Methods(Get)
	authorsRouter.
		Path("/{id:[0-9]+}/quotes").
		HandlerFunc(getRelatedQuotesOfAuthor).
		Methods(Get)
	// Post Methods
	authorsRouter.
		Path("").
		HandlerFunc(postAuthor).
		Methods(Post)
	// Patch Methods
	authorsRouter.
		Path("").
		HandlerFunc(patchAuthor).
		Methods(Patch)

	languagesRouter := root.PathPrefix("/languages").Subrouter()
	// Get Methods
	languagesRouter.
		Path("").
		Queries("q", "{search}").
		HandlerFunc(searchLanguages).
		Methods(Get)
	languagesRouter.
		Path("").
		HandlerFunc(getLanguages).
		Methods(Get)
	languagesRouter.
		Path("/{id:[0-9]+}").
		HandlerFunc(getLanguage).
		Methods(Get)
	languagesRouter.
		Path("/{id:[0-9]+}/books").
		Queries("q", "{filter}").
		HandlerFunc(getRelatedBooksOfLanguage).
		Methods(Get)
	languagesRouter.
		Path("/{id:[0-9]+}/quotes").
		Queries("q", "{filter}").
		HandlerFunc(getRelatedQuotesOfLanguage).
		Methods(Get)
	// Post Methods
	languagesRouter.
		Path("").
		HandlerFunc(postLanguage).
		Methods(Post)
	// Patch Methods
	languagesRouter.
		Path("").
		HandlerFunc(patchLanguage).
		Methods(Patch)

	booksRouter := root.PathPrefix("/books").Subrouter()
	// Get Methods
	booksRouter.
		Path("").
		Queries("q", "{search}").
		HandlerFunc(searchBooks).
		Methods(Get)
	booksRouter.
		Path("").
		HandlerFunc(getBooks).
		Methods(Get)
	booksRouter.
		Path("/{id:[0-9]+}").
		HandlerFunc(getBook).
		Methods(Get)
	booksRouter.
		Path("/{id:[0-9]+}/quotes").
		Queries("q", "{filter}").
		HandlerFunc(getRelatedQuotesOfBook).
		Methods(Get)
	booksRouter.
		Path("/{id:[0-9]+}/quotes").
		HandlerFunc(getRelatedQuotesOfBook).
		Methods(Get)
	// Post Methods
	booksRouter.
		Path("").
		HandlerFunc(postBook).
		Methods(Post)
	// Patch Methods
	booksRouter.
		Path("").
		HandlerFunc(patchBook).
		Methods(Patch)

	quotesRouter := root.PathPrefix("/quotes").Subrouter()
	// Get Methods
	quotesRouter.
		Path("").
		HandlerFunc(getQuotes).
		Methods(Get)
	quotesRouter.
		Path("").
		Queries("q", "{search}").
		HandlerFunc(searchQuotes).
		Methods(Get)
	quotesRouter.
		Path("/{id:[0-9]+}").
		HandlerFunc(getQuote).
		Methods(Get)
	// Post Methods
	quotesRouter.
		Path("").
		HandlerFunc(postQuote).
		Methods(Post)
	// Patch Methods
	quotesRouter.
		Path("").
		HandlerFunc(patchQuote).
		Methods(Patch)
	return
}
