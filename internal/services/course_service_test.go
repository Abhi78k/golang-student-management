package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"project-1/internal/apperrors"
	"project-1/internal/dto"
	"project-1/internal/models"
)

type MockCourseRepository struct {
	Course  *models.Course
	Courses []models.Course

	GetCourseByIDErr   error
	GetCourseByNameErr error
	CreateCourseErr    error
	UpdateCourseErr    error
	DeleteCourseErr    error
	ListCoursesErr     error

	CreateCalled bool
	UpdateCalled bool
	DeleteCalled bool
}

func (m *MockCourseRepository) CreateCourse(
	ctx context.Context,
	course *models.Course,
) error {
	m.CreateCalled = true
	return m.CreateCourseErr
}

func (m *MockCourseRepository) GetCourseByID(
	ctx context.Context,
	id string,
) (*models.Course, error) {
	return m.Course, m.GetCourseByIDErr
}

func (m *MockCourseRepository) GetCourseByName(
	ctx context.Context,
	name string,
) (*models.Course, error) {
	return m.Course, m.GetCourseByNameErr
}

func (m *MockCourseRepository) ListCourses(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Course, error) {
	return m.Courses, m.ListCoursesErr
}

func (m *MockCourseRepository) UpdateCourse(
	ctx context.Context,
	course *models.Course,
) error {
	m.UpdateCalled = true
	return m.UpdateCourseErr
}

func (m *MockCourseRepository) DeleteCourse(
	ctx context.Context,
	id string,
) error {
	m.DeleteCalled = true
	return m.DeleteCourseErr
}

func TestCreateCourse(t *testing.T) {

	tests := []struct {
		name               string
		course             *models.Course
		getCourseByNameErr error
		expectedErr        error
	}{
		{
			name: "course already exists",
			course: &models.Course{
				ID: "123",
			},
			getCourseByNameErr: nil,
			expectedErr:        apperrors.ErrCourseAlreadyExists,
		},
		{
			name:               "course created",
			course:             nil,
			getCourseByNameErr: sql.ErrNoRows,
			expectedErr:        nil,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			repo := &MockCourseRepository{
				Course:             tt.course,
				GetCourseByNameErr: tt.getCourseByNameErr,
			}

			service := NewCourseService(repo)

			err := service.CreateCourse(
				context.Background(),
				&dto.CreateCourseRequest{
					Name:           "Go Backend",
					AvailableSeats: 50,
				},
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
					"expected CreateCourse to be called",
				)
			}
		})
	}
}

func TestGetCourseByID(t *testing.T) {

	expected := &models.Course{
		ID:   "123",
		Name: "Go",
	}

	repo := &MockCourseRepository{
		Course: expected,
	}

	service := NewCourseService(repo)

	course, err := service.GetCourseByID(
		context.Background(),
		"123",
	)

	if err != nil {
		t.Fatal(err)
	}

	if course.ID != expected.ID {
		t.Fatalf(
			"expected %s got %s",
			expected.ID,
			course.ID,
		)
	}
}

func TestListCourses(t *testing.T) {

	courses := []models.Course{
		{ID: "1"},
		{ID: "2"},
	}

	repo := &MockCourseRepository{
		Courses: courses,
	}

	service := NewCourseService(repo)

	result, err := service.ListCourses(
		context.Background(),
		1,
		10,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 courses got %d",
			len(result),
		)
	}
}

func TestUpdateCourse(t *testing.T) {

	repo := &MockCourseRepository{
		Course: &models.Course{
			ID:             "123",
			Name:           "Old Name",
			AvailableSeats: 10,
		},
	}

	service := NewCourseService(repo)

	err := service.UpdateCourse(
		context.Background(),
		"123",
		"New Name",
		25,
	)

	if err != nil {
		t.Fatal(err)
	}

	if !repo.UpdateCalled {
		t.Fatal(
			"expected UpdateCourse to be called",
		)
	}
}

func TestDeleteCourse(t *testing.T) {

	repo := &MockCourseRepository{}

	service := NewCourseService(repo)

	err := service.DeleteCourse(
		context.Background(),
		"123",
	)

	if err != nil {
		t.Fatal(err)
	}

	if !repo.DeleteCalled {
		t.Fatal(
			"expected DeleteCourse to be called",
		)
	}
}
