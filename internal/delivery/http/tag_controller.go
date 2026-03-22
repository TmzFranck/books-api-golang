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

type TagController struct {
	Log       *logrus.Logger
	useCase   *usecase.TagUseCase
	validator *validator.Validate
}

func NewTagController(log *logrus.Logger, useCase *usecase.TagUseCase, validator *validator.Validate) *TagController {
	return &TagController{
		Log:       log,
		useCase:   useCase,
		validator: validator,
	}
}

func (c *TagController) GetTags(w http.ResponseWriter, r *http.Request) {
	c.Log.Infof("retrieving all tags")
	response, err := c.useCase.GetAllTags(r.Context())
	if err != nil {
		c.Log.Errorf("failed to retrieve tags: %+v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to retrieve tags")
		return
	}

	c.Log.Infof("retrieved %d tags", len(response))
	utils.ResponseWithData(w, http.StatusOK, response)
}

func (c *TagController) CreateTag(w http.ResponseWriter, r *http.Request) {
	request := new(model.TagCreateRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("invalid request payload: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("invalid request payload: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format", err.Error())
		return
	}

	c.Log.Infof("creating tag: %+v", request)
	response, err := c.useCase.CreateTag(r.Context(), request)
	if err != nil {
		c.Log.Errorf("failed to create tag: %+v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to create tag")
		return
	}

	utils.ResponseWithData(w, http.StatusCreated, response)
}

func (c *TagController) AddTagToBook(w http.ResponseWriter, r *http.Request) {
	request := new(model.TagAddRequest)
	BookID, err := utils.StringToUint(chi.URLParam(r, "book_id"))
	if err != nil {
		c.Log.Warnf("invalid book id: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid book id")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("invalid request payload: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("invalid request payload: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format", err.Error())
		return
	}

	c.Log.Infof("adding tag to book: %d", BookID)
	response, err := c.useCase.AddTagToBook(r.Context(), BookID, request)
	if err != nil {
		c.Log.Errorf("failed to add tag to book: %+v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to add tag to book")
		return
	}

	utils.ResponseWithData(w, http.StatusOK, response)
}

func (c *TagController) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tagID, err := utils.StringToUint(chi.URLParam(r, "tag_id"))
	if err != nil {
		c.Log.Warnf("invalid tag id: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid tag id")
		return
	}

	if err := c.useCase.DeleteTag(r.Context(), tagID); err != nil {
		c.Log.Errorf("failed to delete tag: %+v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to delete tag")
		return
	}

	utils.ResponseWithData(w, http.StatusOK, map[string]string{"message": "tag deleted successfully"})
}

func (c *TagController) UpdateTag(w http.ResponseWriter, r *http.Request) {
	tagID, err := utils.StringToUint(chi.URLParam(r, "tag_id"))
	if err != nil {
		c.Log.Warnf("invalid tag id: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid tag id")
		return
	}

	request := new(model.TagCreateRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		c.Log.Warnf("invalid request format: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	if err := c.validator.Struct(request); err != nil {
		c.Log.Warnf("invalid request: %+v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	response, err := c.useCase.UpdateTag(r.Context(), tagID, request)
	if err != nil {
		c.Log.Errorf("failed to update tag: %+v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to update tag")
		return
	}

	utils.ResponseWithData(w, http.StatusOK, response)
}
