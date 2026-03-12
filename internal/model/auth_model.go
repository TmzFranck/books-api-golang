package model

type UserInfoResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	User         UserInfoResponse `json:"user"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token,omitempty"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,max=40,email"`
	Password string `json:"password" validate:"required,min=6"`
}
