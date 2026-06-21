package models

import "time"

type Enrollment struct {
	ID        string
	StudentID string
	CourseID  string
	CreatedAt time.Time
}
