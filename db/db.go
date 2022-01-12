package quote

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DAO interface {
	// Commit changes of the DAO object to the Database, returning the
	// associated id of the DAO or an error if Commit failed
	Commit() (int, error)
}

type Quote struct {
	id         int
	Book       Book
	Quote      string
	Page       int
	RecordDate time.Time
	stmt       *sql.Stmt
}

func (db Database) NewQuote(book Book) (quote Quote) {
	quote.stmt = db.insertQuoteStmt
	quote.Book = book
	return
}

func (quote Quote) Commit() (id int, err error) {
	if quote.id == 0 { // Insert
		quote.Book.id, err = quote.Book.Commit()
		if err != nil {
			return -1, err
		}
		res, err := quote.stmt.Exec(quote.Book.id, quote.Quote, quote.Page)
		if err != nil {
			return -1, err
		}
		insertedId, e := res.LastInsertId()
		id = int(insertedId)
		err = e
	} else { // Update
		_, err = quote.stmt.Exec(quote.Book.id, quote.Quote, quote.Page, quote.id)
		id = quote.id
	}
	return
}

type Book struct {
	id          int
	Author      Author
	Topic       Topic
	Title       string
	ISBN        sql.NullString
	Language    Language
	ReleaseDate time.Time
	stmt        *sql.Stmt
}

func (db Database) NewBook(author Author, topic Topic, language Language) (book Book) {
	book.stmt = db.insertBookStmt
	book.Author = author
	book.Topic = topic
	book.Language = language
	return
}

func (book Book) Commit() (id int, err error) {
	if book.id == 0 { // Insert
		book.Author.id, err = book.Author.Commit()
		if err != nil {
			return
		}
		book.Topic.id, err = book.Topic.Commit()
		if err != nil {
			return
		}
		book.Language.id, err = book.Language.Commit()
		if err != nil {
			return
		}
		res, err := book.stmt.Exec(book.Author.id, book.Topic.id, book.Title, book.ISBN, book.Language.id, book.ReleaseDate)
		if err != nil {
			return -1, err
		}
		insertedId, e := res.LastInsertId()
		id = int(insertedId)
		err = e
	} else { // Update
		_, err = book.stmt.Exec(book.Author.id, book.Topic.id, book.Title, book.ISBN, book.Language.id, book.ReleaseDate, book.id)
		id = book.id
	}
	return
}

type Author struct {
	id   int
	Name string
	stmt *sql.Stmt
}

func (db Database) NewAuthor() (author Author) {
	author.stmt = db.insertAuthorStmt
	return
}

func (author *Author) Commit() (id int, err error) {
	if author.id == 0 { // Insert
		res, err := author.stmt.Exec(author.Name)
		if err != nil {
			return -1, err
		}
		insertedId, e := res.LastInsertId()
		id = int(insertedId)
		err = e
	} else { // Update
		_, err = author.stmt.Exec(author.Name, author.id)
		id = author.id
	}
	return
}

type Topic struct {
	id    int
	Topic string
	stmt  *sql.Stmt
}

func (db Database) NewTopic() (topic Topic) {
	topic.stmt = db.insertTopicStmt
	return
}

func (topic Topic) Commit() (id int, err error) {
	if topic.id == 0 { // Insert
		res, err := topic.stmt.Exec(topic.Topic)
		if err != nil {
			return -1, err
		}
		insertedId, e := res.LastInsertId()
		id = int(insertedId)
		err = e
	} else { // Update
		_, err = topic.stmt.Exec(topic.Topic, topic.id)
		id = topic.id
	}
	return
}

type Language struct {
	id       int
	Language string
	stmt     *sql.Stmt
}

func (db Database) NewLanguage() (language Language) {
	language.stmt = db.insertLanguageStmt
	return
}

