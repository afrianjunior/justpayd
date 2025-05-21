package auth

// LoginRequest represents the credentials provided for login
type LoginRequest struct {
	Email string `json:"email" example:"user@example.com"`
}

// LoginResponse represents the response after successful authentication
type LoginResponse struct {
	Token  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserID int    `json:"user_id" example:"1"`
}
