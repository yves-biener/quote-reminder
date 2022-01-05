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
	// TODO create help message for all the available functions
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
	w.WriteHeader(http.StatusOK)
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
	root.HandleFunc("", help)

	topicsRouter := root.PathPrefix("/topics").Subrouter()
	// Get Methods
	topicsRouter.HandleFunc("", getTopics).Methods(Get)
	topicsRouter.HandleFunc("/{id}", getTopic).Methods(Get)
	// Post Methods
	topicsRouter.HandleFunc("", postTopic).Methods(Post)

	// authors := api.PathPrefix("/authors").Subrouter()
	// laguages := api.PathPrefix("/languages").Subrouter()
	// books := api.PathPrefix("/books").Subrouter()
	// quotes := api.PathPrefix("/quotes").Subrouter()
	return
}
