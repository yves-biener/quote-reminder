package quote

import (
	"fmt"
	"io"
	"os"
	"testing"
)

const (
	testSource     = "./../test.sqlite"
	testDatabase   = "./cur_test.sqlite"
	stmtError      = "wrong stmt on dao:\nexpected: %v\nactual: %v\n"
	lenError       = "wrong amount of daos:\nexpected: %d\nactual: %d\n"
	idError        = "wrong id of dao:\nexpected: %d\nactual: %d\n"
	insertionError = "commiting new dao returned unexpected id\nexpected: %d\n got: %d\n"
	contentError   = "content had not the expected value\nexpected: %v\n actual: %v\n"
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

func TestGetTopics(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
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
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
	expectedLen := 2
	actualLen := len(topics)
	if actualLen != expectedLen {
		t.Fatalf(lenError, expectedLen, actualLen)
	}
	for i, topic := range topics {
		expectedTopic := fmt.Sprintf("Topic%d", i+1)
		actualTopic := topic.Topic
		if actualTopic != expectedTopic {
			t.Fatalf(contentError, expectedTopic, actualTopic)
		}
	}
}

func TestGetNonExistingTopic(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
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
	defaultTopic := Topic{}
	if topic != defaultTopic {
		t.Fatal("Got non default topic for non existing topic id")
	}
}

func TestGetExistingTopic(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
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
		t.Fatalf(stmtError, expectedStmt, actualStmt)
	}
	actualId := topic.id
	if actualId != expectedId {
		t.Fatalf(idError, expectedId, actualId)
	}
	expectedTopic := fmt.Sprintf("Topic%d", expectedId)
	actualTopic := topic.Topic
	if actualTopic != expectedTopic {
		t.Fatalf(contentError, expectedTopic, actualTopic)
	}
}

