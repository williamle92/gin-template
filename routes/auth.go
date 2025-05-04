package routes

import (
	"gin-template/database"
	"gin-template/models"
	"gin-template/serializers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

// LoginUser handles the HTTP POST request for user login.
// It validates the incoming credentials (email, password), verifies them against the database,
// and if successful, generates and returns a JSON Web Token (JWT).
//
// Parameters:
//   - c: *gin.Context - The context object for the incoming HTTP request.
func LoginUser(c *gin.Context) {
	// Declare a variable `request` of type serializers.Login to hold the parsed request body.
	var request serializers.Login

	// Attempt to bind the incoming JSON request body to the `request` struct.
	// This also performs validation based on the binding tags in serializers.Login.
	if err := c.ShouldBindJSON(&request); err != nil {
		// If binding or validation fails, respond with 422 Unprocessable Entity.
		c.JSON(http.StatusUnprocessableEntity, serializers.LoginResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Validation failed: " + err.Error(), // Include validation error details.
		})
		// Log the specific validation error.
		log.Println(err.Error())
		return
	}

	// Get the database connection instance.
	db := database.DB

	// Attempt to find a user in the database matching the provided email.
	var ExistingUser models.User
	// GORM query to find the first user where the email matches the request email.
	result := db.Where("email = ?", request.Email).First(&ExistingUser)

	// If there was an error with getting the user
	if result.Error != nil {
		// Check if the error is specifically the user not being found".
		if result.Error == gorm.ErrRecordNotFound {
			// If the user is not found, respond with 401 Unauthorized.
			// Use a generic "Invalid credentials" message for security (don't reveal if email exists).
			c.JSON(http.StatusUnauthorized, serializers.LoginResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid credentials",
			})
			// Log the actual error (record not found).
			log.Println(result.Error)
			return
		} else {
			// If a different database error occurred during lookup, respond with 500 Internal Server Error.
			c.JSON(http.StatusInternalServerError, serializers.LoginResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "An error occurred while verifying credentials",
			})
			// Log the unexpected database error.
			log.Println("Database error during login lookup:", result.Error)
			return
		}
	}

	// If the user was found (result.Error == nil), compare the provided password with the stored hash.
	// bcrypt.CompareHashAndPassword handles the comparison securely.
	err := bcrypt.CompareHashAndPassword([]byte(ExistingUser.Password), []byte(request.Password))
	if err != nil {
		// If the password does not match
		c.JSON(http.StatusUnauthorized, serializers.LoginResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid credentials", // Again, use a generic message.
		})
		// Log the password mismatch error.
		log.Println("password provided does not match", err.Error())
		return
	}

	// --- Password verification successful ---

	// 3. Create JWT Claims: Define the payload of the JWT.
	// Set the token expiration time (e.g., 24 hours from now).
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the claims struct, including custom claims (UserId, Email, PhoneNumber)
	// and standard registered claims (ExpiresAt, IssuedAt).
	claims := &serializers.LoginOut{
		UserId:      ExistingUser.ID,
		Email:       ExistingUser.Email,
		PhoneNumber: ExistingUser.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			// Set the expiration time for the token.
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			// Set the time the token was issued.
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// 4. Create JWT Token object: Create a new token instance using the HS256 signing method and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Retrieve the JWT secret key from environment variables.
	jwtsecret := os.Getenv("JWT_SALT")
	if jwtsecret == "" {
		// Handle the case where the secret is missing.
		c.JSON(http.StatusInternalServerError, serializers.LoginResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "JWT secret configuration error",
		})
		log.Println("CRITICAL: JWT_SALT environment variable not set.")
		return
	}
	jwtbytes := []byte(jwtsecret) // Convert the secret string to a byte slice for signing.

	// 5. Sign the token: Generate the final token string by signing the token object with the secret key.
	tokenString, err := token.SignedString(jwtbytes)
	if err != nil {
		// If signing fails (e.g., issues with the claims or secret), respond with 500 Internal Server Error.
		c.JSON(http.StatusInternalServerError, serializers.LoginResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to generate authentication token",
		})
		// Log the signing error.
		log.Println("Failed to sign JWT:", err)
		return
	}

	// --- Successfully logged in and generated token ---

	// Add the generated token string to the claims struct for inclusion in the response body.
	claims.Token = tokenString

	// Send the successful response (200 OK) including the user data and the JWT.
	c.JSON(http.StatusOK, serializers.LoginResponse{
		StatusCode: http.StatusOK,
		Message:    "Successfully logged in",
		Data:       claims, // Send the LoginOut struct (which now includes the token) as data.
	})

	// Log the successful login event.
	log.Println("User logged in successfully:", ExistingUser.Email)
}
