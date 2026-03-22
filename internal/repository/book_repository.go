package repository

import (
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookRepository struct {
	Repository[entity.Book]
	Log *logrus.Logger
}

func NewBookRepository(log *logrus.Logger) *BookRepository {
	return &BookRepository{
		Log: log,
	}
}

func (b *BookRepository) FindAll(db *gorm.DB) ([]entity.Book, error) {
	var books []entity.Book
	if err := db.Preload("User").Preload("Reviews").Preload("Tags").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (b *BookRepository) GetUserBook(db *gorm.DB, userId uint) ([]entity.Book, error) {
	var books []entity.Book
	err := db.Where("user_id = ?", userId).
		Preload("User").
		Preload("Reviews").
		Preload("Tags").
		Find(&books).Error
	return books, err
}
