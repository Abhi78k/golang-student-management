package services

import (
	"context"
	"database/sql"
	"project-1/internal/apperrors"
	"project-1/internal/models"
	"time"

	"github.com/google/uuid"
)

type EnrollmentService struct {
	db             *sql.DB
	courseRepo     CourseRepository
	enrollmentRepo EnrollmentRepository
}

func NewEnrollmentService(
	db *sql.DB,
	courseRepo CourseRepository,
	enrollmentRepo EnrollmentRepository,
) *EnrollmentService {
	return &EnrollmentService{
		db:             db,
		courseRepo:     courseRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

func (s *EnrollmentService) Enroll(
	ctx context.Context,
	studentID string,
	courseID string,
) error {

	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	course, err := s.courseRepo.GetCourseByIDForUpdateTx(ctx, tx, courseID)

	if err != nil {
		return err
	}

	if course.AvailableSeats <= 0 {
		return apperrors.ErrNoSeatsAvailable
	}

	exists, err := s.enrollmentRepo.ExistsEnrollmentTx(ctx, tx, studentID, courseID)

	if err != nil {
		return err
	}

	if exists {
		return apperrors.ErrAlreadyEnrolled
	}

	enrollment := models.Enrollment{
		ID:        uuid.NewString(),
		StudentID: studentID,
		CourseID:  courseID,
		CreatedAt: time.Now(),
	}

	err = s.enrollmentRepo.CreateEnrollmentTx(ctx, tx, &enrollment)

	if err != nil {
		return err
	}

	course.AvailableSeats--

	err = s.courseRepo.UpdateCourseTx(
		ctx,
		tx,
		course,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}
