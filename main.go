package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	//declaring database connection
	r := gin.Default()
	store := cookie.NewStore([]byte("aoushd1q2y387hiawru12rfsdiuhfa93htgw8rg"))
	r.Use(sessions.Sessions("users", store))

	//images & files
	r.Static("/uploads", "./uploads")

	// in this section will be the methods with the directions, followed of functions that are in the files with their respective title.
	//calendar
	r.GET("/calendar", getCalendar)
	r.GET("/calendar/:class_token", getCalendarWithToken)

	//classes
	r.GET("/classes", getClasses)
	r.POST("/classes", createClass)
	r.POST("/joinClass", joinClass)

	//tasks
	r.GET("/tasks", getTasks)
	r.POST("/tasks", createTask)

	//submissions
	r.GET("/submission", getSubmission)
	r.GET("/submissions", getSubs)
	r.POST("/submission", createSubmission)
	r.PUT("/submission/:id_submission", updateSubmission)

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
		AllowOrigins:     []string{"https://classplanner.mayola.net.ar", "https://mayola.net.ar"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		MaxAge:           12 * time.Hour,
		// Duración de la preconsulta (preflight request)
	}))
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No se han cargado las variables de entorno.")
		panic(err.Error())
	}
	state := os.Getenv("STATE")

	// Iniciar el servidor HTTPS en el puerto 3000
	if state != "development" {
		certFile := "/etc/letsencrypt/live/mayola.net.ar/fullchain.pem"
		keyFile := "/etc/letsencrypt/live/mayola.net.ar/privkey.pem"
		if err := r.RunTLS(":30000", certFile, keyFile); err != nil {
			log.Fatalf("Error al iniciar el servidor HTTPS: %s", err)
		}
	} else {
		r.Run(":3000")
	}
	db := database()

	// Asegurarse de cerrar la conexión cuando main termine
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error al cerrar la conexión:", err)
		} else {
			fmt.Println("Conexión cerrada correctamente.")
		}
	}()

}
