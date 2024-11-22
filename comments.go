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

type Comment struct {
	ID         int64          `json:"id"`
	Task       int64          `json:"id_task"`
	Text       string         `json:"text"`
	UserName   string         `json:"userName"`
	Time       string         `json:"time"`
	User_photo sql.NullString `json:"user_photo"`
}

func createComment(c *gin.Context) {
	db := database()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")
	var newComment Comment
	session := sessions.Default(c)

	if err := c.BindJSON(&newComment); err != nil {
		c.String(http.StatusBadRequest, "El comentario no puede estar vacío")
		return
	}

	Comment, err := db.Prepare("INSERT INTO `comments` (`id_comment`, `id_user`, `id_task`, `comment`) VALUES (NULL, ?, ?, ?);")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer Comment.Close()
	//setting query output
	_, err = Comment.Exec(session.Get("id_user"), newComment.Task, newComment.Text)
	if err != nil {
		fmt.Print(err.Error())
		c.String(http.StatusForbidden, "La cuenta ya existe o los datos ingresados no son correctos.")
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Comentario enviado con éxito", "comment": newComment})
	}
	// Responder con el comentario creado
}
func getComments(c *gin.Context) {
	// Establecer la conexión a la base de datos
	db := database()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")
	// Obtener el ID de la tarea desde los parámetros de la URL
	taskID := c.DefaultQuery("id_task", "")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID de la tarea es obligatorio"})
		return
	}
	// Preparar la consulta SQL para obtener los comentarios de la tarea
	query := `
		SELECT id_comment, id_task, comment, created_on, users.user_name, users.user_photo 
		FROM comments 
		INNER JOIN users ON comments.id_user = users.id_user 
		WHERE id_task = ?
		ORDER BY created_on DESC
	`

	// Ejecutar la consulta
	rows, err := db.Query(query, taskID)
	if err != nil {
		fmt.Println("Error al obtener comentarios:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Crear una lista de comentarios
	var comments []Comment

	// Iterar sobre los resultados de la consulta
	for rows.Next() {
		var comment Comment
		var createdOn []byte
		// Escanear cada fila de resultados
		if err := rows.Scan(&comment.ID, &comment.Task, &comment.Text, &createdOn, &comment.UserName, &comment.User_photo); err != nil {
			fmt.Println("Error al escanear fila:", err)
			continue
		}

		// Convertir la fecha desde formato []byte a un string
		comment.Time = string(createdOn)

		// Agregar el comentario a la lista de comentarios
		comments = append(comments, comment)
	}

	// Verificar si no se encontraron comentarios
	if len(comments) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No se encontraron comentarios para esta tarea"})
		return
	}

	// Devolver los comentarios en la respuesta
	c.JSON(http.StatusOK, comments)
}
