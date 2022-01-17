package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	db "quote/db"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

const (
	testSource   = "./../test.sqlite"
	testDatabase = "./cur_test.sqlite"
	statusError  = "Handler returned wrong status code:\nexpected: %v\nactual: %v\n"
	bodyError    = "Handler returned wrong body:\nexpected: %v\nactual: %v\n"
	headerError  = "Middleware returned wrong header:\nexpected: %v\nactual: %v\n"
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

func TestHeaderMiddleware(t *testing.T) {
	// Arrange
	router := mux.NewRouter()
	router.Use(jsonContentWrapper)
	req, err := http.NewRequest(Get, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecord := httptest.NewRecorder()
	expectedBody := `{"Test": "success"}`
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(expectedBody))
	}
	router.HandleFunc("/", testHandler)
	// Act
	router.ServeHTTP(responseRecord, req)
	// Assert
	expectedHeader := "application/json"
	actualHeader := responseRecord.Header().Get("Content-type")
	if actualHeader == "" {
		t.Error(`Returned header for "Content-type" was empty`)
	} else if actualHeader != expectedHeader {
		t.Errorf(headerError, expectedHeader, actualHeader)
	}
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestHelp(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(help)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	// I do not check the respose body as it will frequently change
}

func TestGetTopics(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getTopics)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedTopics, err := database.GetTopics()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, topic := range expectedTopics {
		jsonTopic, err := json.Marshal(topic)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonTopic))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestSearchTopics(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/?q=Topic", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", searchTopics).Queries("q", "{search}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedTopics, err := database.GetTopics()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, topic := range expectedTopics {
		jsonTopic, err := json.Marshal(topic)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonTopic))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetTopicOfUnknownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/69", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getTopic)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusNotFound
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := ""
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetTopicOfKnownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getTopic)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedTopic, err := database.GetTopic(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedTopic)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := string(expectedJson)
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfUnknownTopic(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfTopic)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfKnownTopic(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/books", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfTopic)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBook, err := database.GetBook(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedBook)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfUnknownTopicFiltered(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/books?q=Book", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfTopic).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfKnownTopicFiltered(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/books?q=Book", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfTopic).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBook, err := database.GetBook(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedBook)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfUnknownTopic(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/quotes", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfTopic)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfKnownTopic(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/quotes", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfTopic)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuote, err := database.GetQuote(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedQuote)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfUnknownTopicFiltered(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/quotes?q=Quote", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfTopic).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfKnownTopicFiltered(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/quotes?q=Quote", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfTopic).Queries("q", "{filtered}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuote, err := database.GetQuote(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedQuote)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

// TODO: add post form variables
func TestPostTopic(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Post, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", postTopic).Methods(Post)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusCreated
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := `{"Id": 3}`
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetAuthors(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getAuthors)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedAuthors, err := database.GetAuthors()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, author := range expectedAuthors {
		jsonAuthor, err := json.Marshal(author)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonAuthor))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestSearchAuthors(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/?q=Author", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", searchAuthors).Queries("q", "{search}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedAuthors, err := database.GetAuthors()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, author := range expectedAuthors {
		jsonAuthor, err := json.Marshal(author)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonAuthor))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetAuthorOfUnknownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/69", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getAuthor)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusNotFound
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := ""
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetAuthorOfKnownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getAuthor)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedAuthor, err := database.GetAuthor(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedAuthor)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := string(expectedJson)
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfUnknownAuthor(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfAuthor)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfKnownAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/books", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfAuthor)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBook, err := database.GetBook(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedBook)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfUnknownAuthorFiltered(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/books?q=Book", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfAuthor).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfKnownAuthorFiltered(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/books?q=Book", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfAuthor).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBook, err := database.GetBook(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedBook)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfUnknownAuthor(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/quotes", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfAuthor)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfKnownAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/quotes", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfAuthor)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuote, err := database.GetQuote(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedQuote)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfUnknownAuthorFiltered(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/quotes?q=Quote", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfAuthor).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfKnownAuthorFiltered(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/quotes?q=Quote", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfAuthor).Queries("q", "{filtered}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuote, err := database.GetQuote(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedQuote)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

// TODO: add post form variables
func TestPostAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Post, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", postAuthor).Methods(Post)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusCreated
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := `{"Id": 3}`
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetLanguages(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getLanguages)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedLanguages, err := database.GetLanguages()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, language := range expectedLanguages {
		jsonLanguage, err := json.Marshal(language)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonLanguage))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestSearchLanguages(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/?q=Language", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", searchLanguages).Queries("q", "{search}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedLanguages, err := database.GetLanguages()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, language := range expectedLanguages {
		jsonLanguage, err := json.Marshal(language)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonLanguage))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetLanguageOfUnknownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/69", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getLanguage)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusNotFound
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := ""
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetLanguageOfKnownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getLanguage)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedLanguage, err := database.GetLanguage(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedLanguage)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := string(expectedJson)
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfUnknownLanguage(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfLanguage)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfKnownLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/books", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfLanguage)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBook, err := database.GetBook(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedBook)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfUnknownLanguageFiltered(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/books?q=Book", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfLanguage).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedBooksOfKnownLanguageFiltered(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/books?q=Book", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/books", getRelatedBooksOfLanguage).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBook, err := database.GetBook(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedBook)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfUnknownLanguage(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/quotes", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfLanguage)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfKnownLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/quotes", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfLanguage)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuote, err := database.GetQuote(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedQuote)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfUnknownLanguageFiltered(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/quotes?q=Quote", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfLanguage).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfKnownLanguageFiltered(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/quotes?q=Quote", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfLanguage).Queries("q", "{filtered}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuote, err := database.GetQuote(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedQuote)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

// TODO: add post form variables
func TestPostLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Post, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", postLanguage).Methods(Post)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusCreated
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := `{"Id": 3}`
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetBooks(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getBooks)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBooks, err := database.GetBooks()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, book := range expectedBooks {
		jsonBook, err := json.Marshal(book)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonBook))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestSearchBooks(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/?q=Book", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", searchBooks).Queries("q", "{search}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBooks, err := database.GetBooks()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, book := range expectedBooks {
		jsonBook, err := json.Marshal(book)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonBook))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetBookOfUnknownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/69", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getBook)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusNotFound
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := ""
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetBookOfKnownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getBook)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBook, err := database.GetBook(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedBook)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := string(expectedJson)
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfUnknownBook(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/quotes", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfBook)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfKnownBook(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/quotes", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfBook)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuote, err := database.GetQuote(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedQuote)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfUnknownBookFiltered(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/69/quotes?q=Quote", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfBook).Queries("q", "{filter}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := "[]"
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetRelatedQuotesOfKnownBookFiltered(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d/quotes?q=Quote", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}/quotes", getRelatedQuotesOfBook).Queries("q", "{filtered}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuote, err := database.GetQuote(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedQuote)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := fmt.Sprintf("[%s]", string(expectedJson))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

// TODO: add post form variables
func TestPostBook(t *testing.T) {
	// Arrange
	initDatabase(t)
	data := url.Values{}
	data.Add("AuthorId", "1")
	data.Add("TopicId", "1")
	data.Add("LanguageId", "1")
	req, err := http.NewRequest(Post, "/", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", postBook).Methods(Post)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusCreated
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := `{"Id": 3}`
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetQuotes(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getQuotes)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuotes, err := database.GetQuotes()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, quote := range expectedQuotes {
		jsonQuote, err := json.Marshal(quote)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonQuote))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestSearchQuotes(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/?q=Quote", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", searchQuotes).Queries("q", "{search}")
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuotes, err := database.GetQuotes()
	if err != nil {
		t.Fatal(err)
	}
	var expectedJson []string
	for _, quote := range expectedQuotes {
		jsonQuote, err := json.Marshal(quote)
		if err != nil {
			t.Fatal(err)
		}
		expectedJson = append(expectedJson, string(jsonQuote))
	}
	expectedBody := fmt.Sprintf("[%s]", strings.Join(expectedJson, ","))
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetQuoteOfUnknownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	req, err := http.NewRequest(Get, "/69", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getQuote)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusNotFound
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := ""
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

func TestGetQuoteOfKnownId(t *testing.T) {
	// Arrange
	initDatabase(t)
	expectedId := 1
	req, err := http.NewRequest(Get, fmt.Sprintf("/%d", expectedId), nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/{id}", getQuote)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedQuote, err := database.GetQuote(expectedId)
	if err != nil {
		t.Fatal(err)
	}
	expectedJson, err := json.Marshal(expectedQuote)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := string(expectedJson)
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}

// TODO: add post form variables
func TestPostQuote(t *testing.T) {
	// Arrange
	initDatabase(t)
	data := url.Values{}
	data.Add("BookId", "1")
	req, err := http.NewRequest(Post, "/", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	database, err = db.Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	routerUnderTest := mux.NewRouter()
	routerUnderTest.HandleFunc("/", postQuote).Methods(Post)
	// Act
	routerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusCreated
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	expectedBody := `{"Id": 3}`
	if actualBody := responseRecord.Body.String(); actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}
