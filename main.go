package main

import (
	"log"
	"melivecode/go-jwt-api/controller/auth"
	"melivecode/go-jwt-api/controller/user"
	"melivecode/go-jwt-api/middleware"
	"melivecode/go-jwt-api/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
func init() {

    err := godotenv.Load(".env")

    if err != nil {
        log.Fatal("Error loading .env file")
    }
}
func main() {
	orm.InitDB()
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)
	authorize:=r.Group("/users",middleware.JWTAuthen())
	authorize.GET("/readall", user.ReadAll)
	authorize.GET("/profile", user.Profile)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
