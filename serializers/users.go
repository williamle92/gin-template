package serializers

type UserIn struct {
	FirstName   string `json:"first_name" binding:"required"`          // FirstName is also likely required? Added binding here.
	LastName    string `json:"last_name" binding:"required"`           // LastName also likely required? Added binding here.
	Email       string `json:"email" binding:"required,email"`         // "required" ensures it's present and non-empty, "email" validates format
	PhoneNumber string `json:"phone_number" binding:"required,len=10"` // "required" ensures it's present and non-empty
	Password    string `json:"password" binding:"required,min=6"`      // Ensure password is required and has minimum length
}
type UserOut struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type UserResponse struct {
	Data       any    `json:"data,omitempty"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
