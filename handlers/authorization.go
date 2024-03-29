package handlers

import (
	"errors"
	"golangfinal/database"
	"golangfinal/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func RegisterHandler(c *gin.Context) {

	firstName := c.PostForm("first_name")
	lastName := c.PostForm("last_name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

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

	if password != confirmPassword {
		c.AbortWithError(http.StatusBadRequest, errors.New("passwords do not match"))
		return
	}

	err = user.HashPassword()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

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
	c.Redirect(http.StatusFound, "/login")

}

func LoginHandler(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user := &models.User{
		Email:    email,
		Password: password,
	}
	err := user.Validate()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

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

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid email or password"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"user_id": existingUser.ID, 
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("token", tokenString, int(time.Hour*24/time.Second), "/", "", false, true)

	c.Redirect(http.StatusFound, "/")
}


func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
