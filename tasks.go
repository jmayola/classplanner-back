package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Tasks struct {
	ID          int16  `json:"id_task"`
	Clase       int16  `json:"id_class"`
	Titulo      string `json:"title"`
	Description string `json:"description"`
	Creado      string `json:"created_on"`
	Limite      string `json:"deliver_until"`
}

func createTask(c *gin.Context) {
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")
	var newClass Classes
	session := sessions.Default(c)
	//getting data sended
	if err := c.BindJSON(&newClass); err != nil {
		c.String(http.StatusForbidden, "Debe enviar datos para poder agregar una clase.")
		return
	}
	if newClass.Name == "" || newClass.Curso == "" || newClass.Color == "" || newClass.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan campos obligatorios."})
		return
	}
	//preparing statement
	ClassN, err := db.Prepare("INSERT INTO `classes` (`id_class`, `class_name`, `class_profesor`, `class_curso`, `class_color`, `class_token`) VALUES (NULL, ?, ?, ?, ?, ?)")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
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

func getTasks(c *gin.Context) {
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")
	session := sessions.Default(c)

	// Preparando la consulta SQL
	tasks, err := db.Prepare("SELECT id_task, tasks.id_class, title, description, created_on, deliver_until FROM tasks INNER JOIN classes ON classes.id_class = tasks.id_class INNER JOIN class_users ON class_users.id_class = classes.id_class INNER JOIN users ON users.id_user = class_users.id_user WHERE users.id_user = ?")
	if err != nil {
		c.Status(505)
		return
	}
	defer tasks.Close()

	var TaskList []Tasks
	rows, err := tasks.Query(session.Get("id_user"))
	if err != nil {
		rows.Close()
		fmt.Println(err.Error())
		c.String(http.StatusBadRequest, "Error al cargar los datos de clases.")
		return
	}
	defer rows.Close()

	// Escanear las filas y convertir las fechas a time.Time
	for rows.Next() {
		var Task Tasks
		var createdOn []byte
		var deliverUntil []byte

		if err := rows.Scan(&Task.ID, &Task.Clase, &Task.Titulo, &Task.Description, &createdOn, &deliverUntil); err != nil {
			fmt.Println("Error al escanear la fila:", err)
			continue
		}

		formated, _ := time.Parse("2006-01-02 15:04:05", string(createdOn))
		Task.Creado = formated.Format(time.DateTime)

		formated2, err := time.Parse("2006-01-02 15:04:05", string(deliverUntil))
		if err != nil {
			Task.Limite = "Sin Limite"
		} else {
			Task.Limite = formated2.Format(time.DateTime)
		}
		TaskList = append(TaskList, Task)
	}

	c.JSON(http.StatusOK, TaskList)
}
