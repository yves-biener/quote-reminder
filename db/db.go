package quote

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DAO interface {
	Commit() (sql.Result, error)
}

type Quote struct {
	id         int
	Book       Book
	Quote      string
	Page       int
	RecordDate time.Time
	stmt       *sql.Stmt
}

func (quote Quote) Commit() (res sql.Result, err error) {
	if quote.id == 0 { // Insert
		quote.Book.Commit()
		res, err = quote.stmt.Exec(quote.Book.id, quote.Quote, quote.Page)
	} else { // Update
		res, err = quote.stmt.Exec(quote.Book.id, quote.Quote, quote.Page, quote.id)
	}
	return
}

type Book struct {
	id          int
	Author      Author
	Topic       Topic
	Title       string
	ISBN        string
	Language    Language
	ReleaseDate time.Time
	stmt        *sql.Stmt
}

func (book Book) Commit() (res sql.Result, err error) {
	if book.id == 0 { // Insert
		book.Author.Commit()
		book.Topic.Commit()
		book.Language.Commit()
		res, err = book.stmt.Exec(book.Author.id, book.Topic.id, book.Title, book.ISBN, book.Language.id, book.ReleaseDate)
	} else { // Update
		res, err = book.stmt.Exec(book.Author.id, book.Topic.id, book.Title, book.ISBN, book.Language.id, book.ReleaseDate, book.id)
	}
	return
}

type Author struct {
	id   int
	Name string
	stmt *sql.Stmt
}

func (author Author) Commit() (res sql.Result, err error) {
	if author.id == 0 { // Insert
		res, err = author.stmt.Exec(author.Name)
	} else { // Update
		res, err = author.stmt.Exec(author.Name, author.id)
	}
	return
}

type Topic struct {
	id    int
	Topic string
	stmt  *sql.Stmt
}

func (topic Topic) Commit() (res sql.Result, err error) {
	if topic.id == 0 { // Insert
		res, err = topic.stmt.Exec(topic.Topic)
	} else { // Update
		res, err = topic.stmt.Exec(topic.Topic, topic.id)
	}
	return
}

type Language struct {
	id       int
	Language string
	stmt     *sql.Stmt
}

func (language Language) Commit() (res sql.Result, err error) {
	if language.id == 0 { // Insert
		res, err = language.stmt.Exec(language.Language)
	} else { // Update
		res, err = language.stmt.Exec(language.Language, language.id)
	}
	return
}

type database struct {
	// do I even need this? Or can I just use sql.Stmt?
	connection *sql.DB
	// insert statements
	insertBookStmt     *sql.Stmt
	insertTopicStmt    *sql.Stmt
	insertAuthorStmt   *sql.Stmt
	insertQuoteStmt    *sql.Stmt
	insertLanguageStmt *sql.Stmt
	// update statements
	updateBookStmt     *sql.Stmt
	updateTopicStmt    *sql.Stmt
	updateAuthorStmt   *sql.Stmt
	updateQuoteStmt    *sql.Stmt
	updateLanguageStmt *sql.Stmt
}

// Connect to an sqlite database located at `filename` This function ensures
// that the file will be created if it does not exist, create the required
// tables if it can successfully open the file
func Connect(filename string) (db *database) {
	var err error
	db = new(database)
	db.connection, err = sql.Open("sqlite3", filename)
	defer db.connection.Close()
	if err != nil {
		log.Fatalf("Could not open database: '%s' due to '%s'\n", filename, err)
	}
	db.Init()
	db.Prepare()
	return
}

