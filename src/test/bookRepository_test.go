package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/DmitriyZhevnov/library/src/app"
)

var a app.App

func TestMain(m *testing.M) {
	a = app.App{}
	a.Initialize("root", "950621", "library")
	code := m.Run()
	os.Exit(code)
}

func TestFindAll(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":1,"name":"book1","price":10,"genre":1,"amount":50},{"id":2,"name":"book2","price":11,"genre":2,"amount":1},{"id":3,"name":"book3","price":20.6,"genre":3,"amount":3},{"id":4,"name":"book4","price":25,"genre":1,"amount":4},{"id":5,"name":"book5","price":30.5,"genre":2,"amount":2}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFindById(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":1,"name":"book1","price":10,"genre":1,"amount":50}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFindNonExistentId(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/books/1000", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Book not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Book not found'. Got '%s'", m["error"])
	}
}

func TestFindByName(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books/name/book1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":1,"name":"book1","price":10,"genre":1,"amount":50}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFindByGenre(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books/genre/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":1,"name":"book1","price":10,"genre":1,"amount":50},{"id":4,"name":"book4","price":25,"genre":1,"amount":4}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFilterByPrices(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/books/price/10/20", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[{"id":1,"name":"book1","price":10,"genre":1,"amount":50},{"id":2,"name":"book2","price":11,"genre":2,"amount":1}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUpdate(t *testing.T) {
	clearTable()
	addUser()
	req, _ := http.NewRequest("GET", "/api/books/1", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	var originalBook map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalBook)

	payload := []byte(`{
		"name": "Some new name",
		"price": 15,
		"genre": 2,
		"amount": 10 }`)
	req, err := http.NewRequest("PUT", "/api/books/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] == originalBook["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalBook["name"], m["name"], m["name"])
	}

	if m["genre"] == originalBook["genre"] {
		t.Errorf("Expected the genre to change from '%v' to '%v'. Got '%v'", originalBook["genre"], m["genre"], m["genre"])
	}
	if m["amount"] == originalBook["amount"] {
		t.Errorf("Expected the amount to change from '%v' to '%v'. Got '%v'", originalBook["amount"], m["amount"], m["amount"])
	}
}

func TestCreate(t *testing.T) {
	jsonStr := []byte(`{
			"name": "The Three Musketeers",
			"price": 10.44,
			"genre": 1,
			"amount": 5 }`)
	req, _ := http.NewRequest("POST", "/api/books", bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)

	checkResponseCode(t, http.StatusNoContent, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	var i = response.Body
	_, err := strconv.ParseInt(i.String(), 10, 64)
	if err != nil {
		t.Errorf("Erorr! %v", err)
	}
}

func TestDelete(t *testing.T) {
	clearTable()
	addUser()
	req, _ := http.NewRequest("GET", "/api/books/1", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/books/1", nil)
	response = httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/api/books/1", nil)
	response = httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addUser() {
	a.DB.Exec(`INSERT INTO library.book (name, price, genre_id, amount) VALUES ('test_book', '25', '1', '30');`)
}

func clearTable() {
	a.DB.Exec("DELETE FROM library.book")
	a.DB.Exec("ALTER TABLE library.book AUTO_INCREMENT = 1")
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
