package services

import (
	"context"
	"database/sql"
	"errors"
	"project-1/internal/apperrors"
	"project-1/internal/auth"
	"project-1/internal/config"
	"project-1/internal/models"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo  UserRepository
	cfg       *config.Config
	authCache AuthCache
}

func NewAuthService(
	userRepo UserRepository,
	cfg *config.Config,
	authCache AuthCache,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		cfg:       cfg,
		authCache: authCache,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	username string,
	email string,
	password string,
) error {
	_, err := s.userRepo.GetUserByEmail(ctx, email)

	if err == nil {
		return apperrors.ErrUserAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	hashedPassword, err := auth.HashPassword(password)

	if err != nil {
		return err
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	err = s.userRepo.CreateUser(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, string, error) {

	user, err := s.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		return "", "", apperrors.ErrInvalidCredentials
	}

	err = auth.CheckPassword(
		user.PasswordHash,
		password,
	)

	if err != nil {
		return "", "", apperrors.ErrInvalidCredentials
	}

	accessToken, err := auth.GenerateAccessToken(
		user.ID,
		s.cfg,
	)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.GenerateRefreshToken(
		user.ID,
		s.cfg,
	)

	if err != nil {
		return "", "", err
	}

	err = s.authCache.StoreRefreshToken(
		ctx,
		user.ID,
		refreshToken,
	)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *AuthService) Logout(
	ctx context.Context,
	userID string,
) error {

	return s.authCache.DeleteRefreshToken(
		ctx,
		userID,
	)
}

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (string, error) {

	token, err := auth.ValidateRefreshToken(
		refreshToken,
		s.cfg,
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*auth.Claims)

	if !ok {
		return "", apperrors.ErrInvalidToken
	}

	storedToken, err := s.authCache.GetRefreshToken(
		ctx,
		claims.UserID,
	)

	if err != nil {
		return "", err
	}

	if storedToken != refreshToken {
		return "", apperrors.ErrInvalidToken
	}

	accessToken, err := auth.GenerateAccessToken(
		claims.UserID,
		s.cfg,
	)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}
