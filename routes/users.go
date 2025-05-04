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

// CreateUser handles the HTTP POST request to register a new user.
// It validates the incoming request data, checks for existing users with the same email or phone number,
// hashes the password, creates the user record in the database, and returns the created user's details.
//
// Parameters:
//   - c: *gin.Context - The context object for the incoming HTTP request, providing access to request data and response writing.
func CreateUser(c *gin.Context) {
	// Declare a variable `request` of type serializers.UserIn to hold the parsed request body.
	var request serializers.UserIn

	// Attempt to bind the incoming JSON request body to the `request` struct.
	// `ShouldBindJSON` also performs validation based on the binding tags in serializers.UserIn.
	if err := c.ShouldBindJSON(&request); err != nil {
		// If binding or validation fails, respond with a 422 Unprocessable Entity status.
		// Include a generic validation error message in the JSON response.
		c.JSON(http.StatusUnprocessableEntity, serializers.UserResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Validation failed: " + err.Error(), // Consider adding more specific error details from err if appropriate
		})
		// Log the actual validation error for debugging purposes.
		log.Println(err.Error())
		return
	}

	// Get the database connection instance established during application startup.
	db := database.DB

	// Check if a user with the provided email or phone number already exists.
	var existingUser models.User
	// GORM query to find the first user matching either the email or phone number.
	result := db.Where("phone_number = ? OR email = ?", request.PhoneNumber, request.Email).First(&existingUser)

	if result.Error == nil {
		// If `result.Error` is nil, it means a user was found (record exists).
		// Respond with a 400 Bad Request status, indicating a user with the
		// same email or phone number already exists.
		c.JSON(http.StatusBadRequest, serializers.UserResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "A user with this phone number or email already exists",
		})
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		// If `result.Error` is not nil AND it's not the specific "record not found" error,
		// it indicates an unexpected database error occurred during the lookup.
		// Respond with a 500 Internal Server Error status.
		c.JSON(http.StatusInternalServerError, serializers.UserResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong with your request",
		})
		// Log the actual database error.
		log.Println("Error checking for existing user:", result.Error)
		return
	}
	// If we reach here, it means `result.Error == gorm.ErrRecordNotFound`,
	// which is the expected outcome for a new user registration (no existing user found).

	// Hash the user's provided password using bcrypt for secure storage.
	// bcrypt.DefaultCost provides a good balance between security and performance.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		// If hashing fails, respond with a 500 Internal Server Error.
		c.JSON(http.StatusInternalServerError, serializers.UserResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to process password",
		})
		// Log the hashing error.
		log.Println("Error hashing password:", err.Error())
		return
	}

	// Create a new models.User instance with the data from the request.
	newUser := models.User{
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
		Password:    string(hashedPassword), // Store the hashed password string.
	}

	// Attempt to save the new user record to the database.
	// The Create method returns a *gorm.DB instance; check its Error field for issues.
	if createResult := db.Create(&newUser); createResult.Error != nil {
		// If there's an error during the database insertion, respond with a 500 Internal Server Error.
		c.JSON(http.StatusInternalServerError, serializers.UserResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong with your request",
		})
		// Log the database creation error.
		log.Println(createResult.Error)
		return
	}

	// If user creation is successful, prepare the response data.
	// Use serializers.UserOut to control which fields are included in the response (excluding the password).
	response := serializers.UserResponse{
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
		Data: serializers.UserOut{
			ID:          newUser.ID,
			FirstName:   newUser.FirstName,
			LastName:    newUser.LastName,
			Email:       newUser.Email,
			PhoneNumber: newUser.PhoneNumber,
		},
	}
	// Send the successful JSON response with a 201 Created status.
	c.JSON(http.StatusCreated, response)

}
