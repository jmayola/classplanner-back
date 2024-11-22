package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Tasks struct {
	ID          int16  `json:"id_task"`
	Clase       int    `json:"id_class"`
	Titulo      string `json:"title"`
	Description string `json:"description"`
	Creado      string `json:"created_on"`
	Limite      string `json:"deliver_until"`
}

func createTask(c *gin.Context) {
	db := database()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")
	var newTask Tasks
	//getting data sended
	if err := c.BindJSON(&newTask); err != nil {
		fmt.Println(err)
		c.String(http.StatusForbidden, "Debe enviar datos para poder agregar una tarea.")
		return
	}
	if newTask.Clase == 0 || newTask.Titulo == "" || newTask.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan campos obligatorios."})
		return
	}
	//preparing statement
	ClassN, err := db.Prepare("INSERT INTO `tasks` (`id_task`, `id_class`, `title`, `description`, `created_on`, `deliver_until`) VALUES (NULL, ?, ?, ?, ?, ?);")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer ClassN.Close()
	//setting query output
	if newTask.Limite == "" {
		_, err = ClassN.Exec(newTask.Clase, newTask.Titulo, newTask.Description, time.Now(), nil)
		if err != nil {
			fmt.Print(err.Error())
			ClassN.Close()
			c.String(http.StatusForbidden, "La Tarea ya existe o los datos ingresados no son correctos.")
		} else {
			ClassN.Close()
			c.String(http.StatusAccepted, "Tarea creada.")
		}
	}
	_, err = ClassN.Exec(newTask.Clase, newTask.Titulo, newTask.Description, time.Now(), newTask.Limite)
	if err != nil {
		fmt.Print(err.Error())
		ClassN.Close()
		c.String(http.StatusForbidden, "La Tarea ya existe o los datos ingresados no son correctos.")
	} else {
		ClassN.Close()
		c.String(http.StatusAccepted, "Tarea creada.")
	}
}

func getTasks(c *gin.Context) {
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

	userType := session.Get("user_type")
	userID := session.Get("id_user")

	if userType == nil || userID == nil {
		c.String(http.StatusUnauthorized, "Usuario no autenticado")
		return
	}
	//had to define query
	var query string
	if userType == "alumno" {
		query = `
			SELECT id_task, tasks.id_class, title, description, created_on, deliver_until
			FROM tasks
			INNER JOIN classes ON classes.id_class = tasks.id_class
			INNER JOIN class_users ON class_users.id_class = classes.id_class
			INNER JOIN users ON users.id_user = class_users.id_user
			WHERE users.id_user = ?`
	} else if userType == "docente" {
		query = `
			SELECT id_task, tasks.id_class, title, description, created_on, deliver_until
			FROM tasks
			INNER JOIN classes ON classes.id_class = tasks.id_class
			WHERE classes.class_profesor = ?`
	} else {
		c.String(http.StatusForbidden, "Tipo de usuario no válido")
		return
	}

	tasksStmt, err := db.Prepare(query)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al preparar la consulta:", err)
		return
	}
	defer tasksStmt.Close()

	rows, err := tasksStmt.Query(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al ejecutar la consulta:", err)
		return
	}
	defer rows.Close()

	var TaskList []Tasks
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
	tasksStmt.Close()
	rows.Close()
	c.JSON(http.StatusOK, TaskList)
}
