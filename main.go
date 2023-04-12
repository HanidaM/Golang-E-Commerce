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
	r.GET("/products/search", handlers.SearchProductByTitle)
	r.GET("/products/sort", handlers.GetSortedProducts)
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)
	r.POST("/products/add", handlers.CreateProductHandler)
	r.POST("/cart", handlers.ShowCartPage)
	r.PUT("/products/:id", handlers.UpdateProductByID)
	r.DELETE("/products/:id", handlers.DeleteProductByID)

	http.ListenAndServe(":8080", r)
}

// {
//     "name": "Example Product",
//     "description": "stuff",
//     "price": 20.99,
//     "rating": 3.0,
//     "image": "https://images-na.ssl-images-amazon.com/images/I/81fyoFoaxlL._AC_UL127_SR127,127_.jpg",
//     "quantity": 100
// }