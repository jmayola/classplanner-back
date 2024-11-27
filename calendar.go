package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Calendar struct {
	Title     string `json:"title"`
	Desc      string `json:"description"`
	IDtask    int    `json:"id_task"`
	Created   string `json:"created_on"`
	Deliver   string `json:"deliver_until"`
	ClassName string `json:"class_name"`
	Curso     string `json:"class_curso"`
}

func getCalendar(c *gin.Context) {
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

	//getting classnames and tasks titles from user
	var query string
	if session.Get("user_type") == "docente" {
		query = `
		SELECT tasks.title, tasks.description, tasks.id_task, tasks.created_on, tasks.deliver_until, classes.class_name, classes.class_curso FROM tasks INNER JOIN classes ON classes.id_class = tasks.id_class WHERE classes.class_profesor = ?;
		`
	} else {
		query = `
		SELECT tasks.title, tasks.description, tasks.id_task, tasks.created_on, tasks.deliver_until, classes.class_name, classes.class_curso
		FROM tasks
		INNER JOIN classes ON classes.id_class = tasks.id_class
		INNER JOIN class_users ON class_users.id_class = classes.id_class
		INNER JOIN users ON users.id_user = class_users.id_user
		WHERE users.id_user = ?;
		`
	}
	calendar, err := db.Prepare(query)
	if err != nil {
		c.Status(505)
	}
	defer calendar.Close()
	//setting query output
	var calendarList []Calendar
	rows, err := calendar.Query(session.Get("id_user"))

	if err != nil {
		rows.Close()
		fmt.Print(err.Error())
		c.String(http.StatusBadRequest, "Error al cargar los datos de calendario.")
	} else {
		for rows.Next() {
			var calend Calendar
			var deliver_until []byte
			var created_on []byte

			if err := rows.Scan(&calend.Title, &calend.Desc, &calend.IDtask, &created_on, &deliver_until, &calend.ClassName, &calend.Curso); err != nil {
				fmt.Println("Error al escanear la fila:", err)
				continue // if there is a error, we just keep with the next one
			}
			calend.Deliver = string(deliver_until)
			calend.Created = string(created_on)
			calendarList = append(calendarList, calend)
		}
		c.JSON(202, calendarList)
	}
}
func getCalendarWithToken(c *gin.Context) {
	db := database()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")
	classID := c.DefaultQuery("class_token", "")
	if classID == "" {
		c.String(http.StatusBadRequest, "El ID de la clase es obligatorio")
		return
	}
	//getting data sended
	session := sessions.Default(c)

	//getting classnames and tasks titles from user

	var query string
	if session.Get("user_type") == "docente" {
		query = `
		SELECT tasks.title, tasks.description, tasks.id_task, tasks.created_on, tasks.deliver_until, classes.class_name, classes.class_curso FROM tasks INNER JOIN classes ON classes.id_class = tasks.id_class WHERE classes.class_profesor = ? AND classes.class_token = ?;
		`
	} else {
		query = `
		SELECT tasks.title, tasks.description, tasks.id_task, tasks.created_on, tasks.deliver_until, classes.class_name, classes.class_curso
		FROM tasks
		INNER JOIN classes ON classes.id_class = tasks.id_class
		INNER JOIN class_users ON class_users.id_class = classes.id_class
		INNER JOIN users ON users.id_user = class_users.id_user
		WHERE users.id_user = ? AND classes.class_token = ?;
		`
	}
	calendar, err := db.Prepare(query)
	if err != nil {
		c.Status(505)
	}
	fmt.Println(classID)
	defer calendar.Close()
	//setting query output
	var calendarList []Calendar
	rows, err := calendar.Query(session.Get("id_user"), classID)

	if err != nil {
		rows.Close()
		fmt.Print(err.Error())
		c.String(http.StatusBadRequest, "Error al cargar los datos de calendario.")
	} else {
		for rows.Next() {
			var calend Calendar
			var deliver_until []byte
			var created_on []byte

			if err := rows.Scan(&calend.Title, &calend.Desc, &calend.IDtask, &created_on, &deliver_until, &calend.ClassName, calend.Curso); err != nil {
				fmt.Println("Error al escanear la fila:", err)
				continue // if there is a error, we just keep with the next one
			}
			calend.Deliver = string(deliver_until)
			calend.Created = string(created_on)
			calendarList = append(calendarList, calend)
		}
		c.JSON(202, calendarList)
	}
}