func (language Language) Commit() (id int, err error) {
	if language.id == 0 { // Insert
		res, err := language.stmt.Exec(language.Language)
		if err != nil {
			return -1, err
		}
		insertedId, e := res.LastInsertId()
		id = int(insertedId)
		err = e
	} else { // Update
		_, err = language.stmt.Exec(language.Language, language.id)
		id = language.id
	}
	return
}

type Database struct {
	connection *sql.DB
	// select statements
	selectBooksStmt     *sql.Stmt
	selectTopicsStmt    *sql.Stmt
	selectAuthorsStmt   *sql.Stmt
	selectQuotesStmt    *sql.Stmt
	selectLanguagesStmt *sql.Stmt
	// select by id statements
	selectBookStmt     *sql.Stmt
	selectTopicStmt    *sql.Stmt
	selectAuthorStmt   *sql.Stmt
	selectQuoteStmt    *sql.Stmt
	selectLanguageStmt *sql.Stmt
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
	// related entries statements
	relatedQuotesOfBookStmt     *sql.Stmt
	relatedBooksOfTopicStmt     *sql.Stmt
	relatedQuotesOfTopicStmt    *sql.Stmt
	relatedBooksOfAuthorStmt    *sql.Stmt
	relatedQuotesOfAuthorStmt   *sql.Stmt
	relatedBooksOfLanguageStmt  *sql.Stmt
	relatedQuotesOfLanguageStmt *sql.Stmt
	// searches
	searchTopicsStmt    *sql.Stmt
	searchAuthorsStmt   *sql.Stmt
	searchLanguagesStmt *sql.Stmt
	searchBooksStmt     *sql.Stmt
	searchQuotesStmt    *sql.Stmt
}

// Connect to an sqlite Database located at `filename` This function ensures
// that the file will be created if it does not exist, create the required
// tables if it can successfully open the file
func Connect(filename string) (db *Database, err error) {
	db = new(Database)
	db.connection, err = sql.Open("sqlite3", filename)
	if err != nil {
		return
	}
	db.Init()
	db.Prepare()
	return
}

// Close the connection to the Database, to a closed Database no statements can
// be executed, meaning that every `Commit` call of any `DAO` will fail
func (db *Database) Close() {
	db.connection.Close()
}

// create tables
const (
	createBook = `CREATE TABLE Books (
Id INTEGER PRIMARY KEY AUTOINCREMENT,
AuthorId INTEGER NOT NULL,
TopicId INTEGER NOT NULL,
ISBN varchar UNIQUE,
Title varchar NOT NULL,
LanguageId INTEGER NOT NULL,
ReleaseDate date NOT NULL,
FOREIGN KEY (AuthorId) REFERENCES Authors(Id),
FOREIGN KEY (TopicId) REFERENCES Topics(Id),
FOREIGN KEY (LanguageId) REFERENCES Languages(Id)
);`
	createTopic = `CREATE TABLE Topics (
Id INTEGER PRIMARY KEY AUTOINCREMENT,
Topic varchar NOT NULL UNIQUE
);`
	createAuthor = `CREATE TABLE Authors (
Id INTEGER PRIMARY KEY AUTOINCREMENT,
Name varchar NOT NULL UNIQUE
);`
	createQuote = `CREATE TABLE Quotes (
Id INTEGER PRIMARY KEY AUTOINCREMENT,
BookId INTEGER NOT NULL,
Quote varchar NOT NULL,
Page INTEGER NOT NULL,
RecordDate date NOT NULL DEFAULT CURRENT_DATE,
FOREIGN KEY (BookId) REFERENCES Books(Id)
);`
	createLanguage = `CREATE TABLE Languages (
Id INTEGER PRIMARY KEY AUTOINCREMENT,
Language varchar NOT NULL UNIQUE
);`
)

