package routes

import (
	"gin-template/database"
	"gin-template/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "The provided request is missing some fields",
			"details": err.Error(),
		})
		return
	}
	db := database.DB

	var existingUser models.User
	result := db.Where("PhoneNumber = ? OR email =?", user.PhoneNumber, user.Email).First(&existingUser)

	// Check if the query found a record (meaning the user exists).
	// If result.Error is nil, it means a record was found.
	// If result.Error is gorm.ErrRecordNotFound, no record was found, which is expected for a new user.
	if result.Error == nil {
		// If result.Error is nil, a user with the same phone or email exists.
		c.JSON(http.StatusConflict, gin.H{ // Use StatusConflict for resource conflicts
			"message": "A user with this phone number or email already exists",
		})
		return
	}

	// If no existing user was found (result.Error is ErrRecordNotFound), proceed to create the new user.
	// The Create method returns a *gorm.DB instance. Check its Error field.
	if createResult := db.Create(&user); createResult.Error != nil {
		// If there's an error during creation, return a 500 Internal Server Error.
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user in the database",
			"details": createResult.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)

}
