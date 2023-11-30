package model

// RegistrationInput: represents expected data from a registration request.
type RegistrationInput struct {
	Fname    string `json:"first_name" binding:"required"`
	Lname    string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone_number" binding:"required,e164"`
	Password string `json:"password" binding:"required"`
}

// LoginInput: represents expected data from a login request.
type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
