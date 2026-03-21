package usecase

import (
	"context"

	"github.com/TmzFranck/books-api-golang/internal/delivery/http/middleware"
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"github.com/TmzFranck/books-api-golang/internal/model/converter"
	"github.com/TmzFranck/books-api-golang/internal/repository"
	"github.com/TmzFranck/books-api-golang/internal/utils"
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
		return nil, utils.ErrInternalServerError
	}

	return converter.BooksToResponse(books), nil

}

func (c *BookUseCase) GetBook(cx context.Context, bookId uint) (*model.BookResponse, error) {
	tx := c.DB.WithContext(cx)

	book := &entity.Book{}
	if err := c.BookRepository.FindByIdWith(tx, book, bookId, "User", "Reviews", "Tags"); err != nil {
		c.Log.Errorf("Error fetching book: %+v", err)
		return nil, utils.ErrNotFound
	}

	return converter.BookToResponse(book), nil
}

func (c *BookUseCase) CreateBook(cx context.Context, request *model.BookCreateRequest) (*model.BookResponse, error) {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, err
	}

	userId := middleware.GetUserID(cx)

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
		return nil, utils.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	db := c.DB.WithContext(cx)
	if err := c.BookRepository.FindByIdWith(db, book, book.ID, "User", "Reviews", "Tags"); err != nil {
		c.Log.Errorf("Error fetching created book: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	return converter.BookToResponse(book), nil
}

func (c *BookUseCase) UpdateBook(cx context.Context, bookId uint, request *model.BookUpdateRequest) (*model.BookResponse, error) {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, utils.ErrBadRequest
	}

	book := &entity.Book{}
	if err := c.BookRepository.FindById(tx, book, bookId); err != nil {
		c.Log.Errorf("Error fetching book: %+v", err)
		return nil, utils.ErrNotFound
	}

	userId := middleware.GetUserID(cx)
	if book.UserID != userId {
		c.Log.Errorf("User %d is not authorized to update book %d", userId, bookId)
		return nil, utils.ErrForbidden
	}
	
	book.Title = request.Title
	book.Author = request.Author
	book.Publisher = request.Publisher
	book.Language = request.Language

	if err := c.BookRepository.Update(tx, book); err != nil {
		c.Log.Errorf("Error updating book: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	db := c.DB.WithContext(cx)
	if err := c.BookRepository.FindByIdWith(db, book, book.ID, "User", "Reviews", "Tags"); err != nil {
		c.Log.Errorf("Error fetching updated book: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	return converter.BookToResponse(book), nil
}

func (c *BookUseCase) DeleteBook(cx context.Context, bookId uint) error {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	book := &entity.Book{}
	if err := c.BookRepository.FindById(tx, book, bookId); err != nil {
		c.Log.Errorf("Error fetching book: %+v", err)
		return utils.ErrNotFound
	}

	if err := c.BookRepository.Delete(tx, book); err != nil {
		c.Log.Errorf("Error deleting book: %+v", err)
		return utils.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return utils.ErrInternalServerError
	}
	return nil
}

func (c *BookUseCase) GetUserBooksSubmissions(cx context.Context, userId uint) ([]model.BookResponse, error) {
	tx := c.DB.WithContext(cx)

	books, err := c.BookRepository.GetUserBook(tx, userId)
	if err != nil {
		c.Log.Errorf("Error fetching user books: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	return converter.BooksToResponse(books), nil
}