// Initialize the Database by creating the tables required for quote.
func (db *Database) Init() (err error) {
	// create tables
	_, err = db.connection.Exec(createTopic)
	if err != nil {
		return
	}
	_, err = db.connection.Exec(createAuthor)
	if err != nil {
		return
	}
	_, err = db.connection.Exec(createLanguage)
	if err != nil {
		return
	}
	_, err = db.connection.Exec(createBook)
	if err != nil {
		return
	}
	_, err = db.connection.Exec(createQuote)
	return
}

// Prepare Statements
const (
	selectBooks = `SELECT * FROM Books
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id;`
	selectTopics  = "SELECT * FROM Topics;"
	selectAuthors = "SELECT * FROM Authors;"
	selectQuotes  = `SELECT * FROM Quotes
JOIN Books ON Quotes.BookId = Books.Id
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id;`
	selectLanguages = "SELECT * FROM Languages;"
)

const (
	selectBook = `SELECT * FROM Books
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Books.Id = ?;`
	selectTopic  = "SELECT * FROM Topics WHERE Id = ?;"
	selectAuthor = "SELECT * FROM Authors WHERE Id = ?;"
	selectQuote  = `SELECT * FROM Quotes
JOIN Books ON Quotes.BookId = Books.Id
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Quotes.Id = ?;`
	selectLanguage = "SELECT * FROM Languages WHERE Id = ?;"
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
	updateQuote    = "UPDATE Quotes SET BookId = ?, Quote = ?, Page = ? WHERE Id = ?;"
	updateLanguage = "UPDATE Languages SET Language = ? WHERE Id = ?;"
)

// related entries
const (
	relatedBooksOfAuthor = `SELECT * FROM Books
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Authors.Id = ?;`
	relatedQuotesOfAuthor = `SELECT * FROM Quotes
JOIN Books ON Quotes.BookId = Books.Id
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Authors.Id = ?;`
	relatedBooksOfLanguage = `SELECT * FROM Books
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Languages.Id = ?;`
	relatedQuotesOfLanguage = `SELECT * FROM Quotes
JOIN Books ON Quotes.BookId = Books.Id
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Languages.Id = ?;`
	relatedBooksOfTopic = `SELECT * FROM Books
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Topics.Id = ?;`
	relatedQuotesOfTopic = `SELECT * FROM Quotes
JOIN Books ON Quotes.BookId = Books.Id
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Topics.Id = ?;`
	relatedQuotesOfBook = `SELECT * FROM Quotes
JOIN Books ON Quotes.BookId = Books.Id
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Books.Id = ?;`
)

// searches
const (
	searchTopics    = `SELECT * FROM Topics WHERE Topic LIKE ?;`
	searchAuthors   = `SELECT * FROM Authors WHERE Name LIKE ?;`
	searchLanguages = `SELECT * FROM Languages WHERE Language LIKE ?;`
	searchBooks     = `SELECT * FROM Books
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Books.Title LIKE ? OR Books.ISBN LIKE ?;`
	searchQuotes = `SELECT * FROM Quotes
JOIN Books ON Quotes.BookId = Books.Id
JOIN Authors ON Books.AuthorId = Authors.Id
JOIN Topics ON Books.TopicId = Topics.Id
JOIN Languages ON Books.LanguageId = Languages.Id
WHERE Quotes.Quote LIKE ?;`
)

