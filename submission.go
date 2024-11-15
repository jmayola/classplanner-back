package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Submission struct {
	ID           string `json:"id_submission" form:"id_submission"`
	ID_user      string `json:"id_user" form:"id_user"`
	ID_task      string `json:"id_task" form:"id_task"`
	File         string `json:"submission_file" form:"submission_file"`
	Comment      string `json:"submission_comment" form:"submission_comment"`
	Date         string `json:"submission_date" form:"submission_date"`
	Calification string `json:"calification" form:"calification"`
	Feedback     string `json:"feedback" form:"feedback"`
}

func createSubmission(c *gin.Context) {
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")
	var submission Submission
	session := sessions.Default(c)
	//getting data sended
	submission.ID_task = c.DefaultPostForm("id_task", "")
	submission.Comment = c.DefaultPostForm("submission_comment", "")

	//preparing statement

	SubM, err := db.Prepare("INSERT INTO `submissions` (`id_submission`, `id_user`, `id_task`, `submission_file`, `submission_comment`, `submission_date`, `calification`, `feedback`) VALUES (NULL, ?, ?, ?, ?, NULL, NULL, NULL);")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer SubM.Close()
	file, err := c.FormFile("submission_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe enviar una imagen"})
		return
	}

	uploadDir := "./uploads"
	filename := filepath.Join(uploadDir, generateFileName()+file.Filename)
	SaveFile(file, uploadDir, c, filename)

	//setting query output
	_, err = SubM.Exec(session.Get("id_user"), submission.ID_task, filename, submission.Comment)
	if err != nil {
		fmt.Print(err.Error())
		c.String(http.StatusForbidden, "No se puedo enviar la tarea")
	} else {
		c.String(http.StatusAccepted, "Tarea enviada.")
	}
}
func getSubmissions(c *gin.Context) {
	// Conectar a la base de datos
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")
	session := sessions.Default(c)

	// Preparar la consulta SQL para obtener las entregas del usuario
	query := `
		SELECT id_submission, id_user, id_task, submission_file, submission_comment, submission_date, calification, feedback
		FROM submissions
		WHERE id_user = ?`

	// Preparar el statement SQL
	submissionsStmt, err := db.Prepare(query)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al preparar la consulta:", err)
		return
	}
	defer submissionsStmt.Close()

	// Ejecutar la consulta
	rows, err := submissionsStmt.Query(session.Get("id_user"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al ejecutar la consulta:", err)
		return
	}
	defer rows.Close()

	// Crear una lista para almacenar las entregas
	var submissions []Submission

	// Iterar sobre las filas y escanear los resultados en la estructura Submission
	for rows.Next() {
		var submission Submission
		var submissionDate []byte // Para almacenar la fecha como []byte antes de convertirla a string

		if err := rows.Scan(&submission.ID, &submission.ID_user, &submission.ID_task, &submission.File, &submission.Comment, &submissionDate, &submission.Calification, &submission.Feedback); err != nil {
			fmt.Println("Error al escanear la fila:", err)
			continue
		}

		// Convertir la fecha de []byte a string
		submission.Date = string(submissionDate)

		// Agregar la entrega a la lista de entregas
		submissions = append(submissions, submission)
	}

	// Devolver las entregas como respuesta JSON
	if len(submissions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No se encontraron entregas."})
	} else {
		c.JSON(http.StatusOK, submissions)
	}
}
