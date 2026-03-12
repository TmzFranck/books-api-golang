package converter

import (
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
)

func TagToResponse(tag *entity.Tag) *model.TagResponse {
	return &model.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
	}
}

func TagsToResponse(tags []entity.Tag) []model.TagResponse {
	result := make([]model.TagResponse, 0, len(tags))
	for _, t := range tags {
		result = append(result, *TagToResponse(&t))
	}
	return result
}
