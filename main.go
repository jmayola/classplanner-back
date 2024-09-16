package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name     string `json:"user_name"`
	Password string `json:"user_password"`
	LastName string `json:"user_lastname"`
	Mail     string `json:"user_mail"`
	Type     string `json:"user_type"`
	Alias    string `json:"user_alias"`
}


func setupRouter() *gin.Engine {
	//declaring database connection
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/classplanner")
	if err != nil {
		panic(err.Error())
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// login
	// login
	// login

	r.POST("/login", func(c *gin.Context) {
		var newUser User
		//getting data sended
		if err := c.BindJSON(&newUser); err != nil {
			return
		}
		//preparing statement
		users, err := db.Prepare("SELECT user_name FROM users WHERE user_mail=? AND user_password=?")
		if err != nil {
			c.Status(505)
		}
		defer users.Close()
		//setting query output
		var outp string
		err = users.QueryRow(newUser.Mail, newUser.Password).Scan(&outp)
		if err != nil {
			c.String(http.StatusForbidden, "Los datos ingresados no son correctos")
		} else {
			c.String(http.StatusAccepted, "Ingreso Exitoso"+outp)
		}
	})

	// registerrr
	// registerrr
	// registerrr

	r.POST("/register", func(c *gin.Context) {
		var newUser User
		//getting data sended
		if err := c.BindJSON(&newUser); err != nil {
			return
		}
		//preparing statement
		users, err := db.Prepare("INSERT INTO `users` (`id_user`, `user_name`, `user_lastname`, `user_password`, `user_mail`, `user_type`, `user_alias`) VALUES (NULL, ?, ?, ?, ?, ?, ?);")
		if err != nil {
			c.Status(505)
		}
		defer users.Close()
		//setting query output

		_, err = users.Exec(newUser.Name, newUser.LastName, newUser.Password, newUser.Mail, newUser.Type, "")
		if err != nil {
			c.Status(403)
		} else {
			c.String(http.StatusAccepted, "Ingreso Exitoso ")
		}
	})
	return r
}
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Use(CORSMiddleware())	
	r.Run(":3000")
}
