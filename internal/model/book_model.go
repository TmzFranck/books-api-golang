package model

import "time"

type BookResponse struct {
	ID            uint             `json:"id"`
	Title         string           `json:"title"`
	Author        string           `json:"author"`
	Publisher     string           `json:"publisher"`
	PublisherDate uint             `json:"publisher_date"`
	PageCount     uint             `json:"page_count"`
	Language      string           `json:"language"`
	Reviews       []ReviewResponse `json:"reviews,omitempty"`
	Tags          []TagResponse    `json:"tags,omitempty"`
	User          UserResponse     `json:"user"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}
type BookUpdateRequest struct {
	Title     string `json:"title" validate:"required"`
	Author    string `json:"author" validate:"required"`
	Publisher string `json:"publisher" validate:"required"`
	Language  string `json:"language" validate:"required"`
}

type BookCreateRequest struct {
	Title         string `json:"title" validate:"required"`
	Author        string `json:"author" validate:"required"`
	Publisher     string `json:"publisher" validate:"required"`
	PublishedDate uint   `json:"published_date" validate:"required"`
	Language      string `json:"language" validate:"required"`
	PageCount     uint   `json:"page_count" validate:"required"`
}
