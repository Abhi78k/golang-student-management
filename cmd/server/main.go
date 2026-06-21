package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	_ "project-1/docs"
	"syscall"
	"time"

	"log"
	"project-1/internal/cache"
	"project-1/internal/config"
	"project-1/internal/database"
	"project-1/internal/handlers"
	"project-1/internal/repositories"
	"project-1/internal/routes"
	"project-1/internal/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Student Management API
// @version 1.0
// @description Student Management System built with Go, Gin, PostgreSQL and Redis.
// @BasePath /

func main() {

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	redisClient := cache.ConnectRedis(cfg)

	defer redisClient.Close()

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

	router := routes.SetupRouter(authHandler, studentHandler, courseHandler, enrollmentHandler, rateLimiter, cfg)

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
		),
	)

	// if err := router.Run(":8080"); err != nil {
	// 	log.Fatal(err)
	// }

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatalf(
				"server error: %v",
				err,
			)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf(
			"server forced to shutdown: %v",
			err,
		)
	}

	db.Close()
	log.Println("database closed.")

	redisClient.Close()
	log.Println("redis closed.")

	log.Println("server exited cleanly")
}
