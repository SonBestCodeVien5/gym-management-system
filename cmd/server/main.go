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

	// Build services.
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, memberRepo, courseRepo, branchRepo)
	memberService := service.NewMemberService(memberRepo)

	// Build HTTP handlers.
	memberHandler := handlers.NewMemberHandler(memberService, subscriptionService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

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
		api.POST("/registration", memberHandler.Register)
		api.GET("/members/:id", memberHandler.GetByID)
		api.PATCH("/members/:id/activate", memberHandler.Activate)
		api.POST("/subscriptions", subscriptionHandler.Create)
		api.GET("/subscriptions/:id", subscriptionHandler.GetByID)
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
