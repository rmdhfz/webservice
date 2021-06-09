package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_Home(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	Server().ServeHTTP(res, req)
	expectedResponse := `{"message":"Hai","status":200}`
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Log(err)
	}
	assert.Equal(t, 200, res.Code, "Invalid Response code")
	assert.Equal(t, expectedResponse, string(bytes.TrimSpace(resBody)))
	t.Log(res.Body.String())
}

func Test_BrowseProduct(t *testing.T) {
	_, mock, err := sqlmock.New()
	if err != nil {
		t.Log(err)
	}
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Bumbu Racik 2", 100)
	mock.ExpectQuery("SELECT * FROM products").WillReturnRows(rows)
	req, _ := http.NewRequest("GET", "/api/products", nil)
	res := httptest.NewRecorder()

	Server().ServeHTTP(res, req)
	expectedResponse := `[{"ID":1,"Name":"Bumbu Racik 2","Price":100}]`
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Log(err)
	}
	assert.Equal(t, 200, "Invalid response code")
	assert.Equal(t, expectedResponse, string(bytes.TrimSpace(resBody)))

}
