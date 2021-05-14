package main

import (
	"fmt"
	"net/http"

	bookapi "github.com/DmitriyZhevnov/library/src/apis/book_api"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/books/", bookapi.FindAll).Methods("GET")
	router.HandleFunc("/api/books/{id}", bookapi.FindById).Methods("GET")
	router.HandleFunc("/api/books/name/{name}", bookapi.FindAllByName).Methods("GET")
	router.HandleFunc("/api/books/genre/{id}", bookapi.FilterByGenre).Methods("GET")
	router.HandleFunc("/api/books/price/{minPrice}/{maxPrice}", bookapi.FilterByPrices).Methods("GET")
	router.HandleFunc("/api/books/", bookapi.Create).Methods("POST")
	router.HandleFunc("/api/books/{id}", bookapi.Update).Methods("PUT")
	router.HandleFunc("/api/books/{id}", bookapi.Delete).Methods("DELETE")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
