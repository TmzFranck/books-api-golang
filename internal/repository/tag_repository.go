package repository

import (
	"log"

	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"gorm.io/gorm"
)

type TagRepository struct {
	Repository[entity.Tag]
	Log *log.Logger
}

func NewTagRepository(log *log.Logger) *TagRepository {
	return &TagRepository{
		Log: log,
	}
}

func (t *TagRepository) GetAll(db *gorm.DB, tags *[]entity.Tag) error {
	return db.Find(&tags).Error
}

func (t *TagRepository) AddTagToBook(db *gorm.DB, bookId uint, tagData model.TagAddRequest) error {
	book := &entity.Book{}
	err := db.First(book, bookId).Error
	if err != nil {
		return err
	}

	var tags []entity.Tag
	for _, data := range tagData.Tags {
		tag := entity.Tag{}
		err = db.Where("name = ?", data.Name).Take(&tag).Error
		if err != nil {
			return err
		}
		tags = append(tags, tag)
	}

	return db.Model(book).Association("Tags").Append(tags)
}
