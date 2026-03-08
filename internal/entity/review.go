package entity

import "time"

type Review struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	Rating     uint      `gorm:"column:rating"`
	ReviewText string    `gorm:"column:review_text"`
	UserID     uint      `gorm:"column:user_id"`
	BookID     uint      `gorm:"column:book_id"`
	User       User      `gorm:"foreignKey:UserID"`
	Book       Book      `gorm:"foreignKey:BookID"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (r *Review) TableName() string { return "reviews" }
