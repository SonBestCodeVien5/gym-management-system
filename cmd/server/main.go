package main

import (
    "context"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    
    // Import đúng tên module dự án của bạn!
    "github.com/SonBestCodeVien5/gym-management-system/pkg/database" 
)

func main() {
    // 1. Tải file cấu hình .env
    err := godotenv.Load()
    if err != nil {
        log.Println("Error: Not found .env file, using environment variables instead")
    }

    // 2. Kết nối Database
    // Trong file .env của chúng ta, nó tên là MONGODB_URI (chứ không phải MONGO_URI)
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        log.Fatal("Error: MONGODB_URI is not set in environment variables")
    }

    // Hàm ConnectMongoDB trả về 2 giá trị: client và err (phải hứng cả err)
    dbClient, err := database.ConnectMongoDB(mongoURI)
    if err != nil {
        log.Fatalf("Error: Failed to connect to MongoDB: %v", err)
    }
    // Đóng kết nối an toàn khi tắt server
    defer dbClient.Disconnect(context.Background())

    // 3. Khởi tạo Gin Engine (Web framework)
    r := gin.Default()

    // 4. Định nghĩa Route đơn giản để test HTTP
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
            "status":  "Backend Go + MongoDB đã sẵn sàng và đang chờ lệnh!",
        })
    })

    // 5. Chạy Server
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