package quote

import (
	"database/sql"
	"fmt"
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
	// database statements
	createDatabaseStmt *sql.Stmt
	useDatabaseStmt    *sql.Stmt
	// create statements
	createBookStmt     *sql.Stmt
	createTopicStmt    *sql.Stmt
	createAuthorStmt   *sql.Stmt
	createQuoteStmt    *sql.Stmt
	createLanguageStmt *sql.Stmt
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
	if err != nil {
		log.Fatalf("Could not open database: '%s' due to '%s'\n", filename, err)
	}
	db.Prepare()
	db.Init()
	return
}

// Prepare Statements
const (
	createDatabase = "IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'QuoteDB') BEGIN CREATE DATABASE QuoteDB END"
	useDatabase    = "USE QuoteDB;"
)

const (
	createBook = `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'Books' and xtype='U')
BEGIN
CREATE TABLE Books (
Id int IDENTITY(1,1) PRIMARY KEY,
AuthorId int FOREIGN KEY REFERENCES Authors(Id),
TopicId int FOREIGN KEY REFERENCES Topics(Id),
ISBN NOT NULL UNIQUE,
Title varchar NOT NULL,
LanguageId int FOREIGN KEY REFERENCES Languages(Id),
ReleaseDate date NOT NULL
)
END`
	createTopic = `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'Topics' and xtype='U')
BEGIN
CREATE TABLE Topics (
Id int IDENTITY(1,1) PRIMARY KEY,
Topic varchar NOT NULL UNIQUE
)
END`
	createAuthor = `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'Authors' and xtype='U')
BEGIN
CREATE TABLE Authors (
Id int IDENTITY(1,1) PRIMARY KEY,
Name varchar NOT NULL UNIQUE
)
END`
	createQuote = `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'Quotes' and xtype='U')
BEGIN
CREATE TABLE Quotes (
Id int IDENTITY(1,1) PRIMARY KEY,
BookId int FOREIGN KEY REFERENCES Books(Id),
Quote varchar NOT NULL,
Page int NOT NULL,
RecordDate date NOT NULL DEFAULT CURRENT_DATE,
)
END`
	createLanguage = `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'Languages' and xtype='U')
BEGIN
CREATE TABLE Languages (
Id int IDENTITY(1,1) PRIMARY KEY,
Language varchar NOT NULL UNIQUE
)
END`
)

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

func (db *database) Prepare() {
	var err error

	// database statements
	db.createDatabaseStmt, err = db.connection.Prepare(createDatabase)
	checkError(err)
	db.useDatabaseStmt, err = db.connection.Prepare(useDatabase)
	checkError(err)

	// create statements
	db.createTopicStmt, err = db.connection.Prepare(createTopic)
	checkError(err)
	db.createAuthorStmt, err = db.connection.Prepare(createAuthor)
	checkError(err)
	db.createLanguageStmt, err = db.connection.Prepare(createLanguage)
	checkError(err)
	db.createBookStmt, err = db.connection.Prepare(createBook)
	checkError(err)
	db.createQuoteStmt, err = db.connection.Prepare(createQuote)
	checkError(err)

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

func (db *database) Init() {
	var err error
	var res sql.Result

	res, err = db.createDatabaseStmt.Exec()
	fmt.Println(res)
	checkError(err)
	// not sure if I should close them after I
	// don't need the anymore...
	db.createDatabaseStmt.Close()
	res, err = db.useDatabaseStmt.Exec()
	fmt.Println(res)
	checkError(err)
	db.useDatabaseStmt.Close()

	res, err = db.createTopicStmt.Exec()
	fmt.Println(res)
	checkError(err)
	db.createTopicStmt.Close()

	res, err = db.createLanguageStmt.Exec()
	fmt.Println(res)
	checkError(err)
	db.createLanguageStmt.Close()

	res, err = db.createAuthorStmt.Exec()
	fmt.Println(res)
	checkError(err)
	db.createAuthorStmt.Close()

	res, err = db.createBookStmt.Exec()
	fmt.Println(res)
	checkError(err)
	db.createBookStmt.Close()

	res, err = db.createQuoteStmt.Exec()
	fmt.Println(res)
	checkError(err)
	db.createQuoteStmt.Close()
}
