package repository

import (
	"errors"
	"log"

	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	Repository[entity.Review]
	Log *log.Logger
}

func NewReviewRepository(log *log.Logger) *ReviewRepository {
	return &ReviewRepository{
		Log: log,
	}
}

func (r *ReviewRepository) AddReviewToBook(db *gorm.DB, userId, bookId uint, reviewData model.ReviewCreateRequest) error {
	bookRepo := NewBookRepository(r.Log)
	book := &entity.Book{}

	err := bookRepo.FindById(db, book, bookId)
	if err != nil {
		return err
	}

	userRepo := NewUserRepository(r.Log)
	user := &entity.User{}

	err = userRepo.FindById(db, user, userId)
	if err != nil {
		return err
	}

	newReview := &entity.Review{
		Rating:     reviewData.Rating,
		ReviewText: reviewData.ReviewText,
		BookID:     bookId,
		UserID:     userId,
	}

	return r.Create(db, newReview)
}

func (r *ReviewRepository) DeleteReviewFromBook(db *gorm.DB, reviewId, userId uint) error {
	review := &entity.Review{}

	err := r.FindById(db, review, reviewId)
	if err != nil {
		return err
	}

	if review.UserID != userId {
		return errors.New("you can only delete your own reviews")
	}

	return r.Delete(db, review)
}
