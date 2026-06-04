package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/SonBestCodeVien5/gym-management-system/internal/app"
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

	db := dbClient.Database(databaseNameFromEnv())
	indexCtx, indexCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer indexCancel()
	if err := database.EnsureIndexes(indexCtx, db); err != nil {
		log.Fatalf("Error: Failed to ensure MongoDB indexes: %v", err)
	}
	log.Println("MongoDB indexes ensured successfully")

	r, err := app.NewRouter(context.Background(), db, app.Config{
		Auth: service.AuthConfig{
			AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
			RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
			AccessTTL:     durationFromEnv("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL:    durationFromEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
		},
		BootstrapAdmin: bootstrapAdminFromEnv(),
		CORSOrigins:    csvFromEnv("CORS_ALLOWED_ORIGINS"),
	})
	if err != nil {
		log.Fatalf("Error: Failed to initialize app router: %v", err)
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

func csvFromEnv(key string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}

func databaseNameFromEnv() string {
	dbName := strings.TrimSpace(os.Getenv("DB_NAME"))
	if dbName == "" {
		return "gym_management"
	}
	return dbName
}
