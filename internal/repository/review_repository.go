package repository

import (
	"errors"

	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	Repository[entity.Review]
	Log *logrus.Logger
}

func NewReviewRepository(log *logrus.Logger) *ReviewRepository {
	return &ReviewRepository{
		Log: log,
	}
}

func (r *ReviewRepository) AddReviewToBook(db *gorm.DB, userId, bookId uint, reviewData model.ReviewCreateRequest) (*entity.Review, error) {
	var bookCount int64
	if err := db.Model(&entity.Book{}).Where("id = ?", bookId).Count(&bookCount).Error; err != nil {
		return nil, err
	}

	if bookCount == 0 {
		return nil, errors.New("book not found")
	}

	var userCount int64
	if err := db.Model(&entity.User{}).Where("id = ?", userId).Count(&userCount).Error; err != nil {
		return nil, err
	}

	if userCount == 0 {
		return nil, errors.New("user not found")
	}

	newReview := &entity.Review{
		Rating:     reviewData.Rating,
		ReviewText: reviewData.ReviewText,
		BookID:     bookId,
		UserID:     userId,
	}

	if err := r.Create(db, newReview); err != nil {
		return nil, err
	}

	return newReview, nil
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

func (r *ReviewRepository) GetAll(db *gorm.DB) ([]entity.Review, error) {
	var reviews []entity.Review
	err := db.Preload("User").Preload("Book").Find(&reviews).Error
	return reviews, err
}
