package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint `gorm:"index"`
	FirstName   string
	LastName    string
	Email       string `gorm:"unique"`
	PhoneNumber *string
	Password    string
}
