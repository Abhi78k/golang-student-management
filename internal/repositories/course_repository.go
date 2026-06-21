package repositories

import (
	"context"
	"database/sql"
	"project-1/internal/apperrors"
	"project-1/internal/models"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

func (r *CourseRepository) CreateCourse(
	ctx context.Context,
	course *models.Course,
) error {

	query := `
		INSERT INTO courses (
			id,
			name,
			available_seats,
			created_at
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		course.ID,
		course.Name,
		course.AvailableSeats,
		course.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *CourseRepository) GetCourseByID(
	ctx context.Context,
	id string,
) (*models.Course, error) {

	query := `
		SELECT
			id,
			name,
			available_seats,
			created_at
		FROM courses
		WHERE id = $1
	`

	var course models.Course

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&course.ID,
		&course.Name,
		&course.AvailableSeats,
		&course.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *CourseRepository) GetCourseByName(
	ctx context.Context,
	name string,
) (*models.Course, error) {

	query := `
		SELECT
			id,
			name,
			available_seats,
			created_at
		FROM courses
		WHERE name = $1
	`

	var course models.Course

	err := r.db.QueryRowContext(
		ctx,
		query,
		name,
	).Scan(
		&course.ID,
		&course.Name,
		&course.AvailableSeats,
		&course.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *CourseRepository) ListCourses(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Course, error) {

	query := `
		SELECT
			id,
			name,
			available_seats,
			created_at
		FROM courses
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`

	var courses []models.Course

	rows, err := r.db.QueryContext(
		ctx,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var course models.Course

		err := rows.Scan(
			&course.ID,
			&course.Name,
			&course.AvailableSeats,
			&course.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, rows.Err()
	}

	return courses, nil
}

func (r *CourseRepository) UpdateCourse(
	ctx context.Context,
	course *models.Course,
) error {

	query := `
		UPDATE courses
		SET
			name = $1,
			available_seats = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		course.Name,
		course.AvailableSeats,
		course.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrCourseNotFound
	}

	return nil
}

func (r *CourseRepository) DeleteCourse(
	ctx context.Context,
	id string,
) error {

	query := `
		DELETE FROM courses
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrCourseNotFound
	}

	return nil
}

// =============================================
// =========TRANSACTION AWARE METHODS===========
// =============================================

func (r *CourseRepository) GetCourseByIDTx(
	ctx context.Context,
	tx *sql.Tx,
	id string,
) (*models.Course, error) {

	query := `
		SELECT
			id,
			name,
			available_seats,
			created_at
		FROM courses
		WHERE id = $1
	`

	var course models.Course

	err := tx.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&course.ID,
		&course.Name,
		&course.AvailableSeats,
		&course.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *CourseRepository) CreateCourseTx(
	ctx context.Context,
	tx *sql.Tx,
	course *models.Course,
) error {

	query := `
		INSERT INTO courses (
			id,
			name,
			available_seats,
			created_at
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		course.ID,
		course.Name,
		course.AvailableSeats,
		course.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *CourseRepository) GetCourseByNameTx(
	ctx context.Context,
	tx *sql.Tx,
	name string,
) (*models.Course, error) {

	query := `
		SELECT
			id,
			name,
			available_seats,
			created_at
		FROM courses
		WHERE name = $1
	`

	var course models.Course

	err := tx.QueryRowContext(
		ctx,
		query,
		name,
	).Scan(
		&course.ID,
		&course.Name,
		&course.AvailableSeats,
		&course.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *CourseRepository) ListCoursesTx(
	ctx context.Context,
	tx *sql.Tx,
	limit int,
	offset int,
) ([]models.Course, error) {

	query := `
		SELECT
			id,
			name,
			available_seats,
			created_at
		FROM courses
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`

	var courses []models.Course

	rows, err := tx.QueryContext(
		ctx,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var course models.Course

		err := rows.Scan(
			&course.ID,
			&course.Name,
			&course.AvailableSeats,
			&course.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, rows.Err()
	}

	return courses, nil
}

func (r *CourseRepository) UpdateCourseTx(
	ctx context.Context,
	tx *sql.Tx,
	course *models.Course,
) error {

	query := `
		UPDATE courses
		SET
			name = $1,
			available_seats = $2
		WHERE id = $3
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		course.Name,
		course.AvailableSeats,
		course.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrCourseNotFound
	}

	return nil
}

func (r *CourseRepository) DeleteCourseTx(
	ctx context.Context,
	tx *sql.Tx,
	id string,
) error {

	query := `
		DELETE FROM courses
		WHERE id = $1
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrCourseNotFound
	}

	return nil
}

// =============================================
// =============FOR UPDATE METHODS==============
// =============================================

func (r *CourseRepository) GetCourseByIDForUpdateTx(
	ctx context.Context,
	tx *sql.Tx,
	id string,
) (*models.Course, error) {
	query := `
		SELECT
			id,
			name,
			available_seats,
			created_at
		FROM courses
		WHERE id = $1
		FOR UPDATE
	`

	var course models.Course

	err := tx.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&course.ID,
		&course.Name,
		&course.AvailableSeats,
		&course.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &course, nil
}