// Prepare the queries used for the tables created by `Init'.
func (db *Database) Prepare() (err error) {
	// select statements
	db.selectTopicsStmt, err = db.connection.Prepare(selectTopics)
	if err != nil {
		return
	}
	db.selectAuthorsStmt, err = db.connection.Prepare(selectAuthors)
	if err != nil {
		return
	}
	db.selectLanguagesStmt, err = db.connection.Prepare(selectLanguages)
	if err != nil {
		return
	}
	db.selectBooksStmt, err = db.connection.Prepare(selectBooks)
	if err != nil {
		return
	}
	db.selectQuotesStmt, err = db.connection.Prepare(selectQuotes)
	if err != nil {
		return
	}

	// select by id statements
	db.selectTopicStmt, err = db.connection.Prepare(selectTopic)
	if err != nil {
		return
	}
	db.selectAuthorStmt, err = db.connection.Prepare(selectAuthor)
	if err != nil {
		return
	}
	db.selectLanguageStmt, err = db.connection.Prepare(selectLanguage)
	if err != nil {
		return
	}
	db.selectBookStmt, err = db.connection.Prepare(selectBook)
	if err != nil {
		return
	}
	db.selectQuoteStmt, err = db.connection.Prepare(selectQuote)
	if err != nil {
		return
	}

	// insert statements
	db.insertTopicStmt, err = db.connection.Prepare(insertTopic)
	if err != nil {
		return
	}
	db.insertAuthorStmt, err = db.connection.Prepare(insertAuthor)
	if err != nil {
		return
	}
	db.insertLanguageStmt, err = db.connection.Prepare(insertLanguage)
	if err != nil {
		return
	}
	db.insertBookStmt, err = db.connection.Prepare(insertBook)
	if err != nil {
		return
	}
	db.insertQuoteStmt, err = db.connection.Prepare(insertQuote)
	if err != nil {
		return
	}

	// update statements
	db.updateTopicStmt, err = db.connection.Prepare(updateTopic)
	if err != nil {
		return
	}
	db.updateAuthorStmt, err = db.connection.Prepare(updateAuthor)
	if err != nil {
		return
	}
	db.updateLanguageStmt, err = db.connection.Prepare(updateLanguage)
	if err != nil {
		return
	}
	db.updateBookStmt, err = db.connection.Prepare(updateBook)
	if err != nil {
		return
	}
	db.updateQuoteStmt, err = db.connection.Prepare(updateQuote)
	if err != nil {
		return
	}

	// related entries statements
	db.relatedBooksOfTopicStmt, err = db.connection.Prepare(relatedBooksOfTopic)
	if err != nil {
		return
	}
	db.relatedQuotesOfTopicStmt, err = db.connection.Prepare(relatedQuotesOfTopic)
	if err != nil {
		return
	}
	db.relatedBooksOfAuthorStmt, err = db.connection.Prepare(relatedBooksOfAuthor)
	if err != nil {
		return
	}
	db.relatedQuotesOfAuthorStmt, err = db.connection.Prepare(relatedQuotesOfAuthor)
	if err != nil {
		return
	}
	db.relatedBooksOfLanguageStmt, err = db.connection.Prepare(relatedBooksOfLanguage)
	if err != nil {
		return
	}
	db.relatedQuotesOfLanguageStmt, err = db.connection.Prepare(relatedQuotesOfLanguage)
	if err != nil {
		return
	}
	db.relatedQuotesOfBookStmt, err = db.connection.Prepare(relatedQuotesOfBook)

	// searches
	db.searchTopicsStmt, err = db.connection.Prepare(searchTopics)
	if err != nil {
		return
	}
	db.searchAuthorsStmt, err = db.connection.Prepare(searchAuthors)
	if err != nil {
		return
	}
	db.searchLanguagesStmt, err = db.connection.Prepare(searchLanguages)
	if err != nil {
		return
	}
	db.searchBooksStmt, err = db.connection.Prepare(searchBooks)
	if err != nil {
		return
	}
	db.searchQuotesStmt, err = db.connection.Prepare(searchQuotes)
	if err != nil {
		return
	}
	return
}

func (db Database) GetTopic(id int) (topic Topic, err error) {
	var res *sql.Rows
	if res, err = db.selectTopicStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			topic.stmt = db.updateTopicStmt
			err = res.Scan(&topic.id, &topic.Topic)
		}
	}
	return
}

func (db Database) GetTopics() (topics []Topic, err error) {
	var res *sql.Rows
	if res, err = db.selectTopicsStmt.Query(); res != nil {
		for res.Next() && err == nil {
			topic := Topic{stmt: db.updateTopicStmt}
			err = res.Scan(&topic.id, &topic.Topic)
			topics = append(topics, topic)
		}
	}
	return
}

