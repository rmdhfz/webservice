package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var mysqlDB *sql.DB

func Server() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleHome).Methods("GET")
	router.HandleFunc("/api/products", BrowseProduct).Methods("GET")
	return router
}

func main() {
	mysqlDB := connect()
	defer mysqlDB.Close()
	router := Server()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func renderJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func handleHome(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	response := map[string]interface{}{
		"status":  200,
		"message": "Hai",
	}
	json.NewEncoder(writer).Encode(response)
}

func BrowseProduct(writer http.ResponseWriter, request *http.Request) {
	rows, err := mysqlDB.Query("SELECT * FROM products")
	renderJson(writer, map[string]interface{}{
		"message": "products",
	})
}
