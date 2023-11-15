package model

// AuthenticationInput: represents the expected data from an authentication request.
type AuthenticationInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
