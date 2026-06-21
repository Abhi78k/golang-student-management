package models

import "time"

type Course struct {
	ID             string
	Name           string
	AvailableSeats int
	CreatedAt      time.Time
}
