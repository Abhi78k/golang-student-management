package dto

import "time"

type CreateCourseRequest struct {
	Name           string `json:"name" binding:"required"`
	AvailableSeats int    `json:"available_seats" binding:"required"`
}

type UpdateCourseRequest struct {
	Name           string `json:"name" binding:"required"`
	AvailableSeats int    `json:"available_seats" binding:"required"`
}

type CourseFilter struct {
	Page  int
	Limit int
	Seats int
}

type CourseResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	AvailableSeats int       `json:"available_seats"`
	CreatedAt      time.Time `json:"created_at"`
}
