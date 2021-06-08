package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Product struct {
	ID    int
	Name  string
	Price int
}

func Server() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleHome).Methods("GET")
	router.HandleFunc("/api/products", BrowseProduct).Methods("GET")
	return router
}

func main() {
	conn := connect()
	defer conn.Close()
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
	conn := connect()
	defer conn.Close()
	rows, err := conn.Query("SELECT * FROM products")
	if err != nil {
		renderJson(writer, map[string]interface{}{
			"message": "not found",
		})
	}
	var products []*Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Print(err)
		} else {
			products = append(products, &product)
		}
	}
	renderJson(writer, products)
}
