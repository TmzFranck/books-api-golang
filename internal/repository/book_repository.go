package repository

import (
	"log"

	"github.com/TmzFranck/books-api-golang/internal/entity"
	"gorm.io/gorm"
)

type BookRepository struct {
	Repository[entity.Book]
	Log *log.Logger
}

func NewBookRepository(log *log.Logger) *BookRepository {
	return &BookRepository{
		Log: log,
	}
}

func (b *BookRepository) FindAll(db *gorm.DB, books *[]entity.Book) error {
	return db.Find(books).Error
}

func (b *BookRepository) GetUserBook(db *gorm.DB, userId uint, books *[]entity.Book) error {
	return db.Where("user_id = ?", userId).Find(books).Error
}
