package converter

import (
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
)

func BookToResponse(book *entity.Book) *model.BookResponse {
	return &model.BookResponse{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		Publisher:     book.Publisher,
		PublisherDate: book.PublisherDate,
		PageCount:     book.PageCount,
		Language:      book.Language,
		Reviews:       ReviewsToResponse(book.Reviews),
		Tags:          TagsToResponse(book.Tags),
		User:          *UserToResponse(&book.User),
		CreatedAt:     book.CreatedAt,
		UpdatedAt:     book.UpdatedAt,
	}
}

func BooksToResponse(books []entity.Book) []model.BookResponse {
	result := make([]model.BookResponse, 0, len(books))
	for _, b := range books {
		result = append(result, *BookToResponse(&b))
	}
	return result
}
