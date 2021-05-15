package book_api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DmitriyZhevnov/library/src/config"
	"github.com/DmitriyZhevnov/library/src/entities"
	"github.com/DmitriyZhevnov/library/src/repository"
	"github.com/gorilla/mux"
)

func FindAll(response http.ResponseWriter, request *http.Request) {
	db, err := config.GetDB()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		bookModel := repository.BookRepository{
			Db: db,
		}
		books, err2 := bookModel.FindAll()
		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err.Error())
		} else {
			respondWithJson(response, http.StatusOK, books)
		}
	}
}

func FindAllByName(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	db, err := config.GetDB()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		bookRepository := repository.BookRepository{
			Db: db,
		}
		books, err2 := bookRepository.FindAllByName(name)
		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err.Error())
		} else {
			respondWithJson(response, http.StatusOK, books)
		}
	}
}

func FilterByPrices(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	minInRequest := vars["minPrice"]
	maxInRequest := vars["maxPrice"]
	min, _ := strconv.ParseFloat(minInRequest, 64)
	max, _ := strconv.ParseFloat(maxInRequest, 64)
	db, err := config.GetDB()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		bookRepository := repository.BookRepository{
			Db: db,
		}
		books, err2 := bookRepository.FilterByPrices(min, max)
		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err.Error())
		} else {
			respondWithJson(response, http.StatusOK, books)
		}
	}
}

func FindById(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idInRequest := vars["id"]
	id, err := strconv.ParseInt(idInRequest, 10, 64)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	}
	db, err := config.GetDB()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		bookRepository := repository.BookRepository{
			Db: db,
		}
		books, err2 := bookRepository.FindById(int(id))
		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err.Error())
		} else {
			respondWithJson(response, http.StatusOK, books)
		}
	}
}

func FilterByGenre(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idInRequest := vars["id"]
	id, err := strconv.ParseInt(idInRequest, 10, 64)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	}
	db, err := config.GetDB()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		bookRepository := repository.BookRepository{
			Db: db,
		}
		books, err2 := bookRepository.FilterByGenre(int(id))
		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err.Error())
		} else {
			respondWithJson(response, http.StatusOK, books)
		}
	}
}

func Create(response http.ResponseWriter, request *http.Request) {
	var book entities.Book
	err := json.NewDecoder(request.Body).Decode(&book)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	}
	db, err := config.GetDB()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		bookRepository := repository.BookRepository{
			Db: db,
		}
		err2 := bookRepository.Create(&book)
		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err2.Error())
		} else {
			respondWithJson(response, http.StatusOK, book.Id)
		}
	}
}

func Update(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idInRequest := vars["id"]
	id, err := strconv.ParseInt(idInRequest, 10, 64)
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	}
	var book entities.Book
	err = json.NewDecoder(request.Body).Decode(&book)
	db, err := config.GetDB()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		bookRepository := repository.BookRepository{
			Db: db,
		}
		_, err2 := bookRepository.Update(id, &book)
		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err.Error())
		} else {
			book.Id = int(id)
			respondWithJson(response, http.StatusOK, book)
		}
	}
}

func Delete(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idInRequest := vars["id"]
	id, _ := strconv.ParseInt(idInRequest, 10, 64)
	db, err := config.GetDB()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		bookRepository := repository.BookRepository{
			Db: db,
		}
		rows, err2 := bookRepository.Delete(id)
		if err2 != nil {
			respondWithError(response, http.StatusBadRequest, err.Error())
		} else {
			respondWithJson(response, http.StatusNoContent, map[string]int64{
				"Rows Affected: ": rows,
			})
		}
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
