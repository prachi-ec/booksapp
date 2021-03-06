package main

import (
	"fmt"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "bookdb"
)

type PostgresBook struct {
	db sql.DB
}

func NewPostgresBook(db *sql.DB) *PostgresBook {
	return &PostgresBook{db: *db}
}

func (pgbook *PostgresBook) CreateBook(book Book) error {

	sqlStatement := `INSERT INTO booktable (id, title, isbnno, author,series,genre,rating,launchdate) VALUES ($1, $2, $3, $4, $5,$6,$7,$8)`
	_, err := pgbook.db.Exec(sqlStatement, book.Id, book.Title, book.ISBNno, book.Author, book.Series, book.Genre, book.Rating, book.Launchdate)
	if err != nil {
		fmt.Println("Error while executing query:-  ", err)
		return err
	}
	fmt.Println("Book Added Successfully!!!")
	return nil
}

func (pgbook *PostgresBook) LaunchedBooks() []Book {

	books := make([]Book, 0)

	rows, err := pgbook.db.Query("SELECT *FROM booktable")

	if err != nil {
		fmt.Println("Error while Quering:--  ", err)
		return nil
	}

	for rows.Next() {
		var b Book

		if err := rows.Scan(&b.Id, &b.Title, &b.ISBNno, &b.Author, &b.Series, &b.Genre, &b.Rating, &b.Launchdate); err != nil {
			fmt.Println("Error while Scanning:- ", err)
			return books
		}

		books = append(books, b)
	}

	return books
}

func (pgbook *PostgresBook) Update(book Book) error {

	sqlStatement := `UPDATE booktable SET title=$2, isbnno=$3, author=$4,series=$5,genre=$6,rating=$7, launchdate=$8 WHERE ID=$1`
	_, err := pgbook.db.Exec(sqlStatement, book.Id, book.Title, book.ISBNno, book.Author, book.Series, book.Genre, book.Rating, book.Launchdate)
	if err != nil {
		fmt.Println("Error while updating: ", err)
		return err

	}
	return nil

}

func (pgbook *PostgresBook) BooksByID(Id int) (Book, error) {

	var title string
	var isbn, author, series, genre sql.NullString
	var id int
	var rating sql.NullInt32
	var launchdate sql.NullTime
	if err := pgbook.db.QueryRow("SELECT * FROM BOOKTABLE WHERE id = $1", Id).Scan(&id, &title, &isbn, &author, &series, &genre, &rating, &launchdate); err != nil {
		fmt.Println("Error while Scanning:- ", err)
	}
	b := Book{
		Id:         id,
		Title:      title,
		ISBNno:     "",
		Author:     "",
		Series:     "",
		Launchdate: time.Now(),
		Genre:      "",
		Rating:     2,
	}
	return b, nil

}

func (pgbook *PostgresBook) Delete(Id int) error {

	if err := pgbook.db.QueryRow("DELETE  FROM BOOKTABLE WHERE id = $1", Id).Scan(); err != nil {
		fmt.Println("Error while Scanning:- ", err)
	}

	return nil
}

// func main() {

// 	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

// 	// open database
// 	db, err := sql.Open("postgres", psqlconn)
// 	if err != nil {
// 		fmt.Println("Error while connecting:- ", err)
// 	}

// 	defer db.Close()

// 	fmt.Println("Connected!")

// 	var pgrepo = NewPostgresBook(db)

// 	// b := Book{
// 	// 	Id:         8,
// 	// 	Title:      "TimeTravel",
// 	// 	ISBNno:     "833H2",
// 	// 	Author:     "Frankk",
// 	// 	Series:     "None",
// 	// 	Launchdate: time.Now(),
// 	// 	Genre:      "Fiction",
// 	// 	Rating:     2.0,
// 	// }

// 	//driver

// 	//er := pgrepo.CreateBook(b)
// 	//er := pgrepo.Books()
// 	//book, er := pgrepo.BooksByID(8)
// 	er := pgrepo.Delete(8)
// 	if er != nil {
// 		fmt.Println("Error:--- ", err)
// 	}
// 	//fmt.Println(book)
// }

/////////////////////////////////////////////////////////////////////
