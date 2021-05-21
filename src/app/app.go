package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/DmitriyZhevnov/library/src/entities"
	repository "github.com/DmitriyZhevnov/library/src/repository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gopkg.in/go-playground/validator.v9"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable ",
		DbHost, DbPort, DbUser, DbName, DbPassword)
	var err error
	a.DB, err = sql.Open(Dbdriver, DBURL)
	fmt.Println(a.DB)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/books", a.FindAll).Methods("GET")
	a.Router.HandleFunc("/api/books/{id}", a.FindById).Methods("GET")
	a.Router.HandleFunc("/api/books/name/{name}", a.FindByName).Methods("GET")
	a.Router.HandleFunc("/api/books/genre/{id:[0-9]+}", a.FilterByGenre).Methods("GET")
	a.Router.HandleFunc("/api/books/price/{minPrice:[0-9]+}/{maxPrice:[0-9]+}", a.FilterByPrices).Methods("GET")
	a.Router.HandleFunc("/api/books", a.Create).Methods("POST")
	a.Router.HandleFunc("/api/books/{id:[0-9]+}", a.Update).Methods("PUT")
	a.Router.HandleFunc("/api/books/{id:[0-9]+}", a.Delete).Methods("DELETE")
}

func (a *App) FindAll(w http.ResponseWriter, r *http.Request) {
	books, err := repository.FindAll(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, books)
}

func (a *App) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid book Id")
		return
	}
	book := []entities.Book{}
	if book, err = repository.FindById(a.DB, id); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Book not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if book == nil {
		respondWithError(w, http.StatusNotFound, "Book not found")
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) FindByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	book := []entities.Book{}
	var err error
	if book, err = repository.FindByName(a.DB, name); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Book not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if book == nil {
		var array [5]string
		respondWithJSON(w, http.StatusOK, array)
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) FilterByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid book Id")
		return
	}
	book := []entities.Book{}
	if book, err = repository.FilterByGenre(a.DB, id); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Books not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if book == nil {
		var array [5]string
		respondWithJSON(w, http.StatusOK, array)
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) FilterByPrices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	minPrice, err := strconv.ParseFloat(vars["minPrice"], 64)
	maxPrice, err := strconv.ParseFloat(vars["maxPrice"], 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid prices")
		return
	}
	book := []entities.Book{}
	if book, err = repository.FilterByPrices(a.DB, minPrice, maxPrice); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Books not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if book == nil {
		var array [5]string
		respondWithJSON(w, http.StatusOK, array)
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	var book entities.Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	err = validateStruct(book)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	if err := repository.Create(a.DB, &book); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusNoContent, book.Id)
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid book Id")
		return
	}
	var book entities.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	err = validateStruct(book)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	book.Id = id
	if _, err := repository.Update(a.DB, int64(id), &book); err != nil {
		respondWithError(w, http.StatusNotModified, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, book)
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Book Id")
		return
	}
	if _, err := repository.Delete(a.DB, int64(id)); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func validateStruct(b entities.Book) error {
	validate := validator.New()
	err := validate.Struct(b)
	if err != nil {
		return err
	}
	return nil
}
