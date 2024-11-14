package main

import (
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DownloadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No se ha enviado ningún archivo."})
		return
	}

	// Crear un directorio donde se guardarán los archivos
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	// Guardar el archivo en el servidor
	filename := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(500, gin.H{"error": "No se pudo guardar el archivo."})
		return
	} else {
		c.String(http.StatusAccepted, "Archivo guardado")
	}
}
func SaveFile(file *multipart.FileHeader, uploadDir string, c *gin.Context, filename string) {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear directorio de carga"})
		return
	}
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el archivo"})
		return
	}
}
