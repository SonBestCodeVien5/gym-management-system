package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Import đúng tên module dự án của bạn!
	"github.com/SonBestCodeVien5/gym-management-system/internal/handlers"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/SonBestCodeVien5/gym-management-system/pkg/database"
)

func main() {
	// Load env config.
	err := godotenv.Load()
	if err != nil {
		log.Println("Error: Not found .env file, using environment variables instead")
	}

	// Connect to MongoDB using MONGODB_URI.
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("Error: MONGODB_URI is not set in environment variables")
	}

	// ConnectMongoDB returns client and error.
	dbClient, err := database.ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Error: Failed to connect to MongoDB: %v", err)
	}
	// Close DB connection on shutdown.
	defer dbClient.Disconnect(context.Background())

	db := dbClient.Database("gym_management")
	indexCtx, indexCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer indexCancel()
	if err := database.EnsureIndexes(indexCtx, db); err != nil {
		log.Fatalf("Error: Failed to ensure MongoDB indexes: %v", err)
	}
	log.Println("MongoDB indexes ensured successfully")

	// Build repositories.
	memberRepo, err := repository.NewMemberRepository(db)
	if err != nil {
		log.Fatalf("Error: Failed to initialize member repository: %v", err)
	}
	courseRepo := repository.NewCourseRepository(db)
	branchRepo, err := repository.NewBranchRepository(db)
	if err != nil {
		log.Fatalf("Error: Failed to initialize branch repository: %v", err)
	}
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	refundRepo := repository.NewRefundRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	employeeRepo, err := repository.NewEmployeeRepository(db)
	if err != nil {
		log.Fatalf("Error: Failed to initialize employee repository: %v", err)
	}
	refreshTokenRepo, err := repository.NewRefreshTokenRepository(db)
	if err != nil {
		log.Fatalf("Error: Failed to initialize refresh token repository: %v", err)
	}

	// Build services.
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, refundRepo, memberRepo, courseRepo, branchRepo)
	memberService := service.NewMemberService(memberRepo)
	courseService := service.NewCourseService(courseRepo)
	branchService := service.NewBranchService(branchRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo, subscriptionRepo, memberRepo)
	sessionService := service.NewSessionService(sessionRepo, subscriptionRepo, attendanceRepo, attendanceService)
	employeeService := service.NewEmployeeService(employeeRepo, branchRepo, refreshTokenRepo)
	authService, err := service.NewAuthService(employeeRepo, refreshTokenRepo, service.AuthConfig{
		AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		AccessTTL:     durationFromEnv("JWT_ACCESS_TTL", 15*time.Minute),
		RefreshTTL:    durationFromEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
	})
	if err != nil {
		log.Fatalf("Error: Failed to initialize auth service: %v", err)
	}
	if err := authService.BootstrapAdmin(context.Background(), bootstrapAdminFromEnv()); err != nil {
		log.Fatalf("Error: Failed to bootstrap admin employee: %v", err)
	}

	// Build HTTP handlers.
	memberHandler := handlers.NewMemberHandler(memberService, subscriptionService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	courseHandler := handlers.NewCourseHandler(courseService)
	branchHandler := handlers.NewBranchHandler(branchService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	sessionHandler := handlers.NewSessionHandler(sessionService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize Gin engine.
	r := gin.Default()

	// Health check.
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "Backend Go + MongoDB đã sẵn sàng và đang chờ lệnh!",
		})
	})
	// API routes.
	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.Refresh)
		api.POST("/auth/logout", authHandler.Logout)

		protected := api.Group("")
		protected.Use(handlers.AuthRequired(authService))

		employeeRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin))
		employeeRoutes.POST("/employees", employeeHandler.Create)
		employeeRoutes.GET("/employees", employeeHandler.List)
		employeeRoutes.GET("/employees/:id", employeeHandler.GetByID)
		employeeRoutes.PATCH("/employees/:id/password", employeeHandler.UpdatePassword)
		employeeRoutes.PATCH("/employees/:id", employeeHandler.Update)

		memberRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleReceptionist))
		memberRoutes.POST("/members", memberHandler.Register)
		memberRoutes.GET("/members/:id", memberHandler.GetByID)
		memberRoutes.GET("/members/:id/subscriptions", memberHandler.ListSubscriptions)
		memberRoutes.PATCH("/members/:id/activate", memberHandler.Activate)

		courseRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager))
		courseRoutes.POST("/courses", courseHandler.Create)
		courseRoutes.GET("/courses", courseHandler.List)
		courseRoutes.GET("/courses/:id", courseHandler.GetByID)
		courseRoutes.PATCH("/courses/:id", courseHandler.Update)
		courseRoutes.DELETE("/courses/:id", courseHandler.Delete)

		branchRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager))
		branchRoutes.POST("/branches", branchHandler.Create)
		branchRoutes.GET("/branches", branchHandler.List)
		branchRoutes.GET("/branches/nearby", branchHandler.Nearby)
		branchRoutes.GET("/branches/:id", branchHandler.GetByID)
		branchRoutes.PATCH("/branches/:id", branchHandler.Update)
		branchRoutes.DELETE("/branches/:id", branchHandler.Delete)

		subscriptionRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleReceptionist))
		subscriptionRoutes.POST("/subscriptions", subscriptionHandler.Create)
		subscriptionRoutes.POST("/subscriptions/:id/refund", subscriptionHandler.Refund)
		subscriptionRoutes.GET("/subscriptions/:id", subscriptionHandler.GetByID)
		subscriptionRoutes.PATCH("/subscriptions/:id/suspend", subscriptionHandler.Suspend)
		subscriptionRoutes.PATCH("/subscriptions/:id/unsuspend", subscriptionHandler.Resume)
		subscriptionRoutes.PATCH("/subscriptions/:id/expire", subscriptionHandler.Expire)

		attendanceRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleReceptionist))
		attendanceRoutes.POST("/attendance/checkin", attendanceHandler.CheckIn)
		attendanceRoutes.POST("/attendance/report", attendanceHandler.ReportMissed)
		attendanceRoutes.POST("/attendance/makeup", attendanceHandler.Makeup)
		attendanceRoutes.GET("/subscriptions/:id/attendance", attendanceHandler.ListBySubscription)

		sessionCreateRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleTrainer))
		sessionCreateRoutes.POST("/sessions", sessionHandler.Create)
		sessionRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleTrainer))
		sessionRoutes.GET("/sessions", sessionHandler.List)
		sessionRoutes.GET("/sessions/:id", sessionHandler.GetByID)
		sessionRoutes.POST("/sessions/:id/enroll", sessionHandler.Enroll)
		sessionRoutes.POST("/sessions/:id/checkin", sessionHandler.CheckIn)
	}

	// Start HTTP server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on port %s...", port)

	// Khởi động server (sẽ chạy vòng lặp vô hạn ở đây)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error: Failed to start server: %v", err)
	}
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	duration, err := time.ParseDuration(raw)
	if err != nil {
		log.Fatalf("Error: %s must be a valid duration: %v", key, err)
	}
	return duration
}

func bootstrapAdminFromEnv() service.BootstrapAdminConfig {
	return service.BootstrapAdminConfig{
		EmployeeID: os.Getenv("BOOTSTRAP_ADMIN_EMPLOYEE_ID"),
		FullName:   os.Getenv("BOOTSTRAP_ADMIN_FULL_NAME"),
		Email:      os.Getenv("BOOTSTRAP_ADMIN_EMAIL"),
		Password:   os.Getenv("BOOTSTRAP_ADMIN_PASSWORD"),
		Phone:      os.Getenv("BOOTSTRAP_ADMIN_PHONE"),
		Level:      os.Getenv("BOOTSTRAP_ADMIN_LEVEL"),
	}
}
