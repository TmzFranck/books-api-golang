package converter

import (
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Firstname:  user.Firstname,
		Lastname:   user.Lastname,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func UserBooksToResponse(user *entity.User, books []entity.Book) *model.UserBooksResponse {
	return &model.UserBooksResponse{
		UserResponse: *UserToResponse(user),
		Books:        BooksToResponse(books),
	}
}
