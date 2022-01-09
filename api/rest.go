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

func getAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := database.GetAuthors()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, authors)))
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
	log.Fatal(http.ListenAndServe(":8000", GetRouter(db)))
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
	root.HandleFunc("", help)

	topicsRouter := root.PathPrefix("/topics").Subrouter()
	// Get Methods
	topicsRouter.HandleFunc("", getTopics).Methods(Get)
	topicsRouter.HandleFunc("/{id}", getTopic).Methods(Get)
	// Post Methods
	topicsRouter.HandleFunc("", postTopic).Methods(Post)

	authorsRouter := root.PathPrefix("/authors").Subrouter()
	// Get Methods
	authorsRouter.HandleFunc("", getAuthors).Methods(Get)
	authorsRouter.HandleFunc("/{id}", getAuthor).Methods(Get)
	// Post Methods
	authorsRouter.HandleFunc("", postAuthor).Methods(Post)

	languagesRouter := root.PathPrefix("/languages").Subrouter()
	// Get Methods
	languagesRouter.HandleFunc("", getLanguages).Methods(Get)
	languagesRouter.HandleFunc("/{id}", getLanguage).Methods(Get)
	// Post Methods
	languagesRouter.HandleFunc("", postLanguage).Methods(Post)

	booksRouter := root.PathPrefix("/books").Subrouter()
	// Get Methods
	booksRouter.HandleFunc("", getBooks).Methods(Get)
	booksRouter.HandleFunc("/{id}", getBook).Methods(Get)
	// Post Methods
	booksRouter.HandleFunc("", postBook).Methods(Post)

	quotesRouter := root.PathPrefix("/quotes").Subrouter()
	// Get Methods
	quotesRouter.HandleFunc("", getQuotes).Methods(Get)
	quotesRouter.HandleFunc("/{id}", getQuote).Methods(Get)
	// Post Methods
	quotesRouter.HandleFunc("", postQuote).Methods(Post)
	return
}
