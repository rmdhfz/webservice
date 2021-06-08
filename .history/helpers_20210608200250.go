package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func connect() *sql.DB {
	user := "root"
	password := ""
	host := "127.0.0.1"
	port := "3306"
	database := "products"

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func renderJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
