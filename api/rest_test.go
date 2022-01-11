package quote

import (
	"net/http"
	"net/http/httptest"
	db "quote/db"
	"testing"
)

const (
	filename    = "./test.sqlite"
	statusError = "handler returned wrong status code:\nexpected: %v\nactual: %v\n"
	bodyError   = "handler returned wrong body:\nexpected: %v\nactual: %v\n"
)

func TestHelp(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
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
	req, err := http.NewRequest(Get, "/topics", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
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
	// TODO: Check response body
	// This requires a test database to be connected
}

func TestSearchTopics(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/topics?q=te st", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(searchTopics)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)

	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	// TODO: Check response body
	// This requires a test database to be connected
}

func TestGetTopicOfUnknownId(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/topics/69", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getTopic)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)

	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	// TODO: Check response body
	// This requires a test database to be connected
}

func TestGetTopicOfKnownId(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/topics/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getTopic)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	// TODO: Check response body
	// This requires a test database to be connected
}

func TestGetRelatedBooksOfUnknownTopic(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/topics/69/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getRelatedBooksOfTopic)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	// TODO: Check response body
	// This requires a test database to be connected
}

func TestGetRelatedBooksOfKnownTopic(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/topics/1/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getRelatedBooksOfTopic)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	// TODO: Check response body
	// This requires a test database to be connected
}

func TestGetRelatedQuotesOfUnknownTopic(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/topics/69/quotes", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getRelatedQuotesOfTopic)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	// TODO: Check response body
	// This requires a test database to be connected
}

func TestGetRelatedQuotesOfKnownTopic(t *testing.T) {
	// Arrange
	req, err := http.NewRequest(Get, "/topics/1/quotes", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(getRelatedQuotesOfTopic)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusOK
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	// TODO: Check response body
	// This requires a test database to be connected
}

func TestPostTopic(t *testing.T) {
	// Arrange
	// TODO: Check where I have to create the post message content
	req, err := http.NewRequest(Post, "/topics", nil)
	if err != nil {
		t.Fatal(err)
	}
	database, err = db.Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	responseRecord := httptest.NewRecorder()
	handlerUnderTest := http.HandlerFunc(postTopic)
	// Act
	handlerUnderTest.ServeHTTP(responseRecord, req)
	// Assert
	expectedStatus := http.StatusCreated
	if actualStatus := responseRecord.Code; actualStatus != expectedStatus {
		t.Errorf(statusError, expectedStatus, actualStatus)
	}
	// TODO: Check that test database will add this accordingly
	// TODO: Revert the database afterwards for consistent results
	expectedBody := `{"id": 3}`
	actualBody := responseRecord.Body.String()
	if actualBody != expectedBody {
		t.Errorf(bodyError, expectedBody, actualBody)
	}
}
