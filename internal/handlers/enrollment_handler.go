package handlers

import (
	"errors"
	"net/http"
	"project-1/internal/apperrors"
	"project-1/internal/services"

	"github.com/gin-gonic/gin"
)

type EnrollmentHandler struct {
	enrollmentService *services.EnrollmentService
}

func NewEnrollmentHandler(enrollmentService *services.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollmentService: enrollmentService,
	}
}

// Enroll godoc
//
// @Summary Enroll in course
// @Description Enroll currently authenticated student into a course
// @Tags Enrollments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /courses/{id}/enroll [post]
func (h *EnrollmentHandler) Enroll(c *gin.Context) {
	ctx := c.Request.Context()

	courseID := c.Param("id")

	userID, exists := c.Get("UserID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found.",
		})
		return
	}

	studentID := userID.(string)

	err := h.enrollmentService.Enroll(ctx, studentID, courseID)

	if errors.Is(
		err,
		apperrors.ErrAlreadyEnrolled,
	) {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	if errors.Is(
		err,
		apperrors.ErrNoSeatsAvailable,
	) {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	if errors.Is(
		err,
		apperrors.ErrCourseNotFound,
	) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "enrolled successfully.",
	})
}
