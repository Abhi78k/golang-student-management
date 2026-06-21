package handlers

import (
	"net/http"
	"project-1/internal/dto"
	"project-1/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService *services.CourseService
}

func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// CreateCourse godoc
//
// @Summary Create course
// @Tags Courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateCourseRequest true "Course"
// @Success 201 {object} map[string]string
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body.",
		})
		return
	}

	err := h.courseService.CreateCourse(
		ctx,
		&req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "course created.",
	})
}

// GetCourseByID godoc
//
// @Summary Get course by ID
// @Tags Courses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} dto.CourseResponse
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	course, err := h.courseService.GetCourseByID(ctx, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.CourseResponse{
		ID:             course.ID,
		Name:           course.Name,
		AvailableSeats: course.AvailableSeats,
		CreatedAt:      course.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// ListCourses godoc
//
// @Summary List courses
// @Tags Courses
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {array} dto.CourseResponse
// @Router /courses [get]
func (h *CourseHandler) ListCourses(c *gin.Context) {
	ctx := c.Request.Context()

	pageStr := c.DefaultQuery(
		"page",
		"1",
	)

	limitStr := c.DefaultQuery(
		"limit",
		"2",
	)

	page, err := strconv.Atoi(pageStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid page value.",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid limit value.",
		})
		return
	}

	courses, err := h.courseService.ListCourses(
		ctx,
		page,
		limit,
	)

	var response []dto.CourseResponse

	for _, course := range courses {
		response = append(
			response,
			dto.CourseResponse{
				ID:             course.ID,
				Name:           course.Name,
				AvailableSeats: course.AvailableSeats,
				CreatedAt:      course.CreatedAt,
			},
		)
	}
	c.JSON(http.StatusOK, response)
}

// UpdateCourse godoc
//
// @Summary Update course
// @Tags Courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param request body dto.UpdateCourseRequest true "Course"
// @Success 200 {object} map[string]string
// @Router /courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	var req dto.UpdateCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body.",
		})
		return
	}

	err := h.courseService.UpdateCourse(
		ctx,
		id,
		req.Name,
		req.AvailableSeats,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "course updated.",
	})
}

// DeleteCourse godoc
//
// @Summary Delete course
// @Tags Courses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Success 204
// @Router /courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	err := h.courseService.DeleteCourse(ctx, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
