package main

import (
	"auth/infrastructure/auth"
	"auth/infrastructure/persistence"
	"auth/interfaces"
	"auth/interfaces/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {

	login := os.Getenv("LOGIN_FE_URL")

	dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	//redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	services, err := persistence.NewRepositories(dbdriver, user, password, port, host, dbname)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.Automigrate()

	redisService, err := auth.NewRedisDB(redis_host, redis_port, redis_password)
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()

	users := interfaces.NewUsers(services.User)
	authenticate := interfaces.NewAuthenticate(services.User, redisService.Auth, tk)
	extAuth := interfaces.NewExtAuth(services.User, redisService.Auth, tk)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware()) //For CORS

	//Home
	r.GET("/", func(c *gin.Context) {
		http.Redirect(c.Writer, c.Request, login, http.StatusTemporaryRedirect)
	})

	//user routes
	r.POST("/users", users.SaveUser)
	r.GET("/users", users.GetUsers)
	r.GET("/users/:user_id", users.GetUser)

	//authentication routes
	r.POST("/login", authenticate.Login)
	r.POST("/logout", authenticate.Logout)
	r.POST("/refresh", authenticate.Refresh)

	//External auth
	r.GET("/auth/facebook", extAuth.HandleFacebookLogin)
	r.GET("/auth/facebook/callback", extAuth.HandleFacebookCallback)
	r.GET("/auth/google", extAuth.HandleGoogleLogin)
	r.GET("/auth/google/callback", extAuth.HandleGoogleCallback)
	r.POST("/auth/verification", extAuth.HanldeExternalClientId)

	//Starting the application
	app_port := os.Getenv("PORT") //using heroku host
	if app_port == "" {
		app_port = "9191" //localhost
	}
	log.Fatal(r.Run(":" + app_port))
}
