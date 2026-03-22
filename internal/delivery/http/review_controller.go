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

type ReviewController struct {
	Log       *logrus.Logger
	useCase   *usecase.ReviewUseCase
	validator *validator.Validate
}

func NewReviewController(log *logrus.Logger, useCase *usecase.ReviewUseCase, validator *validator.Validate) *ReviewController {
	return &ReviewController{
		Log:       log,
		useCase:   useCase,
		validator: validator,
	}
}

func (c *ReviewController) GetReviews(w http.ResponseWriter, r *http.Request) {
	c.Log.Info("Fetching all reviews")
	reviews, err := c.useCase.GetAllReviews(r.Context())
	if err != nil {
		c.Log.Errorf("Error fetching reviews: %+v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "")
		return
	}

	c.Log.Info("Successfully fetched all reviews")
	utils.ResponseWithData(w, http.StatusOK, reviews)
}

func (c *ReviewController) GetReview(w http.ResponseWriter, r *http.Request) {
	ReviewID, err := utils.StringToUint(chi.URLParam(r, "review_id"))
	if err != nil {
		c.Log.Errorf("Invalid review ID: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "Invalid review ID", err.Error())
		return
	}

	c.Log.Infof("fetching info for reviewID: %d", ReviewID)
	review, err := c.useCase.GetReview(r.Context(), ReviewID)
	if err != nil {
		c.Log.Errorf("Error fetching review: %+v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "Error fetching review")
		return
	}

	c.Log.Info("Successfully fetched review")
	utils.ResponseWithData(w, http.StatusOK, review)
}

func (c *ReviewController) AddReviewToBook(w http.ResponseWriter, r *http.Request) {
	BookID, err := utils.StringToUint(chi.URLParam(r, "book_id"))
	if err != nil {
		c.Log.Warnf("Invalid book ID: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "Invalid book ID", err.Error())
		return
	}

	request := new(model.ReviewCreateRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("Invalid review data: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "Invalid review data", err.Error())
		return
	}

	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("Validation error: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	c.Log.Infof("Adding review to book ID: %d", BookID)
	response, err := c.useCase.AddReviewToBook(r.Context(), BookID, request)
	if err != nil {
		c.Log.Errorf("Error adding review: %+v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "Error adding review")
		return
	}

	c.Log.Info("Successfully added review")
	utils.ResponseWithData(w, http.StatusCreated, response)
}

func (c *ReviewController) DeleteReviewFromBook(w http.ResponseWriter, r *http.Request) {
	ReviewID, err := utils.StringToUint(chi.URLParam(r, "review_id"))
	if err != nil {
		c.Log.Warnf("Invalid review ID: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "Invalid review ID", err.Error())
		return
	}

	c.Log.Infof("Deleting review ID: %d", ReviewID)
	if err := c.useCase.DeleteReviewFromBook(r.Context(), ReviewID); err != nil {
		c.Log.Errorf("Error deleting review: %+v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "Error deleting review")
		return
	}

	c.Log.Info("Successfully deleted review")
	utils.ResponseWithData(w, http.StatusOK, map[string]string{"message": "review deleted successfully"})
}
