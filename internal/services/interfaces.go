package services

import (
	"context"
	"database/sql"
	"project-1/internal/dto"
	"project-1/internal/models"
)

type UserRepository interface {
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*models.User, error)

	CreateUser(
		ctx context.Context,
		user *models.User,
	) error

	GetUserByID(
		ctx context.Context,
		userID string,
	) (*models.User, error)
}

type StudentRepository interface {
	CreateStudent(
		ctx context.Context,
		student *models.Student,
	) error

	GetStudentByID(
		ctx context.Context,
		studentID string,
	) (*models.Student, error)

	GetStudentByEmail(
		ctx context.Context,
		email string,
	) (*models.Student, error)

	ListStudents(
		ctx context.Context,
		filter dto.StudentFilter,
		offset int,
	) ([]models.Student, error)

	UpdateStudent(
		ctx context.Context,
		student *models.Student,
	) error

	DeleteStudent(
		ctx context.Context,
		id string,
	) error
}

type CourseRepository interface {
	CreateCourse(
		ctx context.Context,
		course *models.Course,
	) error

	GetCourseByID(
		ctx context.Context,
		id string,
	) (*models.Course, error)

	GetCourseByName(
		ctx context.Context,
		name string,
	) (*models.Course, error)

	ListCourses(
		ctx context.Context,
		limit int,
		offset int,
	) ([]models.Course, error)

	UpdateCourse(
		ctx context.Context,
		course *models.Course,
	) error

	DeleteCourse(
		ctx context.Context,
		id string,
	) error

	GetCourseByIDTx(
		ctx context.Context,
		tx *sql.Tx,
		id string,
	) (*models.Course, error)

	CreateCourseTx(
		ctx context.Context,
		tx *sql.Tx,
		course *models.Course,
	) error

	GetCourseByNameTx(
		ctx context.Context,
		tx *sql.Tx,
		name string,
	) (*models.Course, error)

	ListCoursesTx(
		ctx context.Context,
		tx *sql.Tx,
		limit int,
		offset int,
	) ([]models.Course, error)

	UpdateCourseTx(
		ctx context.Context,
		tx *sql.Tx,
		course *models.Course,
	) error

	GetCourseByIDForUpdateTx(
		ctx context.Context,
		tx *sql.Tx,
		id string,
	) (*models.Course, error)
}

type EnrollmentRepository interface {
	ExistsEnrollmentTx(
		ctx context.Context,
		tx *sql.Tx,
		studentID string,
		courseID string,
	) (bool, error)

	CreateEnrollmentTx(
		ctx context.Context,
		tx *sql.Tx,
		enrollment *models.Enrollment,
	) error
}

type AuthCache interface {
	StoreRefreshToken(
		ctx context.Context,
		userID string,
		token string,
	) error

	GetRefreshToken(
		ctx context.Context,
		userID string,
	) (string, error)

	DeleteRefreshToken(
		ctx context.Context,
		userID string,
	) error
}

type StudentCache interface {
	GetStudent(
		ctx context.Context,
		id string,
	) (*models.Student, error)

	SetStudent(
		ctx context.Context,
		student *models.Student,
	) error

	DeleteStudent(
		ctx context.Context,
		id string,
	) error
}

type EnrollmentCourseRepository interface {
	GetCourseByIDForUpdateTx(
		ctx context.Context,
		tx *sql.Tx,
		id string,
	) (*models.Course, error)

	UpdateCourseTx(
		ctx context.Context,
		tx *sql.Tx,
		course *models.Course,
	) error
}
