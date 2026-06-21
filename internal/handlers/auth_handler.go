package handlers

import (
	"net/http"
	"project-1/internal/dto"
	"project-1/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(
	authService *services.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
//
// @Summary Register a new user
// @Description Create a new account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request.",
		})
		return
	}

	err := h.authService.Register(
		ctx,
		req.Username,
		req.Email,
		req.Password,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"message": "user created.",
		},
	)
}

// Login godoc
//
// @Summary Login user
// @Description Login using email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body.",
		})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(
		ctx,
		req.Email,
		req.Password,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	loginResponse := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, loginResponse)
}

// GetMe godoc
//
// @Summary Get current user
// @Description Returns currently authenticated user
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("UserID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found.",
		})
		return
	}

	uid, ok := userID.(string)

	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid user id",
			},
		)
		return
	}

	user, err := h.authService.GetMe(ctx, uid)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found.",
		})
		return
	}

	c.JSON(
		http.StatusOK,
		dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	)
}

// Refresh godoc
//
// @Summary Refresh access token
// @Description Generate a new access token using a refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh Request"
// @Success 200 {object} dto.RefreshResponse
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body.",
		})
		return
	}

	accessToken, err := h.authService.Refresh(
		c.Request.Context(),
		req.RefreshToken,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,
		dto.RefreshResponse{
			AccessToken: accessToken,
		},
	)
}

// Logout godoc
//
// @Summary Logout user
// @Description Removes refresh token from Redis
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {

	ctx := c.Request.Context()

	userID, exists := c.Get("UserID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found.",
		})
		return
	}

	uid := userID.(string)

	err := h.authService.Logout(
		ctx,
		uid,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully.",
	})
}
