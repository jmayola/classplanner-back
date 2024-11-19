package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
	r.POST("/joinClass", joinClass)

	//tasks
	r.GET("/tasks", getTasks)
	r.POST("/tasks", createTask)

	//submissions
	r.GET("/submission", getSubmission)
	r.POST("/submission", createSubmission)

	//comments
	r.GET("/comments", getComments)
	r.POST("/comments", createComment)

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
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	r.Use(CORSMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5173"}, // CORRECCIÓN: incluye el puerto correctamente con ":"
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},        // Métodos permitidos
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},        // Cabeceras permitidas
		AllowCredentials: true,                                                       // Permite enviar credenciales (cookies, autenticación)
		ExposeHeaders:    []string{"Content-Length", "Authorization"},                // Cabeceras expuestas
		MaxAge:           12 * time.Hour,
		// Duración de la preconsulta (preflight request)
	}))
	// if err := r.RunTLS(":30000", "fullchain.pem", "privkey.pem"); err != nil {
	r.Run(":3000")
	// panic(err.Error())
	// }

}
