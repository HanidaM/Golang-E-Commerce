package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name" validate:"required,min=2,max=30"`
	LastName  string `json:"last_name" validate:"required,min=2,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`

	Cart      []Product      `gorm:"many2many:carts"`
}

type Product struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Rating      float64
	Image       string
	Quantity    int            `gorm:"-"`

}

type CartItem struct {
	gorm.Model
	UserID     uint           `json:"user_id" gorm:"not null"`
	ProductID  uint           `json:"product_id" gorm:"not null"`
	Quantity   int            `json:"quantity" gorm:"not null"`
	TotalPrice float64        `json:"total_price" gorm:"not null"`
	
}


type Comment struct {
    gorm.Model
    UserID    uint
    ProductID uint
    Text      string
}


func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email required")
	}

	if u.Password == "" {
		return errors.New("password required")
	}

	return nil
}

func (p *Product) Validate() error {

	if p.Name == "" {
		return errors.New("name required")
	}
	if p.Description == "" {
		return errors.New("description required")
	}
	if p.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