func (db Database) RelatedBooksOfTopic(id int) (books []Book, err error) {
	var res *sql.Rows
	if res, err = db.relatedBooksOfTopicStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			book := Book{stmt: db.updateBookStmt}
			book.Language.stmt = db.updateLanguageStmt
			book.Author.stmt = db.updateAuthorStmt
			book.Topic.stmt = db.updateTopicStmt
			err = res.Scan(&book.id,
				&book.Author.id,
				&book.Topic.id,
				&book.ISBN,
				&book.Title,
				&book.Language.id,
				&book.ReleaseDate,
				&book.Author.id,
				&book.Author.Name,
				&book.Topic.id,
				&book.Topic.Topic,
				&book.Language.id,
				&book.Language.Language)
			books = append(books, book)
		}
	}
	return
}

func (db Database) RelatedQuotesOfTopic(id int) (quotes []Quote, err error) {
	var res *sql.Rows
	if res, err = db.relatedQuotesOfTopicStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			quote := Quote{stmt: db.updateQuoteStmt}
			quote.Book.stmt = db.updateBookStmt
			quote.Book.Author.stmt = db.updateAuthorStmt
			quote.Book.Topic.stmt = db.updateTopicStmt
			quote.Book.Language.stmt = db.updateLanguageStmt
			err = res.Scan(&quote.id,
				&quote.Book.id,
				&quote.Quote,
				&quote.Page,
				&quote.RecordDate,
				&quote.Book.id,
				&quote.Book.Author.id,
				&quote.Book.Topic.id,
				&quote.Book.ISBN,
				&quote.Book.Title,
				&quote.Book.Language.id,
				&quote.Book.ReleaseDate,
				&quote.Book.Author.id,
				&quote.Book.Author.Name,
				&quote.Book.Topic.id,
				&quote.Book.Topic.Topic,
				&quote.Book.Language.id,
				&quote.Book.Language.Language)
			quotes = append(quotes, quote)
		}
	}
	return
}

func (db Database) SearchTopics(search string) (topics []Topic, err error) {
	var res *sql.Rows
	if res, err = db.searchTopicsStmt.Query("%" + search + "%"); res != nil {
		for res.Next() && err == nil {
			topic := Topic{stmt: db.updateTopicStmt}
			err = res.Scan(&topic.id, &topic.Topic)
			topics = append(topics, topic)
		}
	}
	return
}

func (db Database) GetAuthor(id int) (author Author, err error) {
	var res *sql.Rows
	if res, err = db.selectAuthorStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			author.stmt = db.updateAuthorStmt
			err = res.Scan(&author.id, &author.Name)
		}
	}
	return
}

func (db Database) GetAuthors() (authors []Author, err error) {
	var res *sql.Rows
	if res, err = db.selectAuthorsStmt.Query(); res != nil {
		for res.Next() && err == nil {
			author := Author{stmt: db.updateAuthorStmt}
			err = res.Scan(&author.id, &author.Name)
			authors = append(authors, author)
		}
	}
	return
}

func (db Database) RelatedBooksOfAuthor(id int) (books []Book, err error) {
	var res *sql.Rows
	if res, err = db.relatedBooksOfAuthorStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			book := Book{stmt: db.updateBookStmt}
			book.Language.stmt = db.updateLanguageStmt
			book.Author.stmt = db.updateAuthorStmt
			book.Topic.stmt = db.updateTopicStmt
			err = res.Scan(&book.id,
				&book.Author.id,
				&book.Topic.id,
				&book.ISBN,
				&book.Title,
				&book.Language.id,
				&book.ReleaseDate,
				&book.Author.id,
				&book.Author.Name,
				&book.Topic.id,
				&book.Topic.Topic,
				&book.Language.id,
				&book.Language.Language)
			books = append(books, book)
		}
	}
	return
}

