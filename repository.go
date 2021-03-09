package main

import (
	"errors"

	"time"
)

type Book struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	ISBNno     string
	Author     string
	Series     string
	Genre      string
	Rating     float32
	Launchdate time.Time `json:"LaunchDate,omitempty"`
}

type BookLibrary interface {
	NewBook(book Book) error
	Book(Id int) (Book, error)
	Books() []Book
	Update(book Book) error
	Delete(Id int) error
}

func NewBookLibraryServices(repo BookRepo, noti UpdateNotifier) *BookLibraryServices {
	return &BookLibraryServices{repo: repo, notifier: noti}
}

type BookLibraryServices struct {
	repo     BookRepo
	notifier UpdateNotifier
}

func (lib *BookLibraryServices) NewBook(book Book) error {

	if book.Id == 0 {
		return errors.New("Invalid Book ID")
	}

	if err := lib.repo.CreateBook(book); err != nil {
		return err
	}

	lib.notifier.Notify()
	return nil

}

func (lib *BookLibraryServices) Book(Id int) (Book, error) {

	return lib.repo.BooksByID(Id)
}

func (lib *BookLibraryServices) Books() []Book {

	return lib.repo.LaunchedBooks()
}

func (lib *BookLibraryServices) Update(book Book) error {
	return lib.repo.Update(book)
}

func (lib *BookLibraryServices) Delete(Id int) error {
	return lib.repo.Delete(Id)
}

type BookRepo interface {
	CreateBook(book Book) error
	BooksByID(Id int) (Book, error)
	LaunchedBooks() []Book
	Update(book Book) error
	Delete(Id int) error
}

type BookRepoMemory struct {
	books map[int]Book
}

func NewBookRepo() *BookRepoMemory {
	return &BookRepoMemory{books: make(map[int]Book)}
}

func (repo *BookRepoMemory) CreateBook(book Book) error {
	if _, exists := repo.books[book.Id]; exists {

		return errors.New("Book Already Present in Data!!")
	}

	repo.books[book.Id] = book
	return nil
}

func (repo *BookRepoMemory) BooksByID(Id int) (Book, error) {
	book, exists := repo.books[Id]
	if !exists {
		return Book{}, errors.New(" Book Not Present in Data!!")
	}

	return book, nil
}

func (repo *BookRepoMemory) LaunchedBooks() []Book {
	present_time := time.Now()
	books := make([]Book, 0)

	for _, book := range repo.books { // range returns index,copy_ele

		if book.Launchdate.After(present_time) {
			continue
		}

		books = append(books, book)
	}

	return books
}

func (repo *BookRepoMemory) Update(book Book) error {
	if _, exists := repo.books[book.Id]; !exists {
		return errors.New(" Book NOT Present in Data!!")
	}

	repo.books[book.Id] = book
	return nil
}

func (repo *BookRepoMemory) Delete(Id int) error {
	if Id == 0 {

		return errors.New("Invalid ID!!!")
	}

	delete(repo.books, Id)
	return nil
}

func (repo *BookRepoMemory) Notify() {
	return
}
