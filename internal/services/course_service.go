package services

import (
	"context"
	"database/sql"
	"errors"
	"project-1/internal/apperrors"
	"project-1/internal/dto"
	"project-1/internal/models"
	"time"

	"github.com/google/uuid"
)

type CourseService struct {
	courseRepo CourseRepository
}

func NewCourseService(courseRepo CourseRepository) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
	}
}

func (s *CourseService) CreateCourse(
	ctx context.Context,
	input *dto.CreateCourseRequest,
) error {

	_, err := s.courseRepo.GetCourseByName(ctx, input.Name)

	if err == nil {
		return apperrors.ErrCourseAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	course := models.Course{
		ID:             uuid.NewString(),
		Name:           input.Name,
		AvailableSeats: input.AvailableSeats,
		CreatedAt:      time.Now(),
	}

	return s.courseRepo.CreateCourse(ctx, &course)
}

func (s *CourseService) GetCourseByID(
	ctx context.Context,
	id string,
) (*models.Course, error) {

	return s.courseRepo.GetCourseByID(ctx, id)
}

func (s *CourseService) ListCourses(
	ctx context.Context,
	page int,
	limit int,
) ([]models.Course, error) {

	offset := (page - 1) * limit

	return s.courseRepo.ListCourses(ctx, limit, offset)
}

func (s *CourseService) UpdateCourse(
	ctx context.Context,
	id string,
	name string,
	availableSeats int,
) error {

	course, err := s.GetCourseByID(ctx, id)

	if err != nil {
		return err
	}

	if course.Name == name {
		return apperrors.ErrCourseAlreadyExists
	}

	course.Name = name
	course.AvailableSeats = availableSeats

	return s.courseRepo.UpdateCourse(ctx, course)
}

func (s *CourseService) DeleteCourse(
	ctx context.Context,
	id string,
) error {
	return s.courseRepo.DeleteCourse(ctx, id)
}
