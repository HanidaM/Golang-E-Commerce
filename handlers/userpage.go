package handlers

import (
	"errors"
	"fmt"
	"golangfinal/database"
	"golangfinal/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ShowMainPage(c *gin.Context) {

	var products []models.Product
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	tokenString, err := c.Cookie("token")
	if err != nil {
		c.HTML(http.StatusOK, "mainpage.html", gin.H{
			"Products":        products,
			"IsAuthenticated": false,
		})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte("secret"), nil
	})
	if err != nil {
		c.HTML(http.StatusOK, "mainpage.html", gin.H{
			"Products":        products,
			"IsAuthenticated": false,
		})
		return
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		c.HTML(http.StatusOK, "mainpage.html", gin.H{
			"Products":        products,
			"IsAuthenticated": false,
		})
		return
	}

	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
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
		"Products":        products,
		"IsAuthenticated": true,
		"SortBy":          sortBy,
		"Order":           order,
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

func GetUserIDFromToken(c *gin.Context) (uint, error) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		return 0, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user ID not found in token claims")
	}

	fmt.Println("User ID:", uint(userID))

	return uint(userID), nil
}

func UpdateProductRating(c *gin.Context) {
	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	productIDStr := c.Param("id")

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := db.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var ratingReq struct {
		Rating float64 `json:"rating"`
	}

	if err := c.ShouldBindJSON(&ratingReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating"})
		return
	}

	product.Rating = ratingReq.Rating
	if err := db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "rating": product.Rating})
}

func GetUserCartItems(c *gin.Context) {
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var cartItems []models.CartItem
	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := db.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "cart.html", gin.H{
		"CartItems": cartItems,
	})
}

func AddToCart(c *gin.Context) {

	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var requestData struct {
		ProductID uint `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartItem := models.CartItem{
		UserID:    userID,
		ProductID: requestData.ProductID,
		Quantity:  1,
	}

	if err := db.Create(&cartItem).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
}

func HandleComment(c *gin.Context) {
	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	productID, err := strconv.Atoi(c.PostForm("product_id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tokenString, err := c.Cookie("token")
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("user not authenticated"))
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Replace "secret" with your own secret key
	})
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("failed to parse token"))
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid token claims"))
		return
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid user ID"))
		return
	}

	text := c.PostForm("comment")

	comment := &models.Comment{
		UserID:    uint(userID),
		ProductID: uint(productID),
		Text:      text,
	}

	err = db.Create(comment).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/products/%d", productID))
}
