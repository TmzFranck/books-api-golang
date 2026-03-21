package usecase

import (
	"context"
	"fmt"

	"github.com/TmzFranck/books-api-golang/internal/delivery/http/middleware"
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/jobs"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"github.com/TmzFranck/books-api-golang/internal/model/converter"
	"github.com/TmzFranck/books-api-golang/internal/repository"
	"github.com/TmzFranck/books-api-golang/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) CreateUser(cx context.Context, viper *viper.Viper, worker *jobs.WokerPool, request *model.UserCreateRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, utils.ErrBadRequest
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Errorf("Error hashing password: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	user := &entity.User{
		Username:     request.Username,
		Firstname:    request.FirstName,
		Lastname:     request.Lastname,
		Role:         request.Role,
		Email:        request.Email,
		PasswordHash: string(password),
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Errorf("Error creating user: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	token, err := utils.GenerateURLSafeToken(user.Email)
	if err != nil {
		c.Log.Errorf("Error generating url safe token for user %s: %+v", user.Email, err)
		return nil, utils.ErrInternalServerError
	}

	link := fmt.Sprintf("http://%s/api/v1/auth/verify?token=%s", viper.GetString("server.domain"), token)

	html_message := fmt.Sprintf("<p>Please click the following link to verify your email: <a href=\"%s\">%s</a></p>", link, link)

	mail := &utils.Mail{
		Sender:  "Books API",
		To:      []string{user.Email},
		Subject: "Verify your email",
		Body:    html_message,
	}

	if err = utils.SendMail(worker, mail); err != nil {
		c.Log.Errorf("Error sending verification email: %+v", err)
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) LoginUser(cx context.Context, request model.UserLoginRequest) (*model.LoginResponse, error) {
	tx := c.DB.WithContext(cx)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, utils.ErrBadRequest
	}

	user := &entity.User{}
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Errorf("Error fetching user: %+v", err)
		return nil, utils.ErrBadRequest
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		c.Log.Errorf("Invalid password for user %s: %+v", request.Email, err)
		return nil, utils.ErrUnauthorized
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email)
	if err != nil {
		c.Log.Errorf("Error generating tokens for user %s: %+v", request.Email, err)
		return nil, utils.ErrInternalServerError
	}

	return &model.LoginResponse{
		User:         *converter.UserToInfoResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *UserUseCase) GetNewTokens(cx context.Context, refreshToken string) (*model.LoginResponse, error) {
	tx := c.DB.WithContext(cx)

	claims, err := utils.ValidateToken(refreshToken)
	if err != nil || claims.Refresh == false {
		c.Log.Errorf("Invalid refresh token: %+v", err)
		return nil, utils.ErrUnauthorized
	}

	user := &entity.User{}
	if err := c.UserRepository.FindByEmail(tx, user, claims.UserEmail); err != nil {
		c.Log.Errorf("Error fetching user for token refresh: %+v", err)
		return nil, utils.ErrUnauthorized
	}

	accessToken, err := utils.RefreshToken(refreshToken)

	return &model.LoginResponse{
		User:        *converter.UserToInfoResponse(user),
		AccessToken: accessToken,
	}, nil

}

func (c *UserUseCase) VerifyEmail(cx context.Context, token string) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return utils.ErrUnauthorized
	}

	tx := c.DB.WithContext(cx)
	user := &entity.User{}
	if err := c.UserRepository.FindByEmail(tx, user, claims.UserEmail); err != nil {
		return utils.ErrUnauthorized
	}

	user.IsVerified = true
	if err := c.UserRepository.Update(tx, user); err != nil {
		return utils.ErrInternalServerError
	}

	return nil
}

func (c *UserUseCase) PasswordResetRequest(cx context.Context, viper *viper.Viper, worker *jobs.WokerPool, request *model.PasswordResetRequest) error {
	tx := c.DB.WithContext(cx)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return utils.ErrBadRequest
	}

	user := &entity.User{}
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		return utils.ErrNotFound
	}

	token, err := utils.GenerateURLSafeToken(user.Email)
	if err != nil {
		return utils.ErrInternalServerError
	}

	link := fmt.Sprintf("http://%s/api/auth/password-reset-confirm?token=%s", viper.GetString("server.domain"), token)

	html_message := fmt.Sprintf("<p>Please click the following link to reset your password: <a href=\"%s\">%s</a></p>", link, link)

	mail := utils.Mail{
		Sender:  "Books Api",
		To:      []string{user.Email},
		Subject: "Password Reset",
		Body:    html_message,
	}

	err = utils.SendMail(worker, &mail)
	if err != nil {
		return utils.ErrInternalServerError
	}

	return nil
}

func (c *UserUseCase) PasswordResetConfirm(cx context.Context, token string, request *model.PasswordResetConfirmationRequest) error {
	tx := c.DB.WithContext(cx)

	claims, err := utils.ValidateURLSafeToken(token)
	if err != nil {
		return utils.ErrBadRequest
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return utils.ErrBadRequest
	}

	user := &entity.User{}
	if err := c.UserRepository.FindByEmail(tx, user, claims.Usermail); err != nil {
		return utils.ErrNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Errorf("Error hashing password: %+v", err)
		return utils.ErrInternalServerError
	}
	user.PasswordHash = string(hashedPassword)
	if err := c.UserRepository.Update(tx, user); err != nil {
		return utils.ErrInternalServerError
	}

	return nil
}

func (c *UserUseCase) RevokeToken(cx context.Context, redisClient *redis.Client, token string) error {

	if err := utils.AddJwtToBlacklist(cx, redisClient, token); err != nil {
		return err
	}

	return nil
}

func (c *UserUseCase) SendMail(cx context.Context, worker *jobs.WokerPool, request *model.EmailRequest) error {
	if c.Validate.Struct(request) != nil {
		return utils.ErrBadRequest
	}

	mail := &utils.Mail{
		Sender:  "Books Api",
		To:      request.Addresses,
		Subject: "Welcome message",
		Body:    "<h1>Welcome to Books Api!</h1>",
	}

	if err := utils.SendMail(worker, mail); err != nil {
		return utils.ErrInternalServerError
	}

	return nil
}

func (c *UserUseCase) GetCurrentUser(cx context.Context) (*model.UserResponse, error) {
	tx := c.DB.WithContext(cx)

	usermail := middleware.GetUserEmail(cx)

	user := &entity.User{}
	if err := c.UserRepository.FindByEmail(tx, user, usermail); err != nil {
		return nil, utils.ErrNotFound
	}

	return converter.UserToResponse(user), nil
}
