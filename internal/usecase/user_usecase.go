package usecase

import (
	"context"

	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"github.com/TmzFranck/books-api-golang/internal/model/converter"
	"github.com/TmzFranck/books-api-golang/internal/repository"
	"github.com/TmzFranck/books-api-golang/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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

func (c *UserUseCase) CreateUser(cx context.Context, request *model.UserCreateRequest) (*model.UserResponse, error) {
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

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) LoginUser(cx context.Context, request model.UserLoginRequest) (*model.LoginResponse, error) {
	tx := c.DB.WithContext(cx)

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
