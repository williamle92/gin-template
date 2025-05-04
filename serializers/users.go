package serializers

// UserIn defines the structure for expected data when creating or potentially updating a user.
// It uses Gin's binding tags for input validation.
type UserIn struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,len=10"`
	Password    string `json:"password" binding:"required,min=6"`
}

// UserOut defines the structure for user data returned in API responses.
// It intentionally omits sensitive fields like the password.
type UserOut struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// UserResponse defines the generic structure for JSON responses related to user operations.
// It includes a status code, a message, and optional data payload.
type UserResponse struct {
	Data       any    `json:"data,omitempty"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
