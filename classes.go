package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")
	//getting data sended
	session := sessions.Default(c)
	// docente
	// docente
	// docente
	if session.Get("user_type") == "docente" {
		classes, err := db.Prepare("SELECT classes.id_class, classes.class_name, classes.class_profesor, classes.class_curso, classes.class_color, classes.class_token FROM classes WHERE classes.class_profesor=?")
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
	} else {
		// ALUMNO
		// ALUMNO
		// ALUMNO
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

}
func joinClass(c *gin.Context) {
	db := database()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
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

type UsersClass struct {
	Name     string `json:"user_name"`
	LastName string `json:"user_lastname"`
	Photo    string `json:"user_photo"`
	Type     string `json:"user_type"`
}

func getUsersFromClass(c *gin.Context) {
	db := database()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")
	session := sessions.Default(c)
	class := c.DefaultQuery("id_class", "")
	if class == "" {
		c.String(http.StatusBadRequest, "El ID de la clase es obligatorio")
		return
	}
	userType := session.Get("user_type")
	userID := session.Get("id_user")
	if userType == nil || userID == nil {
		fmt.Println("sin permisos")
		c.String(http.StatusUnauthorized, "Usuario no autenticado")
		return
	}
	var query = `
			SELECT DISTINCT users.user_name, users.user_lastname, users.user_photo, users.user_type
			FROM classes
			INNER JOIN class_users ON class_users.id_class = classes.id_class
			INNER JOIN users ON users.id_user = class_users.id_user OR users.id_user = classes.class_profesor
			WHERE classes.id_class = ?;`

	UsClassStmt, err := db.Prepare(query)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al preparar la consulta:", err)
		return
	}
	defer UsClassStmt.Close()

	rows, err := UsClassStmt.Query(class)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al ejecutar la consulta:", err)
		return
	}
	defer rows.Close()

	var Userlist []UsersClass
	for rows.Next() {
		var list UsersClass
		var photo sql.NullString
		if err := rows.Scan(&list.Name, &list.LastName, &photo, &list.Type); err != nil {
			fmt.Println("Error al escanear la fila:", err)
			continue
		}
		list.Photo = photo.String
		Userlist = append(Userlist, list)
	}
	UsClassStmt.Close()
	rows.Close()
	if len(Userlist) == 0 {
		UsClassStmt.Close()
		rows.Close()
		c.Status(http.StatusNoContent)
	} else {
		UsClassStmt.Close()
		rows.Close()
		c.JSON(http.StatusOK, Userlist)
	}
}
