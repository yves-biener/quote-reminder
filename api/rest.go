package quote

import (
	"fmt"
	"log"
	"net/http"
	db "quote/db"
	"strconv"

	"github.com/gorilla/mux"
)

var database *db.Database

func help(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	// TODO: create help message for all the available functions
	w.Write([]byte(`{"message": "not found"}`))
}

func fail(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
}

func getTopics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	topics, err := database.GetTopics()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v", topics)))
}

func getTopic(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	id := -1
	var err error
	if val, ok := pathParams["id"]; ok {
		id, err = strconv.Atoi(val)
		if err != nil {
			fail(w, err)
			return
		}
	}
	// id should now not be -1
	topic, err := database.GetTopic(id)
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v", topic)))
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

func getLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := database.GetLanguages()
	if err != nil {
		fail(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`%v`, languages)))
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

	// books := api.PathPrefix("/books").Subrouter()
	// quotes := api.PathPrefix("/quotes").Subrouter()
	return
}
