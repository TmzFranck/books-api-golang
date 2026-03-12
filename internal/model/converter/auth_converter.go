package converter

import (
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
)

func UserToInfoResponse(user *entity.User) *model.UserInfoResponse {
	return &model.UserInfoResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}
