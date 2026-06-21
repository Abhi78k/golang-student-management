package handlers

import (
	"errors"
	"net/http"
	"project-1/internal/apperrors"
	"project-1/internal/dto"
	"project-1/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentService *services.StudentService
}

func NewStudentHandler(studentService *services.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// CreateStudent godoc
//
// @Summary Create student
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateStudentRequest true "Student"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /students [post]
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body.",
		})
		return
	}

	err := h.studentService.CreateStudent(
		ctx,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Age,
	)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"message": "student created.",
		},
	)
}

// FindStudentByID godoc
//
// @Summary Get student by ID
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} dto.StudentResponse
// @Failure 404 {object} map[string]string
// @Router /students/{id} [get]
func (h *StudentHandler) FindStudentByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	student, err := h.studentService.GetStudentByID(ctx, id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "student not found.",
		})
		return
	}

	response := dto.StudentResponse{
		ID:        student.ID,
		FirstName: student.FirstName,
		LastName:  student.LastName,
		Email:     student.Email,
		Age:       student.Age,
	}

	c.JSON(http.StatusOK, response)
}

// ListStudents godoc
//
// @Summary List students
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param search query string false "Search"
// @Param age query int false "Age"
// @Param sort query string false "Sort"
// @Success 200 {array} dto.StudentResponse
// @Router /students [get]
func (h *StudentHandler) ListStudents(c *gin.Context) {
	ctx := c.Request.Context()

	pageStr := c.DefaultQuery(
		"page",
		"1",
	)

	limitStr := c.DefaultQuery(
		"limit",
		"2",
	)

	ageStr := c.Query("age")

	search := c.Query("search")

	sort := c.Query("sort")

	var age int

	if ageStr != "" {
		parsedAge, err := strconv.Atoi(ageStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid age value.",
			})
			return
		}

		age = parsedAge
	}

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

	filter := dto.StudentFilter{
		Page:   page,
		Limit:  limit,
		Age:    age,
		Search: search,
		Sort:   sort,
	}

	students, err := h.studentService.ListStudents(ctx, filter)

	var response []dto.StudentResponse

	for _, student := range students {
		response = append(
			response,
			dto.StudentResponse{
				ID:        student.ID,
				FirstName: student.FirstName,
				LastName:  student.LastName,
				Email:     student.Email,
				Age:       student.Age,
			},
		)
	}
	c.JSON(http.StatusOK, response)
}

// UpdateStudent godoc
//
// @Summary Update student
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body dto.UpdateStudentRequest true "Student"
// @Success 200 {object} map[string]string
// @Router /students/{id} [put]
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	var req dto.UpdateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body.",
		})
		return
	}

	err := h.studentService.UpdateStudent(
		ctx,
		id,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Age,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "status updated.",
	})
}

// DeleteStudent godoc
//
// @Summary Delete student
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 204
// @Router /students/{id} [delete]
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := h.studentService.DeleteStudent(ctx, id)

	if errors.Is(
		err,
		apperrors.ErrStudentNotFound,
	) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "student not found.",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateStudentDetail godoc
//
// @Summary Partially update student
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body dto.PatchStudentRequest true "Patch Student"
// @Success 200 {object} map[string]string
// @Router /students/{id} [patch]
func (h *StudentHandler) UpdateStudentDetail(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.PatchStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body.",
		})
		return
	}

	id := c.Param("id")

	err := h.studentService.UpdateStudentDetail(ctx, id, req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "student updated.",
	})
}
