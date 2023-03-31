package handlers

import (
	"errors"
	"golangfinal/database"
	"golangfinal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
func ShowCart(c *gin.Context){
	c.HTML(http.StatusOK,"cart.html",gin.H{})
}

func RegisterHandler(c *gin.Context) {

	// Parse the form data
	firstName := c.PostForm("first_name")
	lastName := c.PostForm("last_name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// Validate the form data
	user := &models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
	err := user.Validate()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check if password and confirm password match
	if password != confirmPassword {
		c.AbortWithError(http.StatusBadRequest, errors.New("passwords do not match"))
		return
	}

	// Hash the user's password
	err = user.HashPassword()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create a new user
	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = db.Create(user).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// // Return success message
	// // Redirect to the login page on success
	// message := gin.H{"message": "user created"}
	// c.JSON(http.StatusCreated, message)

	// time.AfterFunc(3*time.Second, func() {
	// 	c.Redirect(http.StatusFound, "/login")
	// })
	c.Redirect(http.StatusFound, "/login")

}

func LoginHandler(c *gin.Context) {
	// Parse the form data
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate the form data
	user := &models.User{
		Email:    email,
		Password: password,
	}
	err := user.Validate()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check if user exists in the database
	db, err := database.ConnectDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var existingUser models.User
	err = db.Where("email = ?", email).First(&existingUser).Error
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid email or password"))
		return
	}

	// Check if password matches
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid email or password"))
		return
	}

	// Redirect to the dashboard on success
	c.Redirect(http.StatusFound, "/main")
}
