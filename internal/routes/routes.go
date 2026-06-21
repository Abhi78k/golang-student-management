package routes

import (
	"project-1/internal/cache"
	"project-1/internal/config"
	"project-1/internal/handlers"
	"project-1/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	studentHandler *handlers.StudentHandler,
	courseHandler *handlers.CourseHandler,
	enrollmentHandler *handlers.EnrollmentHandler,
	rateLimiter *cache.RateLimiter,
	cfg *config.Config,
) *gin.Engine {

	router := gin.Default()

	protected := router.Group("/auth")
	protected.Use(middleware.AuthMiddleware(cfg))

	auth := router.Group("/auth")
	auth.POST(
		"/register",
		authHandler.Register,
	)

	auth.POST(
		"/login",
		authHandler.Login,
	)

	auth.POST(
		"/refresh",
		authHandler.Refresh,
	)

	protected.GET(
		"/me",
		authHandler.GetMe,
	)

	protected.POST(
		"/logout",
		authHandler.Logout,
	)

	students := router.Group("/students")
	students.Use(
		middleware.AuthMiddleware(cfg),
		middleware.RateLimitMiddleware(rateLimiter),
	)

	students.POST(
		"",
		studentHandler.CreateStudent,
	)

	students.GET(
		"",
		studentHandler.ListStudents,
	)

	students.GET(
		"/:id",
		studentHandler.FindStudentByID,
	)

	students.PUT(
		"/:id",
		studentHandler.UpdateStudent,
	)

	students.DELETE(
		"/:id",
		studentHandler.DeleteStudent,
	)

	students.PATCH(
		"/:id",
		studentHandler.UpdateStudentDetail,
	)

	courses := router.Group("/courses")
	courses.Use(
		middleware.AuthMiddleware(cfg),
		middleware.RateLimitMiddleware(rateLimiter),
	)

	courses.POST(
		"",
		courseHandler.CreateCourse,
	)

	courses.GET(
		"/:id",
		courseHandler.GetCourseByID,
	)

	courses.GET(
		"",
		courseHandler.ListCourses,
	)

	courses.PUT(
		"/:id",
		courseHandler.UpdateCourse,
	)

	courses.DELETE(
		"/:id",
		courseHandler.DeleteCourse,
	)

	courses.POST(
		"/:id/enroll",
		enrollmentHandler.Enroll,
	)

	return router
}
