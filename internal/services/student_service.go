package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"project-1/internal/apperrors"
	"project-1/internal/config"
	"project-1/internal/dto"
	"project-1/internal/models"
	"time"

	"github.com/google/uuid"
)

type StudentService struct {
	studentRepo StudentRepository
	cache       StudentCache
	cfg         *config.Config
}

func NewStudentService(studentRepo StudentRepository, cfg *config.Config, cache StudentCache) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
		cfg:         cfg,
		cache:       cache,
	}
}

func (s *StudentService) CreateStudent(
	ctx context.Context,
	first_name string,
	last_name string,
	email string,
	age int,
) error {

	_, err := s.studentRepo.GetStudentByEmail(ctx, email)

	if err == nil {
		return apperrors.ErrStudentAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	student := models.Student{
		ID:        uuid.NewString(),
		FirstName: first_name,
		LastName:  last_name,
		Email:     email,
		Age:       age,
		CreatedAt: time.Now(),
	}

	return s.studentRepo.CreateStudent(ctx, &student)
}

func (s *StudentService) GetStudentByID(
	ctx context.Context,
	id string,
) (*models.Student, error) {

	student, err := s.cache.GetStudent(ctx, id)

	if err == nil {
		return student, nil
	}

	log.Println("CACHE MISS")

	student, err = s.studentRepo.GetStudentByID(ctx, id)

	if err != nil {
		return nil, err
	}

	_ = s.cache.SetStudent(ctx, student)

	return student, nil
}

func (s *StudentService) ListStudents(
	ctx context.Context,
	filter dto.StudentFilter,
) ([]models.Student, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.Limit < 1 {
		filter.Limit = 10
	}

	if filter.Limit > 100 {
		filter.Limit = 100
	}

	var offset int = ((filter.Page - 1) * filter.Limit)
	return s.studentRepo.ListStudents(ctx, filter, offset)
}

func (s *StudentService) UpdateStudent(
	ctx context.Context,
	id string,
	firstName string,
	lastName string,
	email string,
	age int,
) error {

	student, err := s.studentRepo.GetStudentByID(ctx, id)

	if err != nil {
		return err
	}

	student.FirstName = firstName
	student.LastName = lastName
	student.Email = email
	student.Age = age

	err = s.studentRepo.UpdateStudent(ctx, student)

	if err != nil {
		return err
	}

	_ = s.cache.DeleteStudent(
		ctx,
		id,
	)

	return nil
}

func (s *StudentService) DeleteStudent(
	ctx context.Context,
	id string,
) error {

	err := s.studentRepo.DeleteStudent(ctx, id)

	if err != nil {
		return err
	}

	_ = s.cache.DeleteStudent(
		ctx,
		id,
	)

	return nil
}

func (s *StudentService) UpdateStudentDetail(
	ctx context.Context,
	id string,
	input dto.PatchStudentRequest,
) error {

	student, err := s.GetStudentByID(ctx, id)

	if err != nil {
		return err
	}

	if input.FirstName != nil {
		student.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		student.LastName = *input.LastName
	}

	if input.Email != nil {
		student.Email = *input.Email
	}

	if input.Age != nil && *input.Age >= 1 {
		student.Age = *input.Age
	}

	err = s.studentRepo.UpdateStudent(ctx, student)

	if err != nil {
		return err
	}

	_ = s.cache.DeleteStudent(
		ctx,
		id,
	)

	return nil
}
