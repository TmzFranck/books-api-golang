package http

import (
	"encoding/json"
	"net/http"

	"github.com/TmzFranck/books-api-golang/internal/delivery/http/middleware"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"github.com/TmzFranck/books-api-golang/internal/usecase"
	"github.com/TmzFranck/books-api-golang/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log       *logrus.Logger
	useCase   *usecase.UserUseCase
	validator *validator.Validate
}

func NewUserController(log *logrus.Logger, useCase *usecase.UserUseCase, validator *validator.Validate) *UserController {
	return &UserController{
		Log:       log,
		useCase:   useCase,
		validator: validator,
	}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	request := new(model.UserCreateRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("failed to decode resquest body: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("validation failed for user registration: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	c.Log.Infof("registering user: %s", request.Username)

	response, err := c.useCase.CreateUser(r.Context(), request)
	if err != nil {
		c.Log.Errorf("failed to create user: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	c.Log.Infof("user %s registered successfully", request.Username)
	utils.ResponseWithData(w, http.StatusCreated, response)

}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	request := new(model.UserLoginRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("failed to decode request body: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("validation failed for user login: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	c.Log.Infof("processing login for email: %s", request.Email)

	response, err := c.useCase.LoginUser(r.Context(), request)
	if err != nil {
		c.Log.Errorf("login use case failed: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "failed to procces login")
		return
	}

	c.Log.Infof("user logged in successfully: %v", response.User.ID)
	utils.ResponseWithData(w, http.StatusOK, response)
}

func (c *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetUserToken(r.Context())

	userID := middleware.GetUserID(r.Context())

	c.Log.Infof("processing logout for user: %d", userID)
	if err := c.useCase.RevokeToken(r.Context(), token); err != nil {
		c.Log.Errorf("logout use case failed: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	c.Log.Infof("user logged out successfully user: %d", userID)
	utils.ResponseWithData(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

func (c *UserController) VerifyUserAccount(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	c.Log.Infof("processing verify user account for token: %s", token)
	if err := c.useCase.VerifyEmail(r.Context(), token); err != nil {
		c.Log.Errorf("verify user account failed: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to verify user account")
		return
	}

	c.Log.Infof("email verified successfully for token: %s", token)
	utils.ResponseWithData(w, http.StatusOK, map[string]string{"message": "email verified successfully"})
}

func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {

	c.Log.Infof("fetching current user")
	user, err := c.useCase.GetCurrentUser(r.Context())
	if err != nil {
		c.Log.Errorf("get current user failed: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to get current user")
		return
	}

	c.Log.Infof("current user: %s", user.Username)
	utils.ResponseWithData(w, http.StatusOK, user)
}

func (c *UserController) SendMail(w http.ResponseWriter, r *http.Request) {
	request := new(model.EmailRequest)
	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("validator failed: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	c.Log.Infof("sending mail to addresses: %v", request.Addresses)
	if err := c.useCase.SendMail(r.Context(), request); err != nil {
		c.Log.Errorf("send mail failed: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to send mail")
		return
	}

	c.Log.Infof("mail sent successfully for addresses: %v", request.Addresses)
	utils.ResponseWithData(w, http.StatusOK, map[string]string{"message": "mail sent successfully"})
}

func (c *UserController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetUserToken(r.Context())
	userID := middleware.GetUserID(r.Context())

	c.Log.Infof("refreshing token for user: %d", userID)
	accessToken, err := c.useCase.GetNewAccessToken(r.Context(), token)
	if err != nil {
		c.Log.Errorf("refresh token failed: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to refresh token")
		return
	}

	c.Log.Infof("token refreshed successfully for user: %s", token)
	utils.ResponseWithData(w, http.StatusOK, accessToken)
}

func (c *UserController) PasswordReset(w http.ResponseWriter, r *http.Request) {
	request := new(model.PasswordResetRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Errorf("password reset failed: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	c.Log.Infof("password reset request received for user: %s", request.Email)
	if err := c.useCase.PasswordResetRequest(r.Context(), request); err != nil {
		c.Log.Errorf("password reset failed: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to reset password")
		return
	}

	c.Log.Infof("password reset successfully for user: %s", request.Email)
	utils.ResponseWithData(w, http.StatusOK, map[string]string{"message": "password reset successfully"})
}

func (c *UserController) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	request := new(model.PasswordResetConfirmationRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("password reset confirmation failed: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("password reset confirmation failed: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if request.NewPassword != request.ConfirmNewPassword {
		c.Log.Warnf("password reset confirmation failed: %v", "passwords do not match")
		utils.ResponseWithError(w, http.StatusBadRequest, "passwords do not match")
		return
	}

	c.Log.Infof("password reset confirmed for user: %s", token)
	if err := c.useCase.PasswordResetConfirm(r.Context(), token, request); err != nil {
		c.Log.Errorf("password reset confirmation failed: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to confirm password reset")
		return
	}

	c.Log.Infof("password reset confirmed for user: %s", token)
	utils.ResponseWithData(w, http.StatusOK, map[string]string{"message": "password reset confirmed"})
}
