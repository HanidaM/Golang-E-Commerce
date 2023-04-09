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
	r.GET("/products", handlers.GetAllProducts)
	r.GET("/products/:id", handlers.GetProductByID)
	r.GET("/main/products/search", handlers.SearchProductByTitle)
	r.GET("/products/sort", handlers.GetSortedProducts)
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)
	r.POST("/products/add", handlers.CreateProductHandler)
	r.POST("/cart", handlers.ShowCartPage)
	r.PUT("/products/:id", handlers.UpdateProductByID)
	r.DELETE("/products/:id", handlers.DeleteProductByID)

	http.ListenAndServe(":8080", r)
}
