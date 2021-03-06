package main

import (
	"encoding/json"
	"fmt"
	_ "fmt"

	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is my homepage")

}

func (transport *httptransport) AddBook(w http.ResponseWriter, r *http.Request) {

	var b Book

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		fmt.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := transport.service.NewBook(b)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Book Added")

}

func (transport *httptransport) Bookshow(w http.ResponseWriter, r *http.Request) {

	res := transport.service.Books()
	js, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error Occured: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func (transport *httptransport) BookID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := transport.service.Book(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	js, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func (transport *httptransport) UpdateBook(w http.ResponseWriter, r *http.Request) {

	var b Book

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b.Id = id

	json.NewDecoder(r.Body).Decode(&b)
	fmt.Printf("%+v\n", b)
	er := transport.service.Update(b)

	if er != nil {
		fmt.Println(err)
		return
	}

}

func (transport *httptransport) DeleteBook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	er := transport.service.Delete(id)
	if er != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println("Book with ID", id, "Deleted")

}

// { "id" : "32" , "title": "Treasure Island", "ISBNno":  "4AB234", "Author":  "DK", "Series": "", "Launchdate": "01/02/2021", "Genre": "Fiction", "Rating": "4.8"}

func NewHTTPTransport(l BookLibrary) *httptransport {
	return &httptransport{service: l}
}

type httptransport struct {
	service BookLibrary
}

func notfoundhandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func setupRoutes(t *httptransport) {

	r := mux.NewRouter()

	r.HandleFunc("/", Homepage).Methods("GET")
	r.HandleFunc("/books", t.AddBook).Methods("POST")
	r.HandleFunc("/books", t.Bookshow).Methods("GET")
	r.HandleFunc("/books/{id}", t.BookID).Methods("GET")
	r.HandleFunc("/books/{id}", t.UpdateBook).Methods("PATCH")
	r.HandleFunc("/books/{id}", t.DeleteBook).Methods("DELETE")
	r.NotFoundHandler = http.HandlerFunc(notfoundhandler)
	http.Handle("/", r)
	r.Use(loggingMiddleware)
}

func main() {

	//psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// db, err := sql.Open("postgres", psqlconn)
	// if err != nil {
	// 	fmt.Println("Error while connecting to Database: ", err)
	// }

	var (
		repo = NewBookRepo()
		//pgrepo         = NewPostgresBook(db)
		libraryservice = NewBookLibraryServices(repo)
		tr             = NewHTTPTransport(libraryservice)
	)

	//defer db.Close()

	fmt.Println("Connected!")

	er := libraryservice.NewBook(Book{
		Id:     14,
		Title:  "Sun",
		ISBNno: "A3",
		Author: "lop",

		Series:     "CC",
		Launchdate: time.Now(),
		Genre:      "Fiction",
		Rating:     5.0,
	})

	if er != nil {
		fmt.Println("Error while dealing with librayr service: ", er)
	}
	books := libraryservice.Books()
	fmt.Println(books)
	fmt.Println(repo)

	//Server Side
	fmt.Println(("Server up and Running!!"))
	setupRoutes(tr)
	log.Fatal(http.ListenAndServe(":8088", nil))
	//Database

}
