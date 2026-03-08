package model

import "time"

type UserResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Firstname  string    `json:"firstname"`
	Lastname   string    `json:"lastname"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserCreateRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,max=40,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserBooksResponse struct {
	UserResponse
	Books []BookResponse `json:"books"`
}

type EmailResponse struct {
	Addresses []string `json:"addresses"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,max=40,email"`
}

type PasswordResetConfirmationRequest struct {
	NewPassword        string `json:"new_password" validate:"required,min=6"`
	ConfirmNewPassword string `json:"Confirm_new_password" validate:"required,min=6"`
}
