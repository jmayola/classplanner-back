package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //driver de base de datos
)

func database() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/classplanner")
	if err != nil {
		panic(err.Error())
	}
	return db
}
