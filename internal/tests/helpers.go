package tests

import (
	"database/sql"
	"log"
	"project-1/internal/cache"
	"project-1/internal/config"
	"project-1/internal/database"
	"project-1/internal/handlers"
	"project-1/internal/repositories"
	"project-1/internal/routes"
	"project-1/internal/services"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestConfig() *config.Config {
	return &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "hydrogen",
		DBName:     "student_management_test",

		RedisAddr: "localhost:6379",

		JWTAccessSecret:  "test-access-secret",
		JWTRefreshSecret: "test-refresh-secret",
	}
}

func setupTestRouter(
	t *testing.T,
) (router *gin.Engine, db *sql.DB) {
	cfg := setupTestConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := cache.ConnectRedis(cfg)

	studentCache := cache.NewStudentCache(redisClient)
	authCache := cache.NewAuthCache(redisClient)

	rateLimiter := cache.NewRateLimiter(redisClient)

	userRepo := repositories.NewUserRepository(db)
	studentRepo := repositories.NewStudentRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	enrollmentRepo := repositories.NewEnrollmentRepository(db)

	authService := services.NewAuthService(
		userRepo,
		cfg,
		authCache,
	)

	studentService := services.NewStudentService(
		studentRepo,
		cfg,
		studentCache,
	)

	courseService := services.NewCourseService(
		courseRepo,
	)

	enrollmentService := services.NewEnrollmentService(
		db,
		courseRepo,
		enrollmentRepo,
	)

	authHandler := handlers.NewAuthHandler(authService)
	studentHandler := handlers.NewStudentHandler(studentService)
	courseHandler := handlers.NewCourseHandler(courseService)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService)

	router = routes.SetupRouter(authHandler, studentHandler, courseHandler, enrollmentHandler, rateLimiter, cfg)

	return router, db
}

func cleanDatabase(
	db *sql.DB,
) {

	query := `
		TRUNCATE TABLE
			users,
			students,
			courses,
			enrollments
		RESTART IDENTITY CASCADE
	`

	_, err := db.Exec(
		query,
	)

	if err != nil {
		log.Fatal(err)
	}
}