func (db Database) RelatedQuotesOfAuthor(id int) (quotes []Quote, err error) {
	var res *sql.Rows
	if res, err = db.relatedQuotesOfAuthorStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			quote := Quote{stmt: db.updateQuoteStmt}
			quote.Book.stmt = db.updateBookStmt
			quote.Book.Author.stmt = db.updateAuthorStmt
			quote.Book.Topic.stmt = db.updateTopicStmt
			quote.Book.Language.stmt = db.updateLanguageStmt
			err = res.Scan(&quote.id,
				&quote.Book.id,
				&quote.Quote,
				&quote.Page,
				&quote.RecordDate,
				&quote.Book.id,
				&quote.Book.Author.id,
				&quote.Book.Topic.id,
				&quote.Book.ISBN,
				&quote.Book.Title,
				&quote.Book.Language.id,
				&quote.Book.ReleaseDate,
				&quote.Book.Author.id,
				&quote.Book.Author.Name,
				&quote.Book.Topic.id,
				&quote.Book.Topic.Topic,
				&quote.Book.Language.id,
				&quote.Book.Language.Language)
			quotes = append(quotes, quote)
		}
	}
	return
}

func (db Database) SearchAuthors(search string) (authors []Author, err error) {
	var res *sql.Rows
	if res, err = db.searchAuthorsStmt.Query("%" + search + "%"); res != nil {
		for res.Next() && err == nil {
			author := Author{stmt: db.updateAuthorStmt}
			err = res.Scan(&author.id, &author.Name)
			authors = append(authors, author)
		}
	}
	return
}

func (db Database) GetLanguage(id int) (language Language, err error) {
	var res *sql.Rows
	if res, err = db.selectLanguageStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			language.stmt = db.updateTopicStmt
			err = res.Scan(&language.id, &language.Language)
		}
	}
	return
}

func (db Database) GetLanguages() (languages []Language, err error) {
	var res *sql.Rows
	if res, err = db.selectLanguagesStmt.Query(); res != nil {
		for res.Next() && err == nil {
			language := Language{stmt: db.updateTopicStmt}
			err = res.Scan(&language.id, &language.Language)
			languages = append(languages, language)
		}
	}
	return
}

func (db Database) RelatedBooksOfLanguage(id int) (books []Book, err error) {
	var res *sql.Rows
	if res, err = db.relatedBooksOfLanguageStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			book := Book{stmt: db.updateBookStmt}
			book.Language.stmt = db.updateLanguageStmt
			book.Author.stmt = db.updateAuthorStmt
			book.Topic.stmt = db.updateTopicStmt
			err = res.Scan(&book.id,
				&book.Author.id,
				&book.Topic.id,
				&book.ISBN,
				&book.Title,
				&book.Language.id,
				&book.ReleaseDate,
				&book.Author.id,
				&book.Author.Name,
				&book.Topic.id,
				&book.Topic.Topic,
				&book.Language.id,
				&book.Language.Language)
			books = append(books, book)
		}
	}
	return
}

func (db Database) RelatedQuotesOfLanguage(id int) (quotes []Quote, err error) {
	var res *sql.Rows
	if res, err = db.relatedQuotesOfLanguageStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			quote := Quote{stmt: db.updateQuoteStmt}
			quote.Book.stmt = db.updateBookStmt
			quote.Book.Author.stmt = db.updateAuthorStmt
			quote.Book.Topic.stmt = db.updateTopicStmt
			quote.Book.Language.stmt = db.updateLanguageStmt
			err = res.Scan(&quote.id,
				&quote.Book.id,
				&quote.Quote,
				&quote.Page,
				&quote.RecordDate,
				&quote.Book.id,
				&quote.Book.Author.id,
				&quote.Book.Topic.id,
				&quote.Book.ISBN,
				&quote.Book.Title,
				&quote.Book.Language.id,
				&quote.Book.ReleaseDate,
				&quote.Book.Author.id,
				&quote.Book.Author.Name,
				&quote.Book.Topic.id,
				&quote.Book.Topic.Topic,
				&quote.Book.Language.id,
				&quote.Book.Language.Language)
			quotes = append(quotes, quote)
		}
	}
	return
}

