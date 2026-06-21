package services

import (
	"context"
	"database/sql"
	"errors"
	"project-1/internal/apperrors"
	"project-1/internal/auth"
	"project-1/internal/cache"
	"project-1/internal/config"
	"project-1/internal/models"
	"testing"
)

type MockUserRepository struct {
	User *models.User
	Err  error

	CreateUserCalled bool
	CreateUserErr    error
}

func (m *MockUserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	return m.User, m.Err
}

func (m *MockUserRepository) CreateUser(
	ctx context.Context,
	user *models.User,
) error {
	m.CreateUserCalled = true

	return m.CreateUserErr
}

func (m *MockUserRepository) GetUserByID(
	ctx context.Context,
	userID string,
) (*models.User, error) {

	return m.User, m.Err
}

func TestRegister(
	t *testing.T,
) {

	tests := []struct {
		name        string
		user        *models.User
		err         error
		expectError error
	}{
		{
			name: "user already exists",
			user: &models.User{
				ID: "123",
			},
			err:         nil,
			expectError: apperrors.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {

		t.Run(
			tt.name,
			func(t *testing.T) {
				repo := &MockUserRepository{
					User: tt.user,
					Err:  tt.err,
				}

				redisClient := cache.ConnectRedis()

				authCache := cache.NewAuthCache(redisClient)

				cfg := &config.Config{}

				service := NewAuthService(
					repo,
					cfg,
					authCache,
				)

				err := service.Register(
					context.Background(),
					"abhi",
					"abhi@test.com",
					"password123",
				)

				if !errors.Is(
					err,
					tt.expectError,
				) {
					t.Fatalf(
						"expected %v got %v",
						tt.expectError,
						err,
					)
				}
			},
		)
	}
}

func TestLogin(
	t *testing.T,
) {

	hashedPassword, _ := auth.HashPassword(
		"correct-password",
	)

	tests := []struct {
		name          string
		user          *models.User
		repoErr       error
		password      string
		expectedError error
	}{
		{
			name: "success",
			user: &models.User{
				ID:           "123",
				PasswordHash: hashedPassword,
			},
			repoErr:       nil,
			password:      "correct-password",
			expectedError: nil,
		},
		{
			name: "wrong password",
			user: &models.User{
				ID:           "123",
				PasswordHash: hashedPassword,
			},
			repoErr:       nil,
			password:      "wrong-password",
			expectedError: apperrors.ErrInvalidCredentials,
		},
		{
			name: "user not found",
			user: &models.User{
				ID: "123",
			},
			repoErr:       sql.ErrNoRows,
			password:      "anything",
			expectedError: apperrors.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {

		t.Run(
			tt.name,
			func(t *testing.T) {
				repo := &MockUserRepository{
					User: tt.user,
					Err:  tt.repoErr,
				}

				redisClient := cache.ConnectRedis()

				authCache := cache.NewAuthCache(redisClient)

				cfg := &config.Config{
					JWTAccessSecret:  "access-secret",
					JWTRefreshSecret: "refresh-secret",
				}

				service := NewAuthService(
					repo,
					cfg,
					authCache,
				)

				accessToken, refreshToken, err := service.Login(
					context.Background(),
					"abhi@test.com",
					tt.password,
				)

				if !errors.Is(
					err,
					tt.expectedError,
				) {
					t.Fatalf(
						"expected %v got %v",
						tt.expectedError,
						err,
					)
				}

				if tt.expectedError == nil {
					if accessToken == "" {
						t.Fatalf(
							"expected access token",
						)
					}

					if refreshToken == "" {
						t.Fatalf(
							"expected refresh token",
						)
					}
				}
			},
		)
	}
}
