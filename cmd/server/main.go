package main

import (
	"context"
	"log"
	"os"

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

	// Build repositories.
	memberRepo, err := repository.NewMemberRepository(db)
	if err != nil {
		log.Fatalf("Error: Failed to initialize member repository: %v", err)
	}
	courseRepo := repository.NewCourseRepository(db)
	branchRepo := repository.NewBranchRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Build services.
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, memberRepo, courseRepo, branchRepo)
	memberService := service.NewMemberService(memberRepo)
	courseService := service.NewCourseService(courseRepo)
	branchService := service.NewBranchService(branchRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo, subscriptionRepo, memberRepo)
	sessionService := service.NewSessionService(sessionRepo)

	// Build HTTP handlers.
	memberHandler := handlers.NewMemberHandler(memberService, subscriptionService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	courseHandler := handlers.NewCourseHandler(courseService)
	branchHandler := handlers.NewBranchHandler(branchService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	sessionHandler := handlers.NewSessionHandler(sessionService)

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
		api.POST("/members", memberHandler.Register)
		api.GET("/members/:id", memberHandler.GetByID)
		api.PATCH("/members/:id/activate", memberHandler.Activate)

		api.POST("/courses", courseHandler.Create)
		api.GET("/courses", courseHandler.List)
		api.GET("/courses/:id", courseHandler.GetByID)
		api.PATCH("/courses/:id", courseHandler.Update)
		api.DELETE("/courses/:id", courseHandler.Delete)

		api.POST("/branches", branchHandler.Create)
		api.GET("/branches", branchHandler.List)
		api.GET("/branches/:id", branchHandler.GetByID)
		api.PATCH("/branches/:id", branchHandler.Update)
		api.DELETE("/branches/:id", branchHandler.Delete)

		api.POST("/subscriptions", subscriptionHandler.Create)
		api.GET("/subscriptions/:id", subscriptionHandler.GetByID)
		api.PATCH("/subscriptions/:id/suspend", subscriptionHandler.Suspend)
		api.PATCH("/subscriptions/:id/unsuspend", subscriptionHandler.Resume)
		api.PATCH("/subscriptions/:id/expire", subscriptionHandler.Expire)

		api.POST("/attendance/checkin", attendanceHandler.CheckIn)
		api.GET("/subscriptions/:id/attendance", attendanceHandler.ListBySubscription)

		api.POST("/sessions", sessionHandler.Create)
		api.GET("/sessions", sessionHandler.List)
		api.GET("/sessions/:id", sessionHandler.GetByID)
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
