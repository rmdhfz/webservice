package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Home(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	Server().ServeHTTP(res, req)
	expectedResponse = `{"message":"Hai","status":200}`
	resBody, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, 300, res.Code, "Invalid Response code")
	t.Log(res.Body.String())
}
