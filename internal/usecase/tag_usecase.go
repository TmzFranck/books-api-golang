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

type TagUseCase struct {
	DB            *gorm.DB
	Log           *logrus.Logger
	Validate      *validator.Validate
	TagRepository *repository.TagRepository
}

func NewTagUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, tagRepository *repository.TagRepository) *TagUseCase {
	return &TagUseCase{
		DB:            db,
		Log:           logger,
		Validate:      validate,
		TagRepository: tagRepository,
	}
}

func (c *TagUseCase) GetAllTags(ctx context.Context) ([]model.TagResponse, error) {
	tx := c.DB.WithContext(ctx)

	tags, err := c.TagRepository.GetAll(tx)
	if err != nil {
		c.Log.Errorf("Error fetching tags: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	return converter.TagsToResponse(tags), nil
}

func (c *TagUseCase) CreateTag(cx context.Context, request *model.TagCreateRequest) (*model.TagResponse, error) {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, utils.ErrBadRequest
	}

	tag := &entity.Tag{
		Name: request.Name,
	}
	if err := c.TagRepository.Create(tx, tag); err != nil {
		c.Log.Errorf("Error creating tag: %+v", err)
		return nil, utils.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return nil, utils.ErrInternalServerError
	}
	return converter.TagToResponse(tag), nil
}

func (c *TagUseCase) AddTagToBook(cx context.Context, bookId uint, request *model.TagAddRequest) (*model.BookResponse, error) {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, utils.ErrBadRequest
	}

	book, err := c.TagRepository.AddTagToBook(tx, bookId, request)
	if err != nil {
		c.Log.Errorf("Error adding tag to book: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return nil, utils.ErrInternalServerError
	}

	return converter.BookToResponse(book), nil
}

func (c *TagUseCase) UpdateTag(cx context.Context, tagId uint, request *model.TagCreateRequest) (*model.TagResponse, error) {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Errorf("Validation error: %+v", err)
		return nil, utils.ErrBadRequest
	}
	tag := new(entity.Tag)
	if err := c.TagRepository.FindById(tx, tag, tagId); err != nil {
		c.Log.Errorf("Error fetching tag: %+v", err)
		return nil, utils.ErrNotFound
	}

	tag.Name = request.Name

	if err := c.TagRepository.Update(tx, tag); err != nil {
		c.Log.Errorf("Error updating tag: %+v", err)
		return nil, utils.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return nil, utils.ErrInternalServerError
	}
	return converter.TagToResponse(tag), nil
}

func (c *TagUseCase) DeleteTag(cx context.Context, tagId uint) error {
	tx := c.DB.WithContext(cx).Begin()
	defer tx.Rollback()

	tag := &entity.Tag{
		ID: tagId,
	}
	if err := c.TagRepository.Delete(tx, tag); err != nil {
		c.Log.Errorf("Error deleting tag: %+v", err)
		return utils.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Errorf("Error committing transaction: %+v", err)
		return utils.ErrInternalServerError
	}
	return nil
}
