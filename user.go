package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// defining user attributes disambling json object
type User struct {
	Name       string `json:"user_name"`
	Password   string `json:"user_password"`
	RePassword string `json:"user_password_confirmation"`
	LastName   string `json:"user_lastname"`
	Mail       string `json:"user_mail"`
	Type       string `json:"user_type"`
	Alias      string `json:"user_alias"`
	Photo      string `json:"user_photo"`
}

func getUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")

	session := sessions.Default(c)
	user_name := session.Get("username")
	user_type := session.Get("user_type")
	user_lastname := session.Get("user_lastname")
	user_mail := session.Get("user_mail")
	user_alias := session.Get("user_alias")
	user_photo := session.Get("user_photo")
	if user_name == nil || user_type == nil {
		c.String(http.StatusForbidden, "No tienes una Sesión iniciada.")
	} else {
		c.JSON(http.StatusAccepted, gin.H{"user_name": user_name, "user_type": user_type, "user_lastname": user_lastname, "user_mail": user_mail, "user_alias": user_alias, "user_photo": user_photo})
	}
}
func deleteUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")

	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.String(http.StatusForbidden, "No tienes una Sesión iniciada.")
	} else {
		session.Clear()
		session.Save()
		c.String(http.StatusAccepted, "Has Salido de la sesión.")
	}
}
func register(c *gin.Context) {
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")
	var newUser User
	//getting data sended
	session := sessions.Default(c)
	userName := session.Get("username")
	fmt.Print(userName)
	if userName != nil {
		c.String(http.StatusForbidden, "Ya tienes una sesión ingresada.")
		return
	}
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	fmt.Printf("datos: %s, %s, %s, %s, %s, %s,%s", newUser.Name, newUser.LastName, newUser.Password, newUser.RePassword, newUser.Mail, newUser.Type, "")
	if newUser.Password != newUser.RePassword {
		fmt.Print("las contraseñas no coinciden papi")
		c.String(http.StatusForbidden, "Las contraseñas no coinciden")
	}
	//preparing statement
	users, err := db.Prepare("INSERT INTO `users` (`id_user`, `user_name`, `user_lastname`, `user_password`, `user_mail`, `user_type`, `user_alias`) VALUES (NULL, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		c.Status(505)
	}
	defer users.Close()
	//setting query output

	_, err = users.Exec(newUser.Name, newUser.LastName, createHash(newUser.Password), newUser.Mail, newUser.Type, "")
	if err != nil {
		fmt.Print(err.Error())
		c.String(http.StatusForbidden, "La cuenta ya existe o los datos ingresados no son correctos.")
	} else {
		c.String(http.StatusAccepted, "La cuenta ha sido registrada")
	}
}
func login(c *gin.Context) {
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")
	var newUser User
	//getting data sended
	session := sessions.Default(c)
	userName := session.Get("username")
	if userName != nil {
		c.String(http.StatusForbidden, "Ya tienes una sesión ingresada.")
		return
	}
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	//preparing statement
	users, err := db.Prepare("SELECT user_name, user_lastname,user_mail,user_alias,user_type,id_user,user_photo FROM users WHERE user_mail=? AND user_password=?")
	if err != nil {
		c.Status(505)
	}
	defer users.Close()
	//setting query output
	var user User
	var id = 0
	err = users.QueryRow(newUser.Mail, createHash(newUser.Password)).Scan(&user.Name, &user.LastName, &user.Mail, &user.Alias, &user.Type, &id, &user.Photo)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusForbidden, "Los datos ingresados no son correctos")
	} else {
		session.Set("username", user.Name)
		session.Set("user_lastname", user.LastName)
		session.Set("user_type", user.Type)
		session.Set("user_mail", user.Mail)
		session.Set("user_alias", user.Alias)
		session.Set("id_user", id)
		session.Set("user_photo", user.Photo)
		session.Save()
		c.JSON(http.StatusAccepted, gin.H{"user_name": user.Name, "user_type": user.Type})

	}
}
func updateUser(c *gin.Context) {
	db := database()
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")

	var upUser User
	session := sessions.Default(c)
	id := session.Get("id_user")

	if err := c.ShouldBind(&upUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al parsear datos"})
		return
	}

	user, err := db.Prepare("UPDATE `users` SET `user_photo` = ? WHERE `users`.`id_user` = ?;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al preparar consulta"})
		return
	}
	defer user.Close()

	file, err := c.FormFile("user_photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe enviar una imagen"})
		return
	}

	uploadDir := "./uploads"
	filename := filepath.Join(uploadDir, generateFileName()+file.Filename)
	SaveFile(file, uploadDir, c, filename)

	_, err = user.Exec(filename, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la base de datos"})
		return
	}
	session.Set("user_photo", filename)
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Imagen de perfil actualizada", "file": file.Filename})
}
