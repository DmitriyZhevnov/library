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

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