func TestRelatedBooksOfNonExistingTopic(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
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
	initDatabase(t)
	database, err := Connect(testDatabase)
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
	expectedLen := 1
	actualLen := len(books)
	if actualLen != expectedLen {
		t.Fatalf(contentError, expectedLen, actualLen)
	}
	expectedStmt := database.updateBookStmt
	for _, book := range books {
		expectedContent := fmt.Sprintf("Book%d", expectedId)
		actualContent := book.Title
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		actualId := book.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualStmt := book.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = book.Topic.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Topic%d", expectedId)
		actualContent = book.Topic.Topic
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateTopicStmt
		actualStmt = book.Topic.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = book.Author.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Author%d", expectedId)
		actualContent = book.Author.Name
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateAuthorStmt
		actualStmt = book.Author.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = book.Language.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Language%d", expectedId)
		actualContent = book.Language.Language
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateLanguageStmt
		actualStmt = book.Language.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
}

func TestRelatedQuotesOfNonExistingTopic(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	quotes, err := database.RelatedQuotesOfTopic(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 0 {
		t.Fatal("Found related quotes for non existing topic")
	}
}

func TestRelatedQuotesOfExistingTopic(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	quotes, err := database.RelatedQuotesOfTopic(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedLen := 1
	actualLen := len(quotes)
	if actualLen != expectedLen {
		t.Fatalf(contentError, expectedLen, actualLen)
	}
	expectedStmt := database.updateQuoteStmt
	for _, quote := range quotes {
		expectedContent := fmt.Sprintf("Quote%d", expectedId)
		actualContent := quote.Quote
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedContent = fmt.Sprintf("Book%d", expectedId)
		actualContent = quote.Book.Title
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		actualId := quote.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualStmt := quote.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Topic.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Topic%d", expectedId)
		actualContent = quote.Book.Topic.Topic
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateTopicStmt
		actualStmt = quote.Book.Topic.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Author.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Author%d", expectedId)
		actualContent = quote.Book.Author.Name
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateAuthorStmt
		actualStmt = quote.Book.Author.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Language.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Language%d", expectedId)
		actualContent = quote.Book.Language.Language
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateLanguageStmt
		actualStmt = quote.Book.Language.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
}

func TestInsertNewTopic(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	topic := database.NewTopic()
	topic.Topic = "Test Topic"
	// Act
	actualId, err := topic.Commit()
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedId := 3
	if actualId != expectedId {
		t.Fatalf(insertionError, expectedId, actualId)
	}
}

func TestGetAuthors(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	authors, err := database.GetAuthors()
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedStmt := database.updateAuthorStmt
	for _, author := range authors {
		actualStmt := author.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
	expectedLen := 2
	actualLen := len(authors)
	if actualLen != expectedLen {
		t.Fatalf(lenError, expectedLen, actualLen)
	}
	for i, author := range authors {
		expectedAuthor := fmt.Sprintf("Author%d", i+1)
		actualAuthor := author.Name
		if actualAuthor != expectedAuthor {
			t.Fatalf(contentError, expectedAuthor, actualAuthor)
		}
	}
}

func TestGetNonExistingAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	author, err := database.GetAuthor(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	defaultAuthor := Author{}
	if author != defaultAuthor {
		t.Fatal("Got non default author for non existing author id")
	}
}

func TestGetExistingAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	author, err := database.GetAuthor(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedStmt := database.updateAuthorStmt
	actualStmt := author.stmt
	if actualStmt != expectedStmt {
		t.Fatalf(stmtError, expectedStmt, actualStmt)
	}
	actualId := author.id
	if actualId != expectedId {
		t.Fatalf(idError, expectedId, actualId)
	}
	expectedAuthor := fmt.Sprintf("Author%d", expectedId)
	actualAuthor := author.Name
	if actualAuthor != expectedAuthor {
		t.Fatalf(contentError, expectedAuthor, actualAuthor)
	}
}

func TestRelatedBooksOfNonExistingAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	books, err := database.RelatedBooksOfAuthor(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 0 {
		t.Fatal("Found related books for non existing author")
	}
}

func TestRelatedBooksOfExistingAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	books, err := database.RelatedBooksOfAuthor(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedLen := 1
	actualLen := len(books)
	if actualLen != expectedLen {
		t.Fatalf(contentError, expectedLen, actualLen)
	}
	expectedStmt := database.updateBookStmt
	for _, book := range books {
		expectedContent := fmt.Sprintf("Book%d", expectedId)
		actualContent := book.Title
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		actualId := book.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualStmt := book.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = book.Topic.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Topic%d", expectedId)
		actualContent = book.Topic.Topic
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateTopicStmt
		actualStmt = book.Topic.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = book.Author.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Author%d", expectedId)
		actualContent = book.Author.Name
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateAuthorStmt
		actualStmt = book.Author.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = book.Language.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Language%d", expectedId)
		actualContent = book.Language.Language
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateLanguageStmt
		actualStmt = book.Language.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
}

func TestRelatedQuotesOfNonExistingAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	quotes, err := database.RelatedQuotesOfAuthor(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 0 {
		t.Fatal("Found related quotes for non existing author")
	}
}

func TestRelatedQuotesOfExistingAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	quotes, err := database.RelatedQuotesOfAuthor(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedLen := 1
	actualLen := len(quotes)
	if actualLen != expectedLen {
		t.Fatalf(contentError, expectedLen, actualLen)
	}
	expectedStmt := database.updateQuoteStmt
	for _, quote := range quotes {
		expectedContent := fmt.Sprintf("Quote%d", expectedId)
		actualContent := quote.Quote
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedContent = fmt.Sprintf("Book%d", expectedId)
		actualContent = quote.Book.Title
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		actualId := quote.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualStmt := quote.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Topic.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Topic%d", expectedId)
		actualContent = quote.Book.Topic.Topic
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateTopicStmt
		actualStmt = quote.Book.Topic.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Author.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Author%d", expectedId)
		actualContent = quote.Book.Author.Name
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateAuthorStmt
		actualStmt = quote.Book.Author.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Language.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Language%d", expectedId)
		actualContent = quote.Book.Language.Language
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateLanguageStmt
		actualStmt = quote.Book.Language.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
}

func TestInsertNewAuthor(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	author := database.NewAuthor()
	author.Name = "Test Author"
	// Act
	actualId, err := author.Commit()
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedId := 3
	if actualId != expectedId {
		t.Fatalf(insertionError, expectedId, actualId)
	}
}

func TestGetLanguages(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	languages, err := database.GetLanguages()
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedStmt := database.updateLanguageStmt
	for _, language := range languages {
		actualStmt := language.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
	expectedLen := 2
	actualLen := len(languages)
	if actualLen != expectedLen {
		t.Fatalf(lenError, expectedLen, actualLen)
	}
	for i, language := range languages {
		expectedLanguage := fmt.Sprintf("Language%d", i+1)
		actualLanguage := language.Language
		if actualLanguage != expectedLanguage {
			t.Fatalf(contentError, expectedLanguage, actualLanguage)
		}
	}
}

func TestGetNonExistingLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	language, err := database.GetLanguage(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	defaultLanguage := Language{}
	if language != defaultLanguage {
		t.Fatal("Got non default language for non existing language id")
	}
}

func TestGetExistingLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	language, err := database.GetLanguage(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedStmt := database.updateLanguageStmt
	actualStmt := language.stmt
	if actualStmt != expectedStmt {
		t.Fatalf(stmtError, expectedStmt, actualStmt)
	}
	actualId := language.id
	if actualId != expectedId {
		t.Fatalf(idError, expectedId, actualId)
	}
	expectedLanguage := fmt.Sprintf("Language%d", expectedId)
	actualLanguage := language.Language
	if actualLanguage != expectedLanguage {
		t.Fatalf(contentError, expectedLanguage, actualLanguage)
	}
}

func TestRelatedBooksOfNonExistingLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	books, err := database.RelatedBooksOfLanguage(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 0 {
		t.Fatal("Found related books for non existing language")
	}
}

func TestRelatedBooksOfExistingLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	books, err := database.RelatedBooksOfLanguage(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedLen := 1
	actualLen := len(books)
	if actualLen != expectedLen {
		t.Fatalf(contentError, expectedLen, actualLen)
	}
	expectedStmt := database.updateBookStmt
	for _, book := range books {
		expectedContent := fmt.Sprintf("Book%d", expectedId)
		actualContent := book.Title
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		actualId := book.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualStmt := book.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = book.Topic.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Topic%d", expectedId)
		actualContent = book.Topic.Topic
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateTopicStmt
		actualStmt = book.Topic.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = book.Author.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Author%d", expectedId)
		actualContent = book.Author.Name
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateAuthorStmt
		actualStmt = book.Author.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = book.Language.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Language%d", expectedId)
		actualContent = book.Language.Language
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateLanguageStmt
		actualStmt = book.Language.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
}

func TestRelatedQuotesOfNonExistingLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	quotes, err := database.RelatedQuotesOfLanguage(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 0 {
		t.Fatal("Found related quotes for non existing language")
	}
}

func TestRelatedQuotesOfExistingLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	quotes, err := database.RelatedQuotesOfLanguage(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedLen := 1
	actualLen := len(quotes)
	if actualLen != expectedLen {
		t.Fatalf(contentError, expectedLen, actualLen)
	}
	expectedStmt := database.updateQuoteStmt
	for _, quote := range quotes {
		expectedContent := fmt.Sprintf("Quote%d", expectedId)
		actualContent := quote.Quote
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedContent = fmt.Sprintf("Book%d", expectedId)
		actualContent = quote.Book.Title
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		actualId := quote.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualStmt := quote.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Topic.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Topic%d", expectedId)
		actualContent = quote.Book.Topic.Topic
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateTopicStmt
		actualStmt = quote.Book.Topic.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Author.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Author%d", expectedId)
		actualContent = quote.Book.Author.Name
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateAuthorStmt
		actualStmt = quote.Book.Author.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Language.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Language%d", expectedId)
		actualContent = quote.Book.Language.Language
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateLanguageStmt
		actualStmt = quote.Book.Language.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
}

func TestInsertNewLanguage(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	language := database.NewLanguage()
	language.Language = "Test Language"
	// Act
	actualId, err := language.Commit()
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedId := 3
	if actualId != expectedId {
		t.Fatalf(insertionError, expectedId, actualId)
	}
}

func TestGetBooks(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	books, err := database.GetBooks()
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedStmt := database.updateBookStmt
	for _, book := range books {
		actualStmt := book.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
	expectedLen := 2
	actualLen := len(books)
	if actualLen != expectedLen {
		t.Fatalf(lenError, expectedLen, actualLen)
	}
	for i, book := range books {
		expectedTitle := fmt.Sprintf("Book%d", i+1)
		actualTitle := book.Title
		if actualTitle != expectedTitle {
			t.Fatalf(contentError, expectedTitle, actualTitle)
		}
		actualISBN := book.ISBN.String
		if len(actualISBN) == 0 {
			t.Fatalf("ISBN was empty, the value should be not null")
		}
		expectedId := i + 1
		actualId := book.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualId = book.Author.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualId = book.Topic.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualId = book.Language.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
	}
}

func TestGetNonExistingBook(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	book, err := database.GetBook(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	defaultBook := Book{}
	if book != defaultBook {
		t.Fatal("Got non default book for non existing book id")
	}
}

func TestGetExistingBook(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	book, err := database.GetBook(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedStmt := database.updateBookStmt
	actualStmt := book.stmt
	if actualStmt != expectedStmt {
		t.Fatalf(stmtError, expectedStmt, actualStmt)
	}
	actualId := book.id
	if actualId != expectedId {
		t.Fatalf(idError, expectedId, actualId)
	}
	expectedTitle := fmt.Sprintf("Book%d", expectedId)
	actualTitle := book.Title
	if actualTitle != expectedTitle {
		t.Fatalf(contentError, expectedTitle, actualTitle)
	}
}

func TestRelatedQuotesOfNonExistingBook(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	quotes, err := database.RelatedQuotesOfBook(69)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 0 {
		t.Fatal("Found related quotes for non existing book")
	}
}

func TestRelatedQuotesOfExistingBook(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	// Act
	expectedId := 1
	quotes, err := database.RelatedQuotesOfBook(expectedId)
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedLen := 1
	actualLen := len(quotes)
	if actualLen != expectedLen {
		t.Fatalf(contentError, expectedLen, actualLen)
	}
	expectedStmt := database.updateQuoteStmt
	for _, quote := range quotes {
		expectedContent := fmt.Sprintf("Quote%d", expectedId)
		actualContent := quote.Quote
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedContent = fmt.Sprintf("Book%d", expectedId)
		actualContent = quote.Book.Title
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		actualId := quote.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		actualStmt := quote.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Topic.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Topic%d", expectedId)
		actualContent = quote.Book.Topic.Topic
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateTopicStmt
		actualStmt = quote.Book.Topic.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Author.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Author%d", expectedId)
		actualContent = quote.Book.Author.Name
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateAuthorStmt
		actualStmt = quote.Book.Author.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
		actualId = quote.Book.Language.id
		if actualId != expectedId {
			t.Fatalf(idError, expectedId, actualId)
		}
		expectedContent = fmt.Sprintf("Language%d", expectedId)
		actualContent = quote.Book.Language.Language
		if actualContent != expectedContent {
			t.Fatalf(contentError, expectedContent, actualContent)
		}
		expectedStmt = database.updateLanguageStmt
		actualStmt = quote.Book.Language.stmt
		if actualStmt != expectedStmt {
			t.Fatalf(stmtError, expectedStmt, actualStmt)
		}
	}
}

func TestInsertNewBook(t *testing.T) {
	// Arrange
	initDatabase(t)
	database, err := Connect(testDatabase)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	relatedId := 1
	author, err := database.GetAuthor(relatedId)
	if err != nil {
		t.Fatal(err)
	}
	language, err := database.GetLanguage(relatedId)
	if err != nil {
		t.Fatal(err)
	}
	topic, err := database.GetTopic(relatedId)
	if err != nil {
		t.Fatal(err)
	}
	book := database.NewBook(author, topic, language)
	book.Title = "Test Book"
	book.ISBN.Scan("965-17-650-5-10")
	// Act
	actualId, err := book.Commit()
	// Assert
	if err != nil {
		t.Fatal(err)
	}
	expectedId := 3
	if actualId != expectedId {
		t.Fatalf(insertionError, expectedId, actualId)
	}
}
