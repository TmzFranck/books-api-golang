package http

import (
	"encoding/json"
	"net/http"

	"github.com/TmzFranck/books-api-golang/internal/model"
	"github.com/TmzFranck/books-api-golang/internal/usecase"
	"github.com/TmzFranck/books-api-golang/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type BookController struct {
	Log       *logrus.Logger
	useCase   *usecase.BookUseCase
	validator *validator.Validate
}

func NewBookController(log *logrus.Logger, useCase *usecase.BookUseCase, validator *validator.Validate) *BookController {
	return &BookController{
		Log:       log.WithField("module", "BookController").Logger,
		useCase:   useCase,
		validator: validator,
	}
}

func (c *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	c.Log.Info("fetching books")
	books, err := c.useCase.GetAllBooks(r.Context())
	if err != nil {
		c.Log.Errorf("failed to get books: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to get books")
		return
	}

	c.Log.Info("fetched books successfully")
	utils.ResponseWithData(w, http.StatusOK, books)
}

func (c *BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	BookID, err := utils.StringToUint(chi.URLParam(r, "book_id"))
	if err != nil {
		c.Log.Warnf("invalid book id: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid book id", err.Error())
		return
	}

	c.Log.Infof("fetching book %d", BookID)
	book, err := c.useCase.GetBook(r.Context(), BookID)
	if err != nil {
		c.Log.Errorf("failed to get book: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to get book")
		return
	}

	c.Log.Infof("fetched book %d successfully", BookID)
	utils.ResponseWithData(w, http.StatusOK, book)
}

func (c *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	BookId, err := utils.StringToUint(chi.URLParam(r, "book_id"))
	if err != nil {
		c.Log.Warnf("invalid book id: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid book id", err.Error())
		return
	}

	request := new(model.BookCreateRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("invalid request body: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("invalid request body: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	c.Log.Infof("creating book %d", BookId)
	book, err := c.useCase.CreateBook(r.Context(), request)
	if err != nil {
		c.Log.Errorf("failed to create book: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to create book")
		return
	}

	c.Log.Infof("created book %d successfully", BookId)
	utils.ResponseWithData(w, http.StatusOK, book)

}

func (c *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	BookId, err := utils.StringToUint(chi.URLParam(r, "book_id"))
	if err != nil {
		c.Log.Warnf("invalid book id: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid book id", err.Error())
		return
	}

	request := new(model.BookUpdateRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("invalid request body: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("invalid request body: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	c.Log.Infof("updating book %d", BookId)
	book, err := c.useCase.UpdateBook(r.Context(), BookId, request)
	if err != nil {
		c.Log.Errorf("failed to update book: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to update book")
		return
	}

	c.Log.Infof("updated book %d successfully", BookId)
	utils.ResponseWithData(w, http.StatusOK, book)
}

func (c *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	BookId, err := utils.StringToUint(chi.URLParam(r, "book_id"))
	if err != nil {
		c.Log.Warnf("invalid book id: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid book id", err.Error())
		return
	}

	c.Log.Infof("deleting book %d", BookId)
	if err := c.useCase.DeleteBook(r.Context(), BookId); err != nil {
		c.Log.Errorf("failed to delete book: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to delete book")
		return
	}

	c.Log.Infof("deleted book %d successfully", BookId)
	utils.ResponseWithData(w, http.StatusOK, nil)
}

func (c *BookController) GetUserBooks(w http.ResponseWriter, r *http.Request) {
	c.Log.Infof("getting user books")
	books, err := c.useCase.GetUserBooksSubmissions(r.Context())
	if err != nil {
		c.Log.Errorf("failed to get user books: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to get user books")
		return
	}

	utils.ResponseWithData(w, http.StatusOK, books)
}
