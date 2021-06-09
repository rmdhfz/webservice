package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
)

type Product struct {
	ID    int64  `jsonapi:"primary,products"`
	Name  string `jsonapi:"attr,name"`
	Price int    `jsonapi:"attr,price"`
}

func Server() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleHome).Methods("GET")
	router.HandleFunc("/api/products", BrowseProduct).Methods("GET")
	router.HandleFunc("/api/products", CreateProduct).Methods("POST")
	router.HandleFunc("/api/products/{id}", DeleteProduct).Methods("DELETE")
	return router
}

func main() {
	conn := connect()
	defer conn.Close()
	router := Server()
	log.Fatal(http.ListenAndServe(":8080", router))
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

func CreateProduct(w http.ResponseWriter, req *http.Request) {
	var product Product
	err := jsonapi.UnmarshalPayload(req.Body, &product)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title: "ValidationError",
			Status: strconv.Itoa(http.StatusUnprocessableEntity),
			Detail: "Given request Body was invalid"
		}})
		return
	}
	conn := connect()
	defer conn.Close()
	query, err := conn.Prepare("INSERT INTO products (name, price) VALUES (?, ?)")
	if err != nil {
		log.Print(err)
		return
	}
	result, err := query.Exec(product.Name, product.Price)
	if err != nil {
		log.Print(err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Print(err)
		return
	}
	product.ID = lastId
	renderJson(w, &product)
}

func DeleteProduct(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ProductID := mux.Vars(req)["id"]
	
	conn := connect();
	defer conn.Close()

	result, err := conn.Exec("DELETE FROM products WHERE id = ?", ProductID)
	if err != nil {
		log.Print(err)
		return
	}
}