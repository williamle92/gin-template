package serializers

import "github.com/golang-jwt/jwt/v5"

// Login defines the structure for expected data in a login request body.
// It uses Gin's binding tags for validation.
type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginOut defines the structure for the data included within a successful login response.
// It includes user details and the JWT token, embedding standard JWT claims.
type LoginOut struct {
	UserId      uint   `json:"user_id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
	Token string `json:"token"`
}

// LoginResponse defines the overall structure of the JSON response for the login endpoint.
// It includes status code, a message, and the actual login data (LoginOut) on success.
type LoginResponse struct {
	Data       *LoginOut `json:"data"`
	StatusCode int       `json:"status_code"`
	Message    string    `json:"message"`
}
