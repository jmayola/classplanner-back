package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Classes struct {
	ID       int16  `json:"id_class"`
	Name     string `json:"class_name"`
	Profesor int16  `json:"class_profesor"`
	Curso    string `json:"class_curso"`
	Color    string `json:"class_color"`
	Token    string `json:"class_token"`
}

func createClass(c *gin.Context) {
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")
	var newUser User
	//getting data sended
	session := sessions.Default(c)
	// userName := session.Get("username")
	// user_type := session.Get("user_type")

	//preparing statement
	//INSERT INTO `classes` (`id_class`, `class_name`, `class_profesor`, `class_curso`, `class_color`, `class_token`) VALUES (NULL, 'Artes Visuales', '28', '7mo 2da', '#992233', 'as123sda');
	users, err := db.Prepare("INSERT INTO `classes` (`id_user`, `user_name`, `user_lastname`, `user_password`, `user_mail`, `user_type`, `user_alias`) VALUES (NULL, ?, ?, ?, ?, ?, ?);")
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
func getClasses(c *gin.Context) {
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")
	//getting data sended
	session := sessions.Default(c)
	// userName := session.Get("username")
	// user_type := session.Get("user_type")
	//
	//preparing statement
	classes, err := db.Prepare("SELECT users.id_user, classes.class_name, classes.class_curso FROM class_users INNER JOIN users ON users.id_user=class_users.id_user INNER JOIN classes ON classes.id_class=class_users.id_class WHERE users.id_user=?")
	if err != nil {
		c.Status(505)
	}
	defer classes.Close()
	//setting query output
	var class Classes
	err = classes.QueryRow(session.Get("id_user")).Scan(&class)
	if err != nil {
		fmt.Print(err.Error())
		c.String(http.StatusForbidden, "La cuenta ya existe o los datos ingresados no son correctos.")
	} else {
		c.JSON(202, class)
	}
}
