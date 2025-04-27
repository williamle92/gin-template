package routes

import (
	"gin-template/database"
	"gin-template/models"
	"gin-template/serializers"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var request serializers.UserIn

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, serializers.UserResponse{
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
		c.JSON(http.StatusBadRequest, serializers.UserResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "A user with this phone number or email already exists",
		})
		return
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		// If result.Error is nil, a user with the same phone or email exists.
		c.JSON(http.StatusInternalServerError, serializers.UserResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong with your request",
		})
		log.Println(result.Error)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, serializers.UserResponse{
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
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
		Password:    string(hashedPassword), // Assign the hashed password
	}

	// If no existing user was found (result.Error is ErrRecordNotFound), proceed to create the new user.
	// The Create method returns a *gorm.DB instance. Check its Error field.
	if createResult := db.Create(&newUser); createResult.Error != nil {
		// If there's an error during creation, return a 500 Internal Server Error.
		c.JSON(http.StatusInternalServerError, serializers.UserResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong with your request",
		})
		log.Println(createResult.Error)
		return
	}
	response := serializers.UserResponse{
		StatusCode: 201,
		Message:    "User created successfully",
		Data: serializers.UserOut{
			ID:          newUser.ID,
			FirstName:   newUser.FirstName,
			LastName:    newUser.LastName,
			Email:       newUser.Email,
			PhoneNumber: newUser.PhoneNumber,
		},
	}
	c.JSON(http.StatusCreated, response)

}
