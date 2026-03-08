package usecase

import (
	"log"

	"github.com/TmzFranck/books-api-golang/internal/repository"
	"gorm.io/gorm"
)

type UserUserCase struct {
	DB             *gorm.DB
	Log            *log.Logger
	UserRepository *repository.UserRepository
}

func NewUserUserCase(db *gorm.DB, log *log.Logger, userRepository *repository.UserRepository) *UserUserCase {
	return &UserUserCase{
		DB:             db,
		Log:            log,
		UserRepository: userRepository,
	}
}
