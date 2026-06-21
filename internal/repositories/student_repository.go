package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"project-1/internal/apperrors"
	"project-1/internal/dto"
	"project-1/internal/models"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{
		db: db,
	}
}

func (r *StudentRepository) CreateStudent(
	ctx context.Context,
	student *models.Student,
) error {

	query := `
		INSERT INTO students (
		id,
		first_name,
		last_name,
		email,
		age,
		created_at
		)
		VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		student.ID,
		student.FirstName,
		student.LastName,
		student.Email,
		student.Age,
		student.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *StudentRepository) GetStudentByID(
	ctx context.Context,
	studentID string,
) (*models.Student, error) {

	query := `
		SELECT
		id,
		first_name,
		last_name,
		email,
		age,
		created_at
		FROM students
		WHERE id = $1
	`

	var student models.Student

	err := r.db.QueryRowContext(
		ctx,
		query,
		studentID,
	).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Age,
		&student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *StudentRepository) GetStudentByEmail(
	ctx context.Context,
	email string,
) (*models.Student, error) {

	query := `
		SELECT
		id,
		first_name,
		last_name,
		email,
		age,
		created_at
		FROM students
		WHERE email = $1
	`

	var student models.Student

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Age,
		&student.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *StudentRepository) ListStudents(
	ctx context.Context,
	filter dto.StudentFilter,
	offset int,
) ([]models.Student, error) {

	query := `
		SELECT
			id,
			first_name,
			last_name,
			email,
			age,
			created_at
		FROM students
		WHERE 1=1
	`

	args := []any{}

	placeholder := 1

	if filter.Age > 0 {
		query += fmt.Sprintf(
			" AND age = $%d",
			placeholder,
		)
		args = append(args, filter.Age)

		placeholder++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(
			`
				AND (
					first_name ILIKE $%d
					last_name ILIKE $%d
					OR email LIKE $%d
				)
			`,
			placeholder,
			placeholder,
			placeholder,
		)

		args = append(
			args,
			"%"+filter.Search+"%",
		)

		placeholder++
	}

	allowedSorts := map[string]string{
		"age":         "age ASC",
		"-age":        "age DESC",
		"created_at":  "created_at ASC",
		"-created_at": "created_at DESC",
	}

	orderBy := "created_at DESC"

	if value, ok := allowedSorts[filter.Sort]; ok {
		orderBy = value
	}

	query += fmt.Sprintf(
		" ORDER BY %s",
		orderBy,
	)

	query += fmt.Sprintf(
		" ORDER BY created_at DESC LIMIT %d OFFSET %d",
		placeholder,
		placeholder+1,
	)

	args = append(
		args,
		filter.Limit,
		offset,
	)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student

		err := rows.Scan(
			&student.ID,
			&student.FirstName,
			&student.LastName,
			&student.Email,
			&student.Age,
			&student.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) UpdateStudent(
	ctx context.Context,
	student *models.Student,
) error {

	query := `
		UPDATE students
		SET
			first_name = $1,
			last_name = $2,
			email = $3,
			age = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		student.FirstName,
		student.LastName,
		student.Email,
		student.Age,
		student.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrStudentNotFound
	}

	return nil
}

func (r *StudentRepository) DeleteStudent(
	ctx context.Context,
	id string,
) error {

	query := `
			DELETE FROM students
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
		return apperrors.ErrStudentNotFound
	}

	return nil
}

func (r *StudentRepository) CountStudents(
	ctx context.Context,
	filter dto.StudentFilter,
) (int, error) {

	query := `
		SELECT COUNT(*)
		FROM students
		WHERE 1=1
	`

	args := []any{}
	placeholder := 1

	if filter.Age > 0 {

		query += fmt.Sprintf(
			" AND age = $%d",
			placeholder,
		)

		args = append(
			args,
			filter.Age,
		)

		placeholder++
	}

	if filter.Search != "" {

		query += fmt.Sprintf(
			`
			AND (
				first_name ILIKE $%d
				OR last_name ILIKE $%d
				OR email ILIKE $%d
			)
			`,
			placeholder,
			placeholder,
			placeholder,
		)

		args = append(
			args,
			"%"+filter.Search+"%",
		)

		placeholder++
	}

	var total int

	err := r.db.QueryRowContext(
		ctx,
		query,
		args...,
	).Scan(
		&total,
	)

	if err != nil {
		return 0, err
	}

	return total, nil
}