func (db Database) SearchLanguages(search string) (languages []Language, err error) {
	var res *sql.Rows
	if res, err = db.searchLanguagesStmt.Query("%" + search + "%"); res != nil {
		for res.Next() && err == nil {
			language := Language{stmt: db.updateLanguageStmt}
			err = res.Scan(&language.id, &language.Language)
			languages = append(languages, language)
		}
	}
	return
}

func (db Database) GetBook(id int) (book Book, err error) {
	var res *sql.Rows
	if res, err = db.selectBookStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			book.stmt = db.updateBookStmt
			book.Language.stmt = db.updateLanguageStmt
			book.Author.stmt = db.updateAuthorStmt
			book.Topic.stmt = db.updateTopicStmt
			err = res.Scan(&book.id,
				&book.Author.id,
				&book.Topic.id,
				&book.ISBN,
				&book.Title,
				&book.Language.id,
				&book.ReleaseDate,
				&book.Author.id,
				&book.Author.Name,
				&book.Topic.id,
				&book.Topic.Topic,
				&book.Language.id,
				&book.Language.Language)
		}
	}
	return
}

func (db Database) GetBooks() (books []Book, err error) {
	var res *sql.Rows
	if res, err = db.selectBooksStmt.Query(); res != nil {
		for res.Next() && err == nil {
			book := Book{stmt: db.updateBookStmt}
			book.Language.stmt = db.updateLanguageStmt
			book.Author.stmt = db.updateAuthorStmt
			book.Topic.stmt = db.updateTopicStmt
			err = res.Scan(&book.id,
				&book.Author.id,
				&book.Topic.id,
				&book.ISBN,
				&book.Title,
				&book.Language.id,
				&book.ReleaseDate,
				&book.Author.id,
				&book.Author.Name,
				&book.Topic.id,
				&book.Topic.Topic,
				&book.Language.id,
				&book.Language.Language)
			books = append(books, book)
		}
	}
	return
}

func (db Database) RelatedQuotesOfBook(id int) (quotes []Quote, err error) {
	var res *sql.Rows
	if res, err = db.relatedQuotesOfBookStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			quote := Quote{stmt: db.updateQuoteStmt}
			quote.Book.stmt = db.updateBookStmt
			quote.Book.Author.stmt = db.updateAuthorStmt
			quote.Book.Topic.stmt = db.updateTopicStmt
			quote.Book.Language.stmt = db.updateLanguageStmt
			err = res.Scan(&quote.id,
				&quote.Book.id,
				&quote.Quote,
				&quote.Page,
				&quote.RecordDate,
				&quote.Book.id,
				&quote.Book.Author.id,
				&quote.Book.Topic.id,
				&quote.Book.ISBN,
				&quote.Book.Title,
				&quote.Book.Language.id,
				&quote.Book.ReleaseDate,
				&quote.Book.Author.id,
				&quote.Book.Author.Name,
				&quote.Book.Topic.id,
				&quote.Book.Topic.Topic,
				&quote.Book.Language.id,
				&quote.Book.Language.Language)
			quotes = append(quotes, quote)
		}
	}
	return
}

func (db Database) SearchBooks(search string) (books []Book, err error) {
	var res *sql.Rows
	if res, err = db.searchBooksStmt.Query("%"+search+"%", "%"+search+"%"); res != nil {
		for res.Next() && err == nil {
			book := Book{stmt: db.updateBookStmt}
			book.Language.stmt = db.updateLanguageStmt
			book.Author.stmt = db.updateAuthorStmt
			book.Topic.stmt = db.updateTopicStmt
			err = res.Scan(&book.id,
				&book.Author.id,
				&book.Topic.id,
				&book.ISBN,
				&book.Title,
				&book.Language.id,
				&book.ReleaseDate,
				&book.Author.id,
				&book.Author.Name,
				&book.Topic.id,
				&book.Topic.Topic,
				&book.Language.id,
				&book.Language.Language)
			books = append(books, book)
		}
	}
	return
}

