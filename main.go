package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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

func NewHTTPTransport(l BookLibrary) *httptransport {
	return &httptransport{service: l}
}

type httptransport struct {
	service  BookLibrary
	upgrader websocket.Upgrader
}

func (transport *httptransport) handleWebSocket(pool *Pool, w http.ResponseWriter, r *http.Request) {

	conn, err := transport.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading web socket connection", err)
		return
	}

	defer conn.Close()

	client := &Client{
		ID:   conn.RemoteAddr().String(),
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	go client.Read()

	tick := time.NewTicker(time.Second * 5)

	defer conn.Close()

	for range tick.C {

		books := transport.service.Books()
		if err := conn.WriteJSON(books); err != nil {
			log.Println("error while writing to json:", err)
			return
		}

	}
}

func notfoundhandler(w http.ResponseWriter, r *http.Request) {
	log.Println("not found", r.RequestURI)

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
	pool := NewPool()
	go pool.Start()

	// r.HandleFunc("/", Homepage).Methods("GET")
	r.HandleFunc("/books", t.AddBook).Methods("POST")
	r.HandleFunc("/books", t.Bookshow).Methods("GET")
	r.HandleFunc("/books/{id}", t.BookID).Methods("GET")
	r.HandleFunc("/books/{id}", t.UpdateBook).Methods("PATCH")
	r.HandleFunc("/books/{id}", t.DeleteBook).Methods("DELETE")
	r.NotFoundHandler = http.HandlerFunc(notfoundhandler)

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { t.handleWebSocket(pool, w, r) }).Methods("GET")

	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("static"))))

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
		noOp = noOpNotifier{}
		repo = NewBookRepo()
		//pgrepo         = NewPostgresBook(db)
		libraryservice = NewBookLibraryServices(repo, noOp)
		tr             = NewHTTPTransport(libraryservice)
	)

	//defer db.Close()

	fmt.Println("Connected!")

	er := libraryservice.NewBook(Book{
		Id:     1,
		Title:  "Moon",
		ISBNno: "CH1",
		Author: "WD",

		Series:     "CC",
		Launchdate: time.Now(),
		Genre:      "Fiction",
		Rating:     5.0,
	})

	if er != nil {
		fmt.Println("Error while dealing with library service: ", er)
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
