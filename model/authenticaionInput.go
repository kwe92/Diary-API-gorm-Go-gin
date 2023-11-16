package model

// AuthenticationInput: represents expected data from authentication request.
type AuthenticationInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
