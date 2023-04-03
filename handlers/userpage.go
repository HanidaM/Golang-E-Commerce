package handlers

import (
	"golangfinal/database"
	"golangfinal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowMainPage(c *gin.Context) {
	// Retrieve all products from the database
	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var products []models.Product
	err = db.Find(&products).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Optional filtering based on price and rating
	priceSort := c.Query("price")
	ratingSort := c.Query("rating")
	if priceSort == "low" {
		db.Order("price ASC").Find(&products)
	} else if priceSort == "high" {
		db.Order("price DESC").Find(&products)
	}
	if ratingSort == "low" {
		db.Order("rating ASC").Find(&products)
	} else if ratingSort == "high" {
		db.Order("rating DESC").Find(&products)
	}

	// Render the HTML template with the filtered products
	c.HTML(http.StatusOK, "mainpage.html", gin.H{
		"Products":   products,
		"PriceSort":  priceSort,
		"RatingSort": ratingSort,
	})
}

func ShowCart(c *gin.Context) {
	// Get the user ID from the session
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Retrieve the user's products from the database
	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user models.User
	if err := db.Preload("Cart").First(&user, userID.(uint)).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Calculate the total cost of the user's items
	total := 0.0
	for _, product := range user.Cart {
		total += product.Price * float64(product.Quantity)
	}

	// Render the HTML template with the user's cart items and total cost
	c.HTML(http.StatusOK, "cart.html", gin.H{
		"Items": user.Cart,
		"Total": total,
	})
}

func CreateProductHandler(c *gin.Context) {
	// Parse the request body
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Create the new product
	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = db.Create(&product).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Return the created product in the response
	c.JSON(http.StatusCreated, gin.H{"data": product})
}
