package main

import (
	"fmt"      // formato para strings
	"net/http" // manejador de http (socket)

	// libreria requerida por go-sql-driver

	"github.com/gin-gonic/gin" // server backend y rest api para el manejo de HTTP request

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie" //administrador de sesiones

	_ "github.com/go-sql-driver/mysql" //driver de base de datos
)

// defining user attributes disambling json object
type User struct {
	Name       string `json:"user_name"`
	Password   string `json:"user_password"`
	RePassword string `json:"user_password_confirmation"`
	LastName   string `json:"user_lastname"`
	Mail       string `json:"user_mail"`
	Type       string `json:"user_type"`
	Alias      string `json:"user_alias"`
}

func setupRouter() *gin.Engine {
	//declaring database connection
	db := database()
	r := gin.Default()
	store := cookie.NewStore([]byte("aoushd1q2y387hiawru12rfsdiuhfa93htgw8rg"))
	r.Use(sessions.Sessions("users", store))

	// login
	// login
	// login
	r.GET("/user", getUser)
	r.GET("/classes", getClasses)
	r.POST("/login", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")
		var newUser User
		//getting data sended
		session := sessions.Default(c)
		userName := session.Get("username")
		fmt.Print(userName)
		if userName != nil {
			c.String(http.StatusForbidden, "Ya tienes una sesión ingresada.")
			return
		}
		if err := c.BindJSON(&newUser); err != nil {
			return
		}
		//preparing statement
		users, err := db.Prepare("SELECT user_name, user_lastname,user_mail,user_alias,user_type,id_user FROM users WHERE user_mail=? AND user_password=?")
		if err != nil {
			c.Status(505)
		}
		defer users.Close()
		//setting query output
		var user User
		var id = 0
		err = users.QueryRow(newUser.Mail, createHash(newUser.Password)).Scan(&user.Name, &user.LastName, &user.Mail, &user.Alias, &user.Type, &id)
		if err != nil {
			fmt.Println(err)
			c.String(http.StatusForbidden, "Los datos ingresados no son correctos")
		} else {
			session.Set("username", user.Name)
			session.Set("user_lastname", user.LastName)
			session.Set("user_type", user.Type)
			session.Set("user_mail", user.Mail)
			session.Set("user_alias", user.Alias)
			session.Set("id_user", id)
			session.Save()
			c.String(http.StatusAccepted, user.Name)

		}
	})

	// registerrr
	// registerrr
	// registerrr

	r.POST("/register", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")
		var newUser User
		//getting data sended
		session := sessions.Default(c)
		userName := session.Get("username")
		fmt.Print(userName)
		if userName != nil {
			c.String(http.StatusForbidden, "Ya tienes una sesión ingresada.")
			return
		}
		if err := c.BindJSON(&newUser); err != nil {
			return
		}
		fmt.Printf("datos: %s, %s, %s, %s, %s, %s,%s", newUser.Name, newUser.LastName, newUser.Password, newUser.RePassword, newUser.Mail, newUser.Type, "")
		if newUser.Password != newUser.RePassword {
			fmt.Print("las contraseñas no coinciden papi")
			c.String(http.StatusForbidden, "Las contraseñas no coinciden")
		}
		//preparing statement
		users, err := db.Prepare("INSERT INTO `users` (`id_user`, `user_name`, `user_lastname`, `user_password`, `user_mail`, `user_type`, `user_alias`) VALUES (NULL, ?, ?, ?, ?, ?, ?);")
		if err != nil {
			c.Status(505)
		}
		defer users.Close()
		//setting query output

		_, err = users.Exec(newUser.Name, newUser.LastName, createHash(newUser.Password), newUser.Mail, newUser.Type, "")
		if err != nil {
			fmt.Print(err.Error())
			c.String(http.StatusForbidden, "La cuenta ya existe o los datos ingresados no son correctos.")
		} else {
			session.Set("username", newUser.Name)
			session.Set("user_lastname", newUser.LastName)
			session.Set("user_type", newUser.Type)
			session.Set("user_mail", newUser.Mail)
			session.Set("user_alias", newUser.Alias)
			session.Save()
			c.String(http.StatusAccepted, newUser.Type)
		}
	})
	return r
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
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
