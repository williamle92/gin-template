package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint    `gorm:"index" json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `gorm:"unique" json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Password    string  `json:"password"`
}
