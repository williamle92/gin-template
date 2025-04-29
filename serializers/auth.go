package serializers

import "github.com/golang-jwt/jwt/v5"

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginOut struct {
	UserId      uint   `json:"user_id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
	Token string `json:"token"`
}

type LoginResponse struct {
	Data       *LoginOut `json:"data"`
	StatusCode int       `json:"status_code"`
	Message    string    `json:"message"`
}
