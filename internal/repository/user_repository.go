package repository

import (
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).First(user).Error
}

func (r *UserRepository) FindByEmailWithBooks(db *gorm.DB, user *entity.User, email string) error {
	return db.Preload("Books").Preload("Books.Reviews").Preload("Books.Tags").
		Where("email = ?", email).First(user).Error
}
