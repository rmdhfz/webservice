package main

import (
	"database/sql"
	"fmt"
)

func connect() *sql.DB {
	user := "root"
	password := ""
	host := "127.0.0.1"
	port := "3306"
	database := "products"

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
}
