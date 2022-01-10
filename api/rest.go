package quote

import (
	"fmt"
	"log"
	"net/http"
	db "quote/db"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var database *db.Database

func help(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Help")
	w.WriteHeader(http.StatusNotFound)
	// TODO: create help message for all the available functions
	w.Write([]byte(`{"message": "not found"}`))
}

func fail(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
}

func getTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := database.GetTopics()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, topics)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, topic)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, books)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, quotes)))
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
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func getAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := database.GetAuthors()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, authors)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, author)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, books)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, quotes)))
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
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func getLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := database.GetLanguages()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, languages)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, language)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, books)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, quotes)))
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
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := database.GetBooks()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, books)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, book)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, quotes)))
}

func postBook(w http.ResponseWriter, r *http.Request) {
	var err error
	author := database.NewAuthor()
	author.Name = r.PostFormValue("Name")
	topic := database.NewTopic()
	topic.Topic = r.PostFormValue("Topic")
	language := database.NewLanguage()
	language.Language = r.PostFormValue("Language")
	book := database.NewBook(author, topic, language)
	book.Title = r.PostFormValue("Title")
	book.ISBN.Scan(r.PostFormValue("ISBN"))
	book.ReleaseDate, err = time.Parse(time.ANSIC, r.PostFormValue("ReleaseDate"))
	if err != nil {
		fail(w, err)
		return
	}
	id, err := book.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func getQuotes(w http.ResponseWriter, r *http.Request) {
	quotes, err := database.GetQuotes()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, quotes)))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, quote)))
}

func postQuote(w http.ResponseWriter, r *http.Request) {
	var err error
	author := database.NewAuthor()
	author.Name = r.PostFormValue("Name")
	topic := database.NewTopic()
	topic.Topic = r.PostFormValue("Topic")
	language := database.NewLanguage()
	language.Language = r.PostFormValue("Language")
	book := database.NewBook(author, topic, language)
	book.Title = r.PostFormValue("Title")
	book.ISBN.Scan(r.PostFormValue("ISBN"))
	book.ReleaseDate, err = time.Parse(time.ANSIC, r.PostFormValue("ReleaseDate"))
	if err != nil {
		fail(w, err)
		return
	}
	quote := database.NewQuote(book)
	quote.Quote = r.PostFormValue("Quote")
	quote.RecordDate = time.Now()
	quote.Page, err = strconv.Atoi(r.PostFormValue("Page"))
	if err != nil {
		fail(w, err)
		return
	}
	id, err := book.Commit()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
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
	Patch  = "PATCH"  // -> database update
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
		HandlerFunc(getTopics).
		Methods(Get)
	topicsRouter.
		Path("/{id:[0-9]+}").
		HandlerFunc(getTopic).
		Methods(Get)
	topicsRouter.
		Path("/{id:[0-9]+}/books").
		HandlerFunc(getRelatedBooksOfTopic).
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

	authorsRouter := root.PathPrefix("/authors").Subrouter()
	// Get Methods
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
		HandlerFunc(getRelatedBooksOfAuthor).
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

	languagesRouter := root.PathPrefix("/languages").Subrouter()
	// Get Methods
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
		HandlerFunc(getRelatedBooksOfLanguage).
		Methods(Get)
	languagesRouter.
		Path("/{id:[0-9]+}/quotes").
		HandlerFunc(getRelatedQuotesOfLanguage).
		Methods(Get)
	// Post Methods
	languagesRouter.
		Path("").
		HandlerFunc(postLanguage).
		Methods(Post)

	booksRouter := root.PathPrefix("/books").Subrouter()
	// Get Methods
	booksRouter.
		Path("").
		HandlerFunc(getBooks).
		Methods(Get)
	booksRouter.
		Path("/{id:[0-9]+}").
		HandlerFunc(getBook).
		Methods(Get)
	booksRouter.
		Path("/{id:[0-9]+}").
		HandlerFunc(getRelatedQuotesOfBook).
		Methods(Get)
	// Post Methods
	booksRouter.
		Path("").
		HandlerFunc(postBook).
		Methods(Post)

	quotesRouter := root.PathPrefix("/quotes").Subrouter()
	// Get Methods
	quotesRouter.
		Path("").
		HandlerFunc(getQuotes).
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
	return
}
