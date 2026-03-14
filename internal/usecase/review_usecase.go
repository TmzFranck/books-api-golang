package usecase

import (
	"context"

	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"github.com/TmzFranck/books-api-golang/internal/model/converter"
	"github.com/TmzFranck/books-api-golang/internal/repository"
	"github.com/TmzFranck/books-api-golang/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReviewUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	ReviewRepository *repository.ReviewRepository
}

func NewReviewUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, reviewRepository *repository.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{
		DB:               db,
		Log:              log,
		Validate:         validate,
		ReviewRepository: reviewRepository,
	}
}

func (c *ReviewUseCase) GetAllReviews(cx context.Context) ([]model.ReviewResponse, error) {
	tx := c.DB.WithContext(cx)

	reviews, err := c.ReviewRepository.GetAll(tx)
	if err != nil {
		c.Log.Error("Error fetching reviews: ", err)
		return nil, utils.ErrInternalServerError
	}

	return converter.ReviewsToResponse(reviews), nil
}

func (c *ReviewUseCase) GetReview(cx context.Context, ReviewId uint) (*model.ReviewResponse, error) {
	tx := c.DB.WithContext(cx)

	review := &entity.Review{}
	if err := c.ReviewRepository.FindByIdWith(tx, review, ReviewId, "User", "Book"); err != nil {
		c.Log.Errorf("Error fetching review: %+v", err)
		return nil, utils.ErrNotFound
	}

	return converter.ReviewToResponse(review), nil
}

func (c *ReviewUseCase) DeleteReviewFromBook(cx context.Context, reviewId, userId uint) error {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	if err := c.ReviewRepository.DeleteReviewFromBook(tx, reviewId, userId); err != nil {
		c.Log.Errorf("Error deleting review: %+v", err)
		return utils.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return utils.ErrInternalServerError
	}
	return nil
}

func (c *ReviewUseCase) AddReviewToBook(cx context.Context, userId, bookId uint, request *model.ReviewCreateRequest) (*model.ReviewResponse, error) {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, utils.ErrBadRequest
	}

	review, err := c.ReviewRepository.AddReviewToBook(tx, userId, bookId, *request)
	if err != nil {
		c.Log.Errorf("Error adding review: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return nil, utils.ErrInternalServerError
	}
	return converter.ReviewToResponse(review), nil
}
