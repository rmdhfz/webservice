package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Home(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	Server().ServeHTTP(res, req)
	expectedResponse := `{"message":"Hai","status":200}`
	// resBody, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	t.Log(err)
	// }
	assert.Equal(t, 200, res.Code, "Invalid Response code")
	assert.Equal(t, expectedResponse, res.Body.String())
	t.Log(res.Body.String())
}
