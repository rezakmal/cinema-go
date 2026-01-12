package models

import "time"

type Cinema struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCinemaRequest struct {
	Name     string  `json:"name" binding:"required"`
	Location string  `json:"location" binding:"required"`
	Rating   float64 `json:"rating"`
}

type UpdateCinemaRequest struct {
	Name     *string  `json:"name,omitempty"`
	Location *string  `json:"location,omitempty"`
	Rating   *float64 `json:"rating,omitempty"`
}