// create tables
const (
	createBook = `CREATE TABLE Books (
Id int IDENTITY(1,1) PRIMARY KEY,
AuthorId int,
TopicId int,
ISBN NOT NULL UNIQUE,
Title varchar NOT NULL,
LanguageId int,
ReleaseDate date NOT NULL,
FOREIGN KEY (AuthorId) REFERENCES Authors(Id),
FOREIGN KEY (TopicId) REFERENCES Topics(Id),
FOREIGN KEY (LanguageId) REFERENCES Languages(Id)
);`
	createTopic = `CREATE TABLE Topics (
Id int IDENTITY(1,1) PRIMARY KEY,
Topic varchar NOT NULL UNIQUE
);`
	createAuthor = `CREATE TABLE Authors (
Id int IDENTITY(1,1) PRIMARY KEY,
Name varchar NOT NULL UNIQUE
);`
	createQuote = `CREATE TABLE Quotes (
Id int IDENTITY(1,1) PRIMARY KEY,
BbookId int,
Quote varchar NOT NULL,
Page int NOT NULL,
RecordDate date NOT NULL DEFAULT CURRENT_DATE,
FOREIGN KEY (BbookId) REFERENCES Books(Id)
);`
	createLanguage = `CREATE TABLE Languages (
Id int IDENTITY(1,1) PRIMARY KEY,
Language varchar NOT NULL UNIQUE
);`
)

// Initialize the database by creating the tables required for quote.
func (db *database) Init() {
	// create tables
	db.connection.Exec(createTopic)
	db.connection.Exec(createAuthor)
	db.connection.Exec(createLanguage)
	db.connection.Exec(createBook)
	db.connection.Exec(createQuote)
}

// Prepare Statements
const (
	insertBook     = "INSERT INTO Books (AuthorId, TopicId, ISBN, Title, LanguageId, ReleaseDate) VALUES (?, ?, ?, ?, ?, ?);"
	insertTopic    = "INSERT INTO Topics (Topic) VALUES (?);"
	insertAuthor   = "INSERT INTO Authors (Name) VALUES (?);"
	insertQuote    = "INSERT INTO Quotes (BbookId, Quote, Page) VALUES (?, ?, ?);"
	insertLanguage = "INSERT INTO Languages (Language) VALUES (?);"
)

const (
	updateBook     = "UPDATE Books SET AuthorId = ?, TopicId = ?, ISBN = ?, Title = ?, LanguageId = ?, ReleaseDate = ? WHERE Id = ?;"
	updateTopic    = "UPDATE Topics SET Topic = ? WHERE Id = ?;"
	updateAuthor   = "UPDATE Authors SET NAME = ? WHERE Id = ?;"
	updateQuote    = "UPDATE Quotes SET BbookId = ?, Quote = ?, Page = ? WHERE Id = ?;"
	updateLanguage = "UPDATE Languages SET Language = ? WHERE Id = ?;"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Prepare the queries used for the tables created by `Init'.
func (db *database) Prepare() {
	var err error
	// insert statements
	db.insertTopicStmt, err = db.connection.Prepare(insertTopic)
	checkError(err)
	db.insertAuthorStmt, err = db.connection.Prepare(insertAuthor)
	checkError(err)
	db.insertLanguageStmt, err = db.connection.Prepare(insertLanguage)
	checkError(err)
	db.insertBookStmt, err = db.connection.Prepare(insertBook)
	checkError(err)
	db.insertQuoteStmt, err = db.connection.Prepare(insertQuote)
	checkError(err)

	// update statements
	db.updateTopicStmt, err = db.connection.Prepare(updateTopic)
	checkError(err)
	db.updateAuthorStmt, err = db.connection.Prepare(updateAuthor)
	checkError(err)
	db.updateLanguageStmt, err = db.connection.Prepare(updateLanguage)
	checkError(err)
	db.updateBookStmt, err = db.connection.Prepare(updateBook)
	checkError(err)
	db.updateQuoteStmt, err = db.connection.Prepare(updateQuote)
	checkError(err)
}

func (db database) GetTopicById(id int) (topic Topic, err error) {
	return
}

func (db database) GetTopics() (topics []Topic, err error) {
	return
}

func (db database) GetAuthorById(id int) (author Author, err error) {
	return
}

func (db database) GetAuthors() (authors []Author, err error) {
	return
}

func (db database) GetLanguageById(id int) (language Language, err error) {
	return
}

func (db database) GetLanguages() (languages []Language, err error) {
	return
}

func (db database) GetBookById() (book Book, err error) {
	return
}

func (db database) GetBooks() (books []Book, err error) {
	return
}

func (db database) GetQuoteById() (quote Quote, err error) {
	return
}

func (db database) GetQuotes() (quotes []Quote, err error) {
	return
}
