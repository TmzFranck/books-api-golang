package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[T any] struct{}

func (*Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (*Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Omit(clause.Associations).Save(entity).Error
}

func (*Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (*Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (*Repository[T]) FindById(db *gorm.DB, entity *T, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (*Repository[T]) FindByIdWith(db *gorm.DB, entity *T, id any, preloads ...string) error {
	query := db.Where("id = ?", id)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query.Take(entity).Error
}
