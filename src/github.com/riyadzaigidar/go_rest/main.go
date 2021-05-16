package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// default responses
type defaultResponses struct{
	Msg string `json:"msg"`
}

// struct books
type Book struct{
	Id string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Author *Author `json:"author"`
}

// type book validator
type BookValidator struct{
	Id *string `json:"id"`
	Isbn *string `json:"isbn"`
	Title *string `json:"title"`
	Author *AuthorValidator `json:"author"`
}

// type author validator

type AuthorValidator struct{
	Firstname *string `json:"firstname"`
	Lastname *string `json:"lastname"`
}
// struct author
type Author struct{
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

// init books
var books []Book

func updateItem(id string, data BookValidator) Book{
	// doesnt change item in array
	for _, item := range(books){
		if item.Id == id{
			if data.Isbn != nil {
				item.Isbn = *data.Isbn
			}
			if data.Title != nil {
				item.Title = *data.Title
			}
			// if data.Author.Firstname != nil{
			// 	fmt.Println(data.Author.Firstname)
			// }
			// if data.Author.Lastname != nil{
			// 	fmt.Println(data.Author.Lastname)
			// }
			fmt.Println(item)
			fmt.Println(books)
			return item
		}
	}
	return Book{}
}

// delete item from a slice
func deleteItem(arr []Book, id string) []Book{
	for index, item := range(arr){
		if (item.Id == id){
			books = append(arr[:index], arr[index+1:]...)
			return books
		}
	}
	return arr
}

// Get list of books
func getBooks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get single book detail
func getBook(w http.ResponseWriter, r *http.Request){
	param := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	for _, item := range(books){
		if item.Id == param["id"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(defaultResponses{
		Msg: "item not found",
	})
}

// create book
func addBook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	book := Book{}
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.Id = strconv.Itoa(rand.Intn(1000000))
	books = append(books, book)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&book)
}

// update book
func upDateBook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	book := BookValidator{}
	_ = json.NewDecoder(r.Body).Decode(&book)
	updBook := updateItem(param["id"], book)
	fmt.Println(books)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(updBook)
}

// delete book
func deleteBook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	deleteItem(books, param["id"])
	w.WriteHeader(http.StatusNoContent)

}

func main(){
	// Initialize the router
	router := mux.NewRouter()
	// initialize mock data
	// default author
	author := Author{
		Firstname: "Sherlock",
		Lastname: "Holmes",
	}

	books = append(books, 
		Book{
			Id: "1",
			Isbn: "435",
			Title: "Harry Potter",
			Author: &author,
		},
		Book{
			Id: "2",
			Isbn: "799",
			Title: "Oh my god",
			Author: &author,
		},
		Book{
			Id: "3",
			Isbn: "123123",
			Title: "Last Day",
			Author: &author,
		},
	)	

	// Declare the routes
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books/", addBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", upDateBook).Methods("PATCH")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// listen on 7000 port with this routes
	// add localhost to the string so it runs only on localhost(uses regex)
	http.ListenAndServe(":7000", router)
}