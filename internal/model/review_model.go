package model

import "time"

type ReviewResponse struct {
	ID         uint      `json:"id"`
	Rating     uint      `json:"rating"`
	ReviewText string    `json:"review_text"`
	BookId     uint      `json:"book_id"`
	UserId     uint      `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ReviewCreateRequest struct {
	Rating     uint   `json:"rating" validate:"required"`
	ReviewText string `json:"review_text" validate:"required"`
}
