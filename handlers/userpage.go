package handlers

import (
	"golangfinal/database"
	"golangfinal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowMainPage(c *gin.Context) {
	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var products []models.Product
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	switch sortBy {
	case "cost":
		if order == "asc" {
			db.Order("price ASC").Find(&products)
		} else if order == "desc" {
			db.Order("price DESC").Find(&products)
		} else {
			db.Find(&products)
		}
	case "rating":
		if order == "asc" {
			db.Order("rating ASC").Find(&products)
		} else if order == "desc" {
			db.Order("rating DESC").Find(&products)
		} else {
			db.Find(&products)
		}
	default:
		db.Find(&products)
	}

	c.HTML(http.StatusOK, "mainpage.html", gin.H{
		"Products": products,
		"SortBy":   sortBy,
		"Order":    order,
	})
}


func ShowCart(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

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

	total := 0.0
	for _, product := range user.Cart {
		total += product.Price * float64(product.Quantity)
	}

	c.HTML(http.StatusOK, "cart.html", gin.H{
		"Items": user.Cart,
		"Total": total,
	})
}

func CreateProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

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

	c.JSON(http.StatusCreated, gin.H{"data": product})
}
