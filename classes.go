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
	var newClass Classes
	session := sessions.Default(c)
	//getting data sended
	if err := c.BindJSON(&newClass); err != nil {
		c.String(http.StatusForbidden, "Debe enviar datos para poder agregar una clase.")
	}
	if newClass.Name == "" || newClass.Curso == "" || newClass.Color == "" || newClass.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan campos obligatorios."})
		return
	}
	//preparing statement
	ClassN, err := db.Prepare("INSERT INTO `classes` (`id_class`, `class_name`, `class_profesor`, `class_curso`, `class_color`, `class_token`) VALUES (NULL, ?, ?, ?, ?, ?)")
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	defer ClassN.Close()
	//setting query output

	_, err = ClassN.Exec(newClass.Name, session.Get("id_user"), newClass.Curso, newClass.Color, newClass.Token)
	if err != nil {
		fmt.Print(err.Error())
		c.String(http.StatusForbidden, "La cuenta ya existe o los datos ingresados no son correctos.")
	} else {
		c.String(http.StatusAccepted, "Clase creada.")
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
	classes, err := db.Prepare("SELECT classes.id_class, classes.class_name, classes.class_profesor, classes.class_curso, classes.class_color, classes.class_token FROM class_users INNER JOIN users ON users.id_user=class_users.id_user INNER JOIN classes ON classes.id_class=class_users.id_class WHERE users.id_user=?")
	if err != nil {
		c.Status(505)
	}
	defer classes.Close()
	//setting query output
	var classList []Classes
	rows, err := classes.Query(session.Get("id_user"))

	if err != nil {
		rows.Close()
		fmt.Print(err.Error())
		c.String(http.StatusBadRequest, "Error al cargar los datos de clases.")
	} else {
		for rows.Next() {
			var class Classes
			if err := rows.Scan(&class.ID, &class.Name, &class.Profesor, &class.Curso, &class.Color, &class.Token); err != nil {
				fmt.Println("Error al escanear la fila:", err)
				continue // if there is a error, we just keep with the next one
			}
			classList = append(classList, class)
		}
		c.JSON(202, classList)
	}
}
func joinClass(c *gin.Context) {
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")
	var class_token Classes
	session := sessions.Default(c)
	//getting data sended

	if err := c.BindJSON(&class_token); err != nil {
		c.String(http.StatusForbidden, "Debe enviar datos para poder ser añadido a una clase.")
	}
	if class_token.Token == "" {
		c.String(http.StatusBadRequest, "No puedes ingresar campos nulos.")
		return
	}
	//getting class_id from token
	classId, err := db.Prepare("SELECT classes.id_class FROM `classes` WHERE classes.class_token=?")
	if err != nil {
		fmt.Println(err.Error())
		c.String(404, err.Error())
	}
	defer classId.Close()
	//getting id class
	var class_id Classes
	err = classId.QueryRow(class_token.Token).Scan(&class_id.ID)
	if err != nil {
		fmt.Println(err.Error())
		c.String(404, "No existe esa clase.")
	}
	defer classId.Close()

	//preparing statement
	NewMember, err := db.Prepare("INSERT INTO `class_users` (`id`, `id_user`, `id_class`) VALUES (NULL, ?, ?);")
	if err != nil {
		c.String(http.StatusForbidden, "Ya has sido añadido a la clase.")
	}
	defer NewMember.Close()
	//setting query output
	_, err = NewMember.Exec(session.Get("id_user"), class_id.ID)
	if err != nil {
		fmt.Print(err.Error())
		c.String(http.StatusForbidden, "La cuenta ingresada o la clase no exiten.")
	} else {
		c.String(http.StatusAccepted, "Has sido añadido a la clase.")
	}
}
