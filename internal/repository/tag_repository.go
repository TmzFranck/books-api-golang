package repository

import (
	"github.com/TmzFranck/books-api-golang/internal/entity"
	"github.com/TmzFranck/books-api-golang/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TagRepository struct {
	Repository[entity.Tag]
	Log *logrus.Logger
}

func NewTagRepository(log *logrus.Logger) *TagRepository {
	return &TagRepository{
		Log: log,
	}
}

func (t *TagRepository) GetAll(db *gorm.DB) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := db.Find(&tags).Error
	return tags, err
}

func (t *TagRepository) AddTagToBook(db *gorm.DB, bookId uint, tagData *model.TagAddRequest) (*entity.Book, error) {
	book := &entity.Book{}
	err := db.Preload("User").
		Preload("Reviews").
		Preload("Tags").
		Where("id = ?", bookId).
		First(book).Error
	if err != nil {
		return nil, err
	}

	var tags []entity.Tag
	for _, data := range tagData.Tags {
		tag := entity.Tag{}
		err = db.Where("name = ?", data.Name).Take(&tag).Error
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	err = db.Model(book).Association("Tags").Append(tags)
	if err != nil {
		return nil, err
	}

	err = db.Preload("User").Preload("Reviews").Preload("Tags").Where("id = ?", bookId).First(book).Error
	if err != nil {
		return nil, err
	}
	return book, nil
}
