package main

import (
	"golangfinal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("sessions/*.html")

	r.GET("/register", handlers.ShowRegisterPage)
	r.GET("/login", handlers.ShowLoginPage)
	r.GET("/main", handlers.ShowMainPage)
	r.GET("/cart", handlers.ShowCartPage)
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)
	r.POST("/main", handlers.CreateProductHandler)
	r.POST("/cart", handlers.ShowCartPage)

	http.ListenAndServe(":8080", r)
}