func (db Database) GetQuote(id int) (quote Quote, err error) {
	var res *sql.Rows
	if res, err = db.selectQuoteStmt.Query(id); res != nil {
		for res.Next() && err == nil {
			quote.stmt = db.updateQuoteStmt
			quote.Book.stmt = db.updateBookStmt
			quote.Book.Author.stmt = db.updateAuthorStmt
			quote.Book.Topic.stmt = db.updateTopicStmt
			quote.Book.Language.stmt = db.updateLanguageStmt
			err = res.Scan(&quote.id,
				&quote.Book.id,
				&quote.Quote,
				&quote.Page,
				&quote.RecordDate,
				&quote.Book.id,
				&quote.Book.Author.id,
				&quote.Book.Topic.id,
				&quote.Book.ISBN,
				&quote.Book.Title,
				&quote.Book.Language.id,
				&quote.Book.ReleaseDate,
				&quote.Book.Author.id,
				&quote.Book.Author.Name,
				&quote.Book.Topic.id,
				&quote.Book.Topic.Topic,
				&quote.Book.Language.id,
				&quote.Book.Language.Language)
		}
	}
	return
}

func (db Database) GetQuotes() (quotes []Quote, err error) {
	var res *sql.Rows
	if res, err = db.selectQuotesStmt.Query(); res != nil {
		for res.Next() && err == nil {
			quote := Quote{stmt: db.updateQuoteStmt}
			quote.Book.stmt = db.updateBookStmt
			quote.Book.Author.stmt = db.updateAuthorStmt
			quote.Book.Topic.stmt = db.updateTopicStmt
			quote.Book.Language.stmt = db.updateLanguageStmt
			err = res.Scan(&quote.id,
				&quote.Book.id,
				&quote.Quote,
				&quote.Page,
				&quote.RecordDate,
				&quote.Book.id,
				&quote.Book.Author.id,
				&quote.Book.Topic.id,
				&quote.Book.ISBN,
				&quote.Book.Title,
				&quote.Book.Language.id,
				&quote.Book.ReleaseDate,
				&quote.Book.Author.id,
				&quote.Book.Author.Name,
				&quote.Book.Topic.id,
				&quote.Book.Topic.Topic,
				&quote.Book.Language.id,
				&quote.Book.Language.Language)
			quotes = append(quotes, quote)
		}
	}
	return
}

func (db Database) SearchQuotes(search string) (quotes []Quote, err error) {
	var res *sql.Rows
	if res, err = db.searchQuotesStmt.Query("%" + search + "%"); res != nil {
		for res.Next() && err == nil {
			quote := Quote{stmt: db.updateQuoteStmt}
			quote.Book.stmt = db.updateBookStmt
			quote.Book.Author.stmt = db.updateAuthorStmt
			quote.Book.Topic.stmt = db.updateTopicStmt
			quote.Book.Language.stmt = db.updateLanguageStmt
			err = res.Scan(&quote.id,
				&quote.Book.id,
				&quote.Quote,
				&quote.Page,
				&quote.RecordDate,
				&quote.Book.id,
				&quote.Book.Author.id,
				&quote.Book.Topic.id,
				&quote.Book.ISBN,
				&quote.Book.Title,
				&quote.Book.Language.id,
				&quote.Book.ReleaseDate,
				&quote.Book.Author.id,
				&quote.Book.Author.Name,
				&quote.Book.Topic.id,
				&quote.Book.Topic.Topic,
				&quote.Book.Language.id,
				&quote.Book.Language.Language)
			quotes = append(quotes, quote)
		}
	}
	return
}
