package model

import "time"

type TagResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type TagCreateRequest struct {
	Name string `json:"name" validate:"required"`
}

type TagAddRequest struct {
	Tags []TagCreateRequest `json:"tags" validate:"required"`
}
