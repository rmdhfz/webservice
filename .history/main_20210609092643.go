package main

import (
	"encoding/json"
	"fmt"
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
	router.HandleFunc("/api/products/{id}", UpdateProduct).Methods("PATCH")
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
			Title:  "ValidationError",
			Status: strconv.Itoa(http.StatusUnprocessableEntity),
			Detail: "Given request Body was invalid",
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

	conn := connect()
	defer conn.Close()

	result, err := conn.Exec("DELETE FROM products WHERE id = ?", ProductID)
	if err != nil {
		log.Print(err)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		return
	}

	if affected == 0 {
		w.WriteHeader(http.StatusNotFound)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Not Found",
			Status: strconv.Itoa(http.StatusNotFound),
			Detail: fmt.Sprintf("Product with id %s not found", ProductID),
		}})
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateProduct(w http.ResponseWriter, req *http.Request) {
	productID := mux.Vars(req)["id"]
	var product Product
	err := jsonapi.UnmarshalPayload(req.Body, &product)
	if err != nil {
		w.Header().Set("Content-Type", jsonapi.MediaType)
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "ValidationError",
			Detail: "Given request is invalid",
			Status: strconv.Itoa(http.StatusUnprocessableEntity),
		}})
		return
	}
	conn := connect()
	defer conn.Close()
	var exist bool
	row := conn.Exec("SELECT id FROM products WHERE id = ?", productID)
	if condition {

	}
	query, err := conn.Prepare("UPDATE products SET name = ?, price = ? WHERE id = ?")
	if err != nil {
		log.Print(err)
		return
	}
	query.Exec(product.Name, product.Price, productID)
	product.ID, _ = strconv.ParseInt(productID, 10, 64)
	renderJson(w, &product)
}
