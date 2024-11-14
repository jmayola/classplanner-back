package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie" //administrador de sesiones
	"github.com/gin-gonic/gin"               // server backend y rest api para el manejo de HTTP request
	_ "github.com/go-sql-driver/mysql"       //driver de base de datos
)

func setupRouter() *gin.Engine {
	//declaring database connection
	r := gin.Default()
	store := cookie.NewStore([]byte("aoushd1q2y387hiawru12rfsdiuhfa93htgw8rg"))
	r.Use(sessions.Sessions("users", store))

	//images & files
	r.Static("/uploads", "./uploads")

	// in this section will be the methods with the directions, followed of functions that are in the files with their respective title.

	//classes
	r.GET("/classes", getClasses)
	r.POST("/classes", createClass)

	//user
	r.GET("/user", getUser)
	r.DELETE("/user", deleteUser)
	r.PUT("/user", updateUser)
	r.POST("/login", login)
	r.POST("/register", register)

	//returing methods
	return r
}
func main() {
	r := setupRouter()
	r.Use(CORSMiddleware())
	r.Run(":3000")
}
