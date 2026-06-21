package repositories

import (
	"context"
	"database/sql"
	"project-1/internal/models"
)

type EnrollmentRepository struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) *EnrollmentRepository {
	return &EnrollmentRepository{
		db: db,
	}
}

func (r *EnrollmentRepository) ExistsEnrollmentTx(
	ctx context.Context,
	tx *sql.Tx,
	studentID string,
	courseID string,
) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM enrollments
			WHERE student_id = $1
			AND course_id = $2
		)
	`

	var isPresent bool

	err := tx.QueryRowContext(
		ctx,
		query,
		studentID,
		courseID,
	).Scan(
		&isPresent,
	)

	if err != nil {
		return false, err
	}

	return isPresent, err
}

func (r *EnrollmentRepository) CreateEnrollmentTx(
	ctx context.Context,
	tx *sql.Tx,
	enrollment *models.Enrollment,
) error {

	query := `
		INSERT INTO enrollments (
			id,
			student_id,
			course_id,
			created_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		enrollment.ID,
		enrollment.StudentID,
		enrollment.CourseID,
		enrollment.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
