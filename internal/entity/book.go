package entity

import "time"

type Book struct {
	ID            uint      `gorm:"column:id;primaryKey"`
	Title         string    `gorm:"column:title"`
	Author        string    `gorm:"column:author"`
	Publisher     string    `gorm:"column:publisher"`
	PublisherDate uint      `gorm:"column:publisher_date"`
	PageCount     uint      `gorm:"column:page_count"`
	Language      string    `gorm:"column:language"`
	UserID        uint      `gorm:"column:user_id"`
	User          User      `gorm:"foreignKey:UserID"`
	Reviews       []Review  `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags          []Tag     `gorm:"many2many:book_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (b *Book) TableName() string {
	return "books"
}
