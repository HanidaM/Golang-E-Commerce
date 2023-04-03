package handlers

import (
	"golangfinal/database"
	"golangfinal/models"
	"math/rand"
	"net/http"
	"strconv"

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
		"Products": products,
		"PriceSort": priceSort,
		"RatingSort": ratingSort,
	})
}

func CreateProductHandler(c *gin.Context) {
	// Parse the request body
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Generate SKU for the new product
	product.SKU = generateSKU()

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

// generateSKU generates a unique SKU for a new product
func generateSKU() string {
	// Generate a random 6-digit number
	num := rand.Intn(900000) + 100000

	// Convert the number to a string and return it as the SKU
	return strconv.Itoa(num)
}
