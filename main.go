package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//var db = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/login", func(c *gin.Context) {
		var newUser User
		if err := c.BindJSON(&newUser); err != nil {
			return
		}
		//data := c.BindJSON(&newUser)
		c.IndentedJSON(http.StatusAccepted, newUser)
	})

	return r
}

func main() {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/reports")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer rows.Close()
	for rows.Next() {
		var name string

		// Escanea los datos de la fila en variables
		if err := rows.Scan(&name); err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		// Imprime los valores de la fila
		fmt.Printf("Name: %s\n", name)
	}
	defer rows.Close()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":3000")
}
