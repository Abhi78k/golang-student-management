package apperrors

import "errors"

var ErrInvalidToken = errors.New("invalid token.")

var ErrInvalidCredentials = errors.New("invalid credentials.")

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrStudentNotFound = errors.New("student not found")

var ErrStudentAlreadyExists = errors.New("student already exists")

var ErrCourseAlreadyExists = errors.New("course already exists.")

var ErrCourseNotFound = errors.New("course not found.")

var ErrNoSeatsAvailable = errors.New("no seats available.")

var ErrAlreadyEnrolled = errors.New("course already enrolled.")
