package handlers

import (
	"golangfinal/database"
	"golangfinal/models"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var product models.Product
	err = db.First(&product, id).Error
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func GetAllProducts(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func UpdateProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var product models.Product
	err = db.First(&product, id).Error
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = db.Save(&product).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func DeleteProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var product models.Product
	err = db.First(&product, id).Error
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	err = db.Delete(&product).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func SearchProductByTitle(c *gin.Context) {
	searchQuery := c.Query("title")

	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var products []models.Product
	err = db.Where("title LIKE ?", "%"+searchQuery+"%").Find(&products).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "mainpage.html", gin.H{
		"search_query": searchQuery,
		"products":     products})

	c.JSON(http.StatusOK, gin.H{"data": products})

}

func GetSortedProducts(c *gin.Context) {
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
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	switch sortBy {
	case "cost":
		if order == "desc" {
			sort.Slice(products, func(i, j int) bool {
				return products[i].Price > products[j].Price
			})
		} else {
			sort.Slice(products, func(i, j int) bool {
				return products[i].Price < products[j].Price
			})
		}
	default:
		if order == "desc" {
			sort.Slice(products, func(i, j int) bool {
				return products[i].ID > products[j].ID
			})
		} else {
			sort.Slice(products, func(i, j int) bool {
				return products[i].ID < products[j].ID
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": products})

}
