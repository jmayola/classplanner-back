package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	Username     string `json:"user_name"`
	Lastname     string `json:"user_lastname"`
	Alias        string `json:"user_alias"`
	Photo        string `json:"user_photo"`
}

func createSubmission(c *gin.Context) {
	db := database()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")
	var submission Submission
	session := sessions.Default(c)
	//getting data sended
	submission.ID_task = c.DefaultPostForm("id_task", "")
	submission.Comment = c.DefaultPostForm("submission_comment", "")

	//preparing statement

	SubM, err := db.Prepare("INSERT INTO `submissions` (`id_submission`, `id_user`, `id_task`, `submission_file`, `submission_comment`, `submission_date`) VALUES (NULL, ?, ?, ?, ?, ?);")
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
	_, err = SubM.Exec(session.Get("id_user"), submission.ID_task, filename, submission.Comment, time.Now())
	if err != nil {
		fmt.Print(err.Error())
		SubM.Close()
		c.String(http.StatusForbidden, "No se puedo enviar la tarea")
	} else {
		SubM.Close()
		c.String(http.StatusAccepted, "Tarea enviada.")
	}
}
func getSubmission(c *gin.Context) {
	// Conectar a la base de datos
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
	taskID := c.DefaultQuery("id_task", "")
	if taskID == "" {
		c.String(http.StatusBadRequest, "El ID de la tarea es obligatorio")
		return
	}
	// Preparar la consulta SQL para obtener las entregas del usuario
	query := `
		SELECT submission_file, submission_comment, submission_date, calification, feedback
		FROM submissions
		WHERE id_user = ? AND id_task = ?`

	// Preparar el statement SQL
	submissionsStmt, err := db.Prepare(query)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al preparar la consulta:", err)
		return
	}
	defer submissionsStmt.Close()

	var submissions Submission
	var calification sql.NullString
	var feedback sql.NullString
	err = submissionsStmt.QueryRow(session.Get("id_user"), taskID).Scan(&submissions.File, &submissions.Comment, &submissions.Date, &calification, &feedback)
	if err != nil {
		fmt.Println("Error en la consulta", err.Error())
		c.Status(http.StatusInternalServerError)
	}
	// Ejecutar la consulta
	if calification.Valid {
		submissions.Calification = calification.String
	} else {
		submissions.Calification = ""
	}
	if feedback.Valid {
		submissions.Feedback = feedback.String
	} else {
		submissions.Feedback = ""
	}
	if err != nil {
		submissionsStmt.Close()
		c.Status(http.StatusNoContent)
		return
	} else {
		submissionsStmt.Close()
		c.JSON(http.StatusOK, submissions)
	}
}
func getSubs(c *gin.Context) {
	db := database()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")
	taskID := c.DefaultQuery("id_task", "")
	if taskID == "" {
		c.String(http.StatusBadRequest, "El ID de la tarea es obligatorio")
		return
	}
	// Preparar la consulta SQL para obtener las entregas del usuario
	query := `
		SELECT id_submission, submissions.id_user, id_task, submission_file, submission_comment, submission_date, calification, feedback,
		us.user_name, us.user_lastname, us.user_alias, us.user_photo
		FROM submissions INNER JOIN users us ON us.id_user=submissions.id_user WHERE id_task = ?`

	// Preparar el statement SQL
	submissionsStmt, err := db.Prepare(query)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al preparar la consulta:", err)
		return
	}
	defer submissionsStmt.Close()

	// Ejecutar la consulta
	rows, err := submissionsStmt.Query(taskID)
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

		if err := rows.Scan(&submission.ID, &submission.ID_user, &submission.ID_task, &submission.File, &submission.Comment, &submissionDate, &submission.Calification, &submission.Feedback, &submission.Username, &submission.Lastname, &submission.Alias, &submission.Photo); err != nil {
			fmt.Println("Error al escanear la fila:", err)
			continue
		}

		// Convertir la fecha de []byte a string
		submission.Date = string(submissionDate)

		// Agregar la entrega a la lista de entregas
		submissions = append(submissions, submission)
	}

	if len(submissions) == 0 {
		submissionsStmt.Close()
		rows.Close()
		c.JSON(http.StatusOK, "No se encontraron entregas.")
	} else {
		submissionsStmt.Close()
		rows.Close()
		c.JSON(http.StatusOK, submissions)
	}
}
func updateSubmission(c *gin.Context) {
	db := database()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	ORIGIN := os.Getenv("ORIGIN")
	c.Header("Access-Control-Allow-Origin", ORIGIN)
	c.Header("Access-Control-Allow-Credentials", "true")

	submissionID := c.Param("id_submission")
	if submissionID == "" {
		c.String(http.StatusBadRequest, "El ID del envío es obligatorio")
		return
	}

	var updatedSubmission struct {
		Comment      string `json:"submission_comment"`
		File         string `json:"submission_file"`
		Calification string `json:"calification"`
		Feedback     string `json:"feedback"`
	}

	if err := c.BindJSON(&updatedSubmission); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println("Error al analizar el JSON:", err)
		return
	}

	query := `
		UPDATE submissions
		SET submission_comment = ?, submission_file = ?, calification = ?, feedback = ?
		WHERE id_submission = ?`

	updateStmt, err := db.Prepare(query)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al preparar la consulta:", err)
		return
	}
	defer updateStmt.Close()

	_, err = updateStmt.Exec(updatedSubmission.Comment, updatedSubmission.File, updatedSubmission.Calification, updatedSubmission.Feedback, submissionID)
	if err != nil {
		updateStmt.Close()
		c.Status(http.StatusInternalServerError)
		fmt.Println("Error al ejecutar la actualización:", err)
		return
	}

	updateStmt.Close()
	c.JSON(http.StatusOK, gin.H{"message": "Envío actualizado correctamente"})
}
