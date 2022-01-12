package quote

import (
	"testing"
)

const (
	filename       = "./test.sqlite"
	wrongStmt      = "wrong stmt on dao:\nexpected: %v\nactual: %v\n"
	wrongLen       = "wrong amount of daos:\nexpected: %d\nactual: %d\n"
	wrongId        = "wrong id of dao:\nexpected: %d\nactual: %d\n"
	wrongInsertion = "commiting a new dao returned non 0 id:\n got: %d\n"
)

func TestGetTopics(t *testing.T) {
	// Arrange
	database, err := Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	topics, err := database.GetTopics()
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedStmt := database.updateTopicStmt
	for _, topic := range topics {
		actualStmt := topic.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(wrongStmt, expectedStmt, actualStmt)
		}
	}
	// TODO: check database content which should be consistent
	expectedLen := 2
	actualLen := len(topics)
	if actualLen != expectedLen {
		t.Fatalf(wrongLen, expectedLen, actualLen)
	}
}

func TestGetNonExistingTopic(t *testing.T) {
	// Arrange
	database, err := Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	topic, err := database.GetTopic(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	// TODO: check database content as this should not exist
	defaultTopic := Topic{}
	if topic != defaultTopic {
		t.Fatal("Got non default topic for non existing topic id")
	}
}

func TestGetExistingTopic(t *testing.T) {
	// Arrange
	database, err := Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	topic, err := database.GetTopic(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedStmt := database.updateTopicStmt
	actualStmt := topic.stmt
	if actualStmt != expectedStmt {
		t.Fatalf(wrongStmt, expectedStmt, actualStmt)
	}
	// TODO: check database content should be consistent
	actualId := topic.id
	if actualId != expectedId {
		t.Fatalf(wrongId, expectedId, actualId)
	}
}

func TestRelatedBooksOfNonExistingTopic(t *testing.T) {
	// Arrange
	database, err := Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	books, err := database.RelatedBooksOfTopic(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 0 {
		t.Fatal("Found related books for non existing topic")
	}
}

func TestRelatedBooksOfExistingTopic(t *testing.T) {
	// Arrange
	database, err := Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	books, err := database.RelatedBooksOfTopic(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Check database
	expectedStmt := database.updateBookStmt
	for _, book := range books {
		actualStmt := book.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(wrongStmt, expectedStmt, actualStmt)
		}
		actualStmt = book.Topic.stmt
		expectedStmt = database.updateTopicStmt
		if actualStmt != expectedStmt {
			t.Fatalf(wrongStmt, expectedStmt, actualStmt)
		}
		actualStmt = book.Author.stmt
		expectedStmt = database.updateAuthorStmt
		if actualStmt != expectedStmt {
			t.Fatalf(wrongStmt, expectedStmt, actualStmt)
		}
		actualStmt = book.Language.stmt
		expectedStmt = database.updateLanguageStmt
		if actualStmt != expectedStmt {
			t.Fatalf(wrongStmt, expectedStmt, actualStmt)
		}
	}
}

func TestInsertNewTopic(t *testing.T) {
	// Arrange
	database, err := Connect(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	topic := database.NewTopic()
	topic.Topic = "Test Topic"
	// Act
	id, err := topic.Commit()
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if id == 0 {
		t.Fatalf(wrongInsertion, id)
	}
}
