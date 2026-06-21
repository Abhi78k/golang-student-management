package models

import "time"

type Student struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Age       int
	CreatedAt time.Time
}
