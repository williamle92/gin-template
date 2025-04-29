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

func LoginUser(c *gin.Context) {
	var request serializers.Login

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, serializers.LoginResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Validation failed: " + err.Error(),
		})
		log.Println(err.Error())
		return
	}

	db := database.DB

	var ExistingUser models.User
	result := db.Where("email = ?", request.Email).First(&ExistingUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, serializers.LoginResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid credentials",
			})
			log.Println(result.Error)
			return
		} else {
			// Handle other potential database errors
			c.JSON(http.StatusInternalServerError, serializers.LoginResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "An error occured with your request",
			})
			log.Println(result.Error)
			return

		}
	}

	// encrypt password
	err := bcrypt.CompareHashAndPassword([]byte(ExistingUser.Password), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, serializers.LoginResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid password",
		})
		log.Println("password provided does not match", err.Error())
		return
	}

	// 3. Create JWT Claims
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &serializers.LoginOut{
		UserId:      ExistingUser.ID, // Include user ID in claims
		Email:       ExistingUser.Email,
		PhoneNumber: ExistingUser.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 4. Create JWT Token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtsecret := os.Getenv("JWT_SALT")
	jwtbytes := []byte(jwtsecret)
	// 5. Sign the token with the secret key
	tokenString, err := token.SignedString(jwtbytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, serializers.LoginResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "There was an issue with your request",
		})
		log.Println("Failed to sign JWT:", err)
		return
	}

	// --- Successfully logged in and generated token ---

	// --- Choose ONE way to send the token back ---
	claims.Token = tokenString
	// OPTION A: Send token in the response body
	c.JSON(http.StatusOK, serializers.LoginResponse{
		StatusCode: http.StatusOK,
		Message:    "Successfully logged in",
		Data:       claims,
	})

	log.Println("User logged in successfully:", ExistingUser.Email)
}
