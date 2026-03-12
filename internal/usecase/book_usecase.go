package usecase

import (
	"context"

	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"github.com/TmzFranck/books-api-golang/internal/model/converter"
	"github.com/TmzFranck/books-api-golang/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BookUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	BookRepository *repository.BookRepository
}

func NewBookUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, bookRepository *repository.BookRepository) *BookUseCase {
	return &BookUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		BookRepository: bookRepository,
	}
}

func (c *BookUseCase) GetAllBooks(cx context.Context) ([]model.BookResponse, error) {
	tx := c.DB.WithContext(cx)

	var books []entity.Book
	if err := c.BookRepository.FindAll(tx, &books); err != nil {
		c.Log.Errorf("Error fetching books: %+v", err)
		return nil, err
	}

	return converter.BooksToResponse(books), nil

}

func (c *BookUseCase) GetBook(cx context.Context, bookId uint) (*model.BookResponse, error) {
	tx := c.DB.WithContext(cx)

	book := &entity.Book{}
	if err := c.BookRepository.FindByIdWith(tx, book, bookId, "User", "Reviews", "Tags"); err != nil {
		c.Log.Errorf("Error fetching book: %+v", err)
		return nil, err
	}

	return converter.BookToResponse(book), nil
}

func (c *BookUseCase) CreateBook(cx context.Context, userId uint, request *model.BookCreateRequest) (*model.BookResponse, error) {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, err
	}

	book := &entity.Book{
		Title:         request.Title,
		Author:        request.Author,
		Publisher:     request.Publisher,
		PublisherDate: request.PublishedDate,
		PageCount:     request.PageCount,
		Language:      request.Language,
		UserID:        userId,
	}
	if err := c.BookRepository.Create(tx, book); err != nil {
		c.Log.Errorf("Error creating book: %+v", err)
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return nil, err
	}

	db := c.DB.WithContext(cx)
	if err := c.BookRepository.FindByIdWith(db, book, book.ID, "User", "Reviews", "Tags"); err != nil {
		c.Log.Errorf("Error fetching created book: %+v", err)
		return nil, err
	}

	return converter.BookToResponse(book), nil
}

func (c *BookUseCase) UpdateBook(cx context.Context, bookId uint, request *model.BookUpdateRequest) (*model.BookResponse, error) {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, err
	}

	book := &entity.Book{}
	if err := c.BookRepository.FindById(tx, book, bookId); err != nil {
		c.Log.Errorf("Error fetching book: %+v", err)
		return nil, err
	}

	book.Title = request.Title
	book.Author = request.Author
	book.Publisher = request.Publisher
	book.Language = request.Language

	if err := c.BookRepository.Update(tx, book); err != nil {
		c.Log.Errorf("Error updating book: %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return nil, err
	}

	db := c.DB.WithContext(cx)
	if err := c.BookRepository.FindByIdWith(db, book, book.ID, "User", "Reviews", "Tags"); err != nil {
		c.Log.Errorf("Error fetching updated book: %+v", err)
		return nil, err
	}

	return converter.BookToResponse(book), nil
}

func (c *BookUseCase) DeleteBook(cx context.Context, bookId uint) error {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	book := &entity.Book{}
	if err := c.BookRepository.FindById(tx, book, bookId); err != nil {
		c.Log.Errorf("Error fetching book: %+v", err)
		return err
	}

	if err := c.BookRepository.Delete(tx, book); err != nil {
		c.Log.Errorf("Error deleting book: %+v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return err
	}
	return nil
}

func (c *BookUseCase) GetUserBooksSubmissions(cx context.Context, userId uint) ([]model.BookResponse, error) {
	tx := c.DB.WithContext(cx)

	books, err := c.BookRepository.GetUserBook(tx, userId)
	if err != nil {
		c.Log.Errorf("Error fetching user books: %+v", err)
		return nil, err
	}

	return converter.BooksToResponse(books), nil
}
