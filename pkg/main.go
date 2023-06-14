package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

const (
	API_PATH = "/apis/v1/books"
)

type Book struct {
	Id, Name, Isbn string
}

type library struct {
	dbHost, dbPass, dbUser, dbName string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "password"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbUser: dbUser,
		dbName: dbName,
	}

	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet)
	r.HandleFunc(apiPath, l.postBook).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)
}

func (l library) postBook(w http.ResponseWriter, r *http.Request) {
	//read the request into an instance of Book
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)

	// open connection
	db := l.openConnection()

	// insert the books into db
	insertQuery, err := db.Prepare("insert into books values (?, ? , ?)")
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	// close connection
	l.closeConnection(db)
}

func (l library) getBooks(w http.ResponseWriter, r *http.Request) {
	// open connection
	db := l.openConnection()

	// read all the books
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("querying the books table %s\n", err.Error())
	}

	books := []Book{}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("while scanning the row %s\n", err.Error())
		}
		aBook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, aBook)
	}

	json.NewEncoder(w).Encode(books)

	// close connection
	l.closeConnection(db)
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", l.dbUser, l.dbPass, l.dbHost, l.dbName))

	if err != nil {
		log.Fatalf("Error while opening the connection to the database %s\n", err.Error())
	}

	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()

	if err != nil {
		panic(err)
	}
}
