package model

import "time"

type ReviewResponse struct {
	ID         uint      `json:"id"`
	Rating     uint      `json:"rating"`
	ReviewText string    `json:"review_text"`
	CreateAt   time.Time `json:"create_at"`
	UpdateAt   time.Time `json:"update_at"`
}

type ReviewCreateRequest struct {
	Rating     uint   `json:"rating" validate:"required"`
	ReviewText string `json:"review_text" validate:"required"`
}
