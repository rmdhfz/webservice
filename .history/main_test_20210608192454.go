package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
