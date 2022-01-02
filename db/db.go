package quote

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	Id         int
	BookId     int
	Value      string
	Page       int
	RecordDate time.Time
}

type Book struct {
	Id          int
	AuthorId    int
	TopicId     int
	LanguageId  int
	Title       string
	ISBN        string
	ReleaseDate time.Time
}

type Author struct {
	Id   int
	Name string
}

type Topic struct {
	Id    int
	Topic string
}

type Language struct {
	Id       int
	Language string
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
BookId int,
Quote varchar NOT NULL,
Page int NOT NULL,
RecordDate date NOT NULL DEFAULT CURRENT_DATE,
FOREIGN KEY (BookId) REFERENCES Books(Id)
);`
	createLanguage = `CREATE TABLE Languages (
Id int IDENTITY(1,1) PRIMARY KEY,
Language varchar NOT NULL UNIQUE
);`
)

// Prepare Statements
const (
	insertBook     = "INSERT INTO Books (AuthorId, TopicId, ISBN, Title, LanguageId, ReleaseDate) VALUES (?, ?, ?, ?, ?, ?);"
	insertTopic    = "INSERT INTO Topics (Topic) VALUES (?);"
	insertAuthor   = "INSERT INTO Authors (Name) VALUES (?);"
	insertQuote    = "INSERT INTO Quotes (BookId, Quote, Page) VALUES (?, ?, ?);"
	insertLanguage = "INSERT INTO Languages (Language) VALUES (?);"
)

const (
	updateBook     = "UPDATE Books SET AuthorId = ?, TopicId = ?, ISBN = ?, Title = ?, LanguageId = ?, ReleaseDate = ? WHERE Id = ?;"
	updateTopic    = "UPDATE Topics SET Topic = ? WHERE Id = ?;"
	updateAuthor   = "UPDATE Authors SET NAME = ? WHERE Id = ?;"
	updateQuote    = "UPDATE Quotes SET Quote = ?, Page = ? WHERE Id = ?;"
	updateLanguage = "UPDATE Languages SET Language = ? WHERE Id = ?;"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Initialize the database by creating the tables required for quote.
func (db *database) Init() {
	// create tables
	db.connection.Exec(createTopic)
	db.connection.Exec(createAuthor)
	db.connection.Exec(createLanguage)
	db.connection.Exec(createBook)
	db.connection.Exec(createQuote)
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
