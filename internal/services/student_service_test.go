// package services

// import (
// 	"context"
// 	"project-1/internal/models"
// )

// type MockStudentRepository struct {
// 	Student *models.Student
// 	Err     error

// 	CreateStudentCalled bool
// }

// func (m *MockStudentRepository) CreateStudent(
// 	ctx context.Context,
// 	student *models.Student,
// ) error {

// 	m.CreateStudentCalled = true

//		return nil
//	}
package services

import (
	"context"
	"database/sql"
	"errors"
	"project-1/internal/apperrors"
	"project-1/internal/cache"
	"project-1/internal/config"
	"project-1/internal/dto"
	"project-1/internal/models"
	"testing"
)

type MockStudentRepository struct {
	Student  *models.Student
	Students []models.Student

	GetStudentByIDErr    error
	GetStudentByEmailErr error
	ListStudentsErr      error
	UpdateStudentErr     error
	DeleteStudentErr     error
	CreateStudentErr     error

	CreateCalled bool
	UpdateCalled bool
	DeleteCalled bool
}

func (m *MockStudentRepository) CreateStudent(
	ctx context.Context,
	student *models.Student,
) error {
	m.CreateCalled = true
	return m.CreateStudentErr
}

func (m *MockStudentRepository) GetStudentByID(
	ctx context.Context,
	id string,
) (*models.Student, error) {
	return m.Student, m.GetStudentByIDErr
}

func (m *MockStudentRepository) GetStudentByEmail(
	ctx context.Context,
	email string,
) (*models.Student, error) {
	return m.Student, m.GetStudentByEmailErr
}

func (m *MockStudentRepository) ListStudents(
	ctx context.Context,
	filter dto.StudentFilter,
	offset int,
) ([]models.Student, error) {
	return m.Students, m.ListStudentsErr
}

func (m *MockStudentRepository) UpdateStudent(
	ctx context.Context,
	student *models.Student,
) error {
	m.UpdateCalled = true
	return m.UpdateStudentErr
}

func (m *MockStudentRepository) DeleteStudent(
	ctx context.Context,
	id string,
) error {
	m.DeleteCalled = true
	return m.DeleteStudentErr
}

func TestCreateStudent(t *testing.T) {

	tests := []struct {
		name                 string
		student              *models.Student
		getStudentByEmailErr error
		createStudentErr     error
		expectedErr          error
	}{
		{
			name: "student already exists",
			student: &models.Student{
				ID: "123",
			},
			getStudentByEmailErr: nil,
			expectedErr:          apperrors.ErrStudentAlreadyExists,
		},
		{
			name:                 "student created",
			student:              nil,
			getStudentByEmailErr: sql.ErrNoRows,
			createStudentErr:     nil,
			expectedErr:          nil,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			repo := &MockStudentRepository{
				Student:              tt.student,
				GetStudentByEmailErr: tt.getStudentByEmailErr,
				CreateStudentErr:     tt.createStudentErr,
			}

			redisClient := cache.ConnectRedis()

			studentCache := cache.NewStudentCache(redisClient)

			service := NewStudentService(
				repo,
				&config.Config{},
				studentCache,
			)

			err := service.CreateStudent(
				context.Background(),
				"Abhi",
				"K",
				"abhi@test.com",
				20,
			)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf(
					"expected %v got %v",
					tt.expectedErr,
					err,
				)
			}

			if tt.expectedErr == nil &&
				!repo.CreateCalled {
				t.Fatal(
					"expected CreateStudent to be called",
				)
			}
		})
	}
}

func TestGetStudentByID(t *testing.T) {

	expected := &models.Student{
		ID:        "123",
		FirstName: "Abhi",
	}

	repo := &MockStudentRepository{
		Student:           expected,
		GetStudentByIDErr: nil,
	}

	redisClient := cache.ConnectRedis()

	studentCache := cache.NewStudentCache(redisClient)

	service := NewStudentService(
		repo,
		&config.Config{},
		studentCache,
	)

	student, err := service.GetStudentByID(
		context.Background(),
		"123",
	)

	if err != nil {
		t.Fatal(err)
	}

	if student.ID != expected.ID {
		t.Fatalf(
			"expected %s got %s",
			expected.ID,
			student.ID,
		)
	}
}

func TestListStudents(t *testing.T) {

	students := []models.Student{
		{
			ID: "1",
		},
		{
			ID: "2",
		},
	}

	repo := &MockStudentRepository{
		Students:        students,
		ListStudentsErr: nil,
	}

	redisClient := cache.ConnectRedis()

	studentCache := cache.NewStudentCache(redisClient)

	service := NewStudentService(
		repo,
		&config.Config{},
		studentCache,
	)

	result, err := service.ListStudents(
		context.Background(),
		dto.StudentFilter{
			Page:  1,
			Limit: 10,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 students got %d",
			len(result),
		)
	}
}

func TestUpdateStudent(t *testing.T) {

	repo := &MockStudentRepository{
		Student: &models.Student{
			ID:        "123",
			FirstName: "Old",
		},
		GetStudentByIDErr: nil,
		UpdateStudentErr:  nil,
	}

	redisClient := cache.ConnectRedis()

	studentCache := cache.NewStudentCache(redisClient)

	service := NewStudentService(
		repo,
		&config.Config{},
		studentCache,
	)

	err := service.UpdateStudent(
		context.Background(),
		"123",
		"New",
		"User",
		"new@test.com",
		22,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !repo.UpdateCalled {
		t.Fatal(
			"expected UpdateStudent to be called",
		)
	}
}

func TestDeleteStudent(t *testing.T) {

	repo := &MockStudentRepository{
		DeleteStudentErr: nil,
	}

	redisClient := cache.ConnectRedis()

	studentCache := cache.NewStudentCache(redisClient)

	service := NewStudentService(
		repo,
		&config.Config{},
		studentCache,
	)

	err := service.DeleteStudent(
		context.Background(),
		"123",
	)

	if err != nil {
		t.Fatal(err)
	}

	if !repo.DeleteCalled {
		t.Fatal(
			"expected DeleteStudent to be called",
		)
	}
}
