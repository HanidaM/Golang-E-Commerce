package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	FirstName string `json:"first_name" validate:"required,min=2,max=30"`
	LastName  string `json:"last_name" validate:"required,min=2,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Product struct {
	gorm.Model
	SKU         string `gorm:"uniqueIndex;not null"`
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Price       float64
	Rating      float64
	Image       string
}

// Validate validates the user fields
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
	if p.SKU == "" {
		return errors.New("sku required")
	}
	if p.Name == "" {
		return errors.New("name required")
	}
	if p.Description == "" {
		return errors.New("description required")
	}
	return nil
}

// HashPassword hashes the user's password and sets it as the hashed password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
