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

type Calification struct {
	ID       string `json:"id_task"`
	Title    string `json:"title"`
	Calif    int    `json:"calification"`
	IDuser   string `json:"id_user"`
	Name     string `json:"user_name"`
	LastName string `json:"user_lastname"`
	Photo    string `json:"user_photo"`
}

func getCalifications(c *gin.Context) {
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
	fmt.Println(userID, userType, class)
	if userType == nil || userID == nil {
		fmt.Println("sin permisos")
		c.String(http.StatusUnauthorized, "Usuario no autenticado")
		return
	}
	var query = `
			SELECT tasks.id_task, title, calification, users.id_user, users.user_name, users.user_lastname,users.user_photo
			FROM tasks
			INNER JOIN classes ON classes.id_class = tasks.id_class
			INNER JOIN submissions on submissions.id_task = tasks.id_task
			INNER JOIN users on users.id_user = submissions.id_user	
			WHERE classes.class_profesor = ? AND classes.id_class = ?`

	califStmt, err := db.Prepare(query)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al preparar la consulta:", err)
		return
	}
	defer califStmt.Close()

	rows, err := califStmt.Query(userID, class)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al ejecutar la consulta:", err)
		return
	}
	defer rows.Close()

	var CalifiList []Calification
	for rows.Next() {
		var califi Calification
		var photo sql.NullString
		if err := rows.Scan(&califi.ID, &califi.Title, &califi.Calif, &califi.IDuser, &califi.Name, &califi.LastName, &photo); err != nil {
			fmt.Println("Error al escanear la fila:", err)
			continue
		}
		califi.Photo = photo.String
		CalifiList = append(CalifiList, califi)
	}
	califStmt.Close()
	rows.Close()
	if len(CalifiList) == 0 {
		califStmt.Close()
		rows.Close()
		c.JSON(http.StatusOK, "No se encontraron entregas.")
	} else {
		califStmt.Close()
		rows.Close()
		c.JSON(http.StatusOK, CalifiList)
	}
}
