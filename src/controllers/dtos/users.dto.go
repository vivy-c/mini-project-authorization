package dtos

import "time"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type AmendProfileInputRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AmendProfilePasswordRequest struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
