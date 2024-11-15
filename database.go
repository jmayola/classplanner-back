package main

import (
	"database/sql"
	"fmt"

	"os"

	_ "github.com/go-sql-driver/mysql" //driver de base de datos
	"github.com/joho/godotenv"
)

func database() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	direction := os.Getenv("DIRECTION")
	db, err := sql.Open("mysql", direction)
	if err != nil {
		panic(err.Error())
	}
	return db
}
