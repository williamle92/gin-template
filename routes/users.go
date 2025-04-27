package routes

import (
	"gin-template/database"
	"gin-template/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type UserIn struct {
	FirstName   string `json:"first_name" binding:"required"`          // FirstName is also likely required? Added binding here.
	LastName    string `json:"last_name" binding:"required"`           // LastName also likely required? Added binding here.
	Email       string `json:"email" binding:"required,email"`         // "required" ensures it's present and non-empty, "email" validates format
	PhoneNumber string `json:"phone_number" binding:"required,len=10"` // "required" ensures it's present and non-empty
	Password    string `json:"password" binding:"required,min=6"`      // Ensure password is required and has minimum length
}
type UserOut struct {
	ID          uint    `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type UserResponse struct {
	Data       any    `json:"data,omitempty"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func CreateUser(c *gin.Context) {
	var request UserIn

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, UserResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Validation failed",
		})
		log.Println(err.Error())
		return
	}
	db := database.DB

	var existingUser models.User
	result := db.Where("phone_number = ? OR email =?", request.PhoneNumber, request.Email).First(&existingUser)

	// Check if the query found a record (meaning the user exists).
	// If result.Error is nil, it means a record was found.
	// If result.Error is gorm.ErrRecordNotFound, no record was found, which is expected for a new user.
	if result == nil {
		// If result.Error is nil, a user with the same phone or email exists.
		c.JSON(http.StatusBadRequest, UserResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "A user with this phone number or email already exists",
		})
		return
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		// If result.Error is nil, a user with the same phone or email exists.
		c.JSON(http.StatusInternalServerError, UserResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong with your request",
		})
		log.Println(result.Error)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, UserResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Invalid password",
		})
		log.Println(err.Error())
		return
	}
	newUser := models.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		// Take the address of the input string to assign to the *string model field
		Email:       &request.Email,
		PhoneNumber: &request.PhoneNumber,
		Password:    string(hashedPassword), // Assign the hashed password
	}

	// If no existing user was found (result.Error is ErrRecordNotFound), proceed to create the new user.
	// The Create method returns a *gorm.DB instance. Check its Error field.
	if createResult := db.Create(&newUser); createResult.Error != nil {
		// If there's an error during creation, return a 500 Internal Server Error.
		c.JSON(http.StatusInternalServerError, UserResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong with your request",
		})
		log.Println(createResult.Error)
		return
	}
	response := UserResponse{
		StatusCode: 201,
		Message:    "User created successfully",
		Data: UserOut{
			ID:          newUser.ID,
			FirstName:   newUser.FirstName,
			LastName:    newUser.LastName,
			Email:       newUser.Email,
			PhoneNumber: newUser.PhoneNumber,
		},
	}
	c.JSON(http.StatusCreated, response)

}
