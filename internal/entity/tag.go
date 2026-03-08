package entity

import "time"

type Tag struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;uniqueIndex"`
	Books     []Book    `gorm:"many2many:book_tags;"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (t *Tag) TableName() string { return "tags" }
