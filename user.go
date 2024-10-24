package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Credentials", "true")

	session := sessions.Default(c)
	username := session.Get("username")
	user_type := session.Get("user_type")
	if username == nil || user_type == nil {
		fmt.Print(username)
		fmt.Print(user_type)
		c.String(http.StatusForbidden, "No tienes una Sesión iniciada.")
	} else {
		c.JSON(200, gin.H{"username": username, "user_type": user_type})
	}
}
