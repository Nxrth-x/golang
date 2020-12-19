package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// StandardMessage is the server response
type StandardMessage struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

// Books is a list of books used as data for the API
var Books []Book

var booksIndex int = 3

// Get list of all books
func getBooks(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(Books)
}

// Get a single book
func getBook(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request) // Get all parameters
	// Loop through books and find id
	for _, item := range Books {
		if item.ID == params["id"] {
			json.NewEncoder(response).Encode(item)
			return
		}
	}
	json.NewEncoder(response).Encode(&StandardMessage{
		Error:   true,
		Message: "Not found"})
}

// Create a new book
func createBook(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(request.Body).Decode(&book)
	book.ID = strconv.Itoa(booksIndex)
	booksIndex++
	Books = append(Books, book)
	json.NewEncoder(response).Encode(&StandardMessage{
		Error:   false,
		Message: "Created"})
}

// Update a book
func updateBook(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, item := range Books {
		if item.ID == params["id"] {
			var book Book
			_ = json.NewDecoder(request.Body).Decode(&book)
			Books[index] = book
			json.NewEncoder(response).Encode(&book)
			return
		}
	}
	json.NewEncoder(response).Encode(&StandardMessage{Error: true,
		Message: "Not found"})
}

// Delete a existing book
func deleteBook(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	for index, item := range Books {
		if item.ID == params["id"] {
			Books = append(Books[:index], Books[index+1:]...)
			break
		}
	}
	json.NewEncoder(response).Encode(Books)
}

func main() {
	// Initializing router
	router := mux.NewRouter()

	// Mock data - @Todo - implement Database
	Books = append(Books, Book{
		ID:    "1",
		Isbn:  "44832425",
		Title: "Wonderfull World of Nothing",
		Author: &Author{
			Firstname: "Eder",
			Lastname:  "Lima"}})
	Books = append(Books, Book{
		ID:    "2",
		Isbn:  "44112425",
		Title: "Hello, world!",
		Author: &Author{
			Firstname: "Eder",
			Lastname:  "Lima"}})

	// Route handlers
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Running server
	log.Fatal(http.ListenAndServe(":8000", router))
}
