package converter

import (
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
)

func ReviewToResponse(review *entity.Review) *model.ReviewResponse {
	return &model.ReviewResponse{
		ID:         review.ID,
		Rating:     review.Rating,
		ReviewText: review.ReviewText,
		BookId:     review.BookID,
		UserId:     review.UserID,
		CreatedAt:  review.CreatedAt,
		UpdatedAt:  review.UpdatedAt,
	}
}

func ReviewsToResponse(reviews []entity.Review) []model.ReviewResponse {
	result := make([]model.ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		result = append(result, *ReviewToResponse(&r))
	}
	return result
}
