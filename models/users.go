package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement;index" json:"id"`
	FirstName   string `json:"first_name" gorm:"not null"`
	LastName    string `json:"last_name" gorm:"not null"`
	Email       string `gorm:"unique;index;not null" json:"email"`
	PhoneNumber string `gorm:"unique;index;not null" json:"phone_number"`
	Password    string `json:"-"`
}
