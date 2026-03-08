package repository

import (
	"log"

	"github.com/TmzFranck/books-api-golang/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *log.Logger
}

func NewUserRepository(log *log.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).First(user).Error
}
