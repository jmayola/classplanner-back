package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func createClass(c *gin.Context) {
	db := database()
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
		session.Set("user_type", newUser.Type)
		c.String(http.StatusAccepted, newUser.Type)
	}
}
