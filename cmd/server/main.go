package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 1. In ra lời chào Hello World
	fmt.Println("--- GYM MANAGEMENT SYSTEM ---")
	fmt.Println("Status: Initializing Backend...")

	// 2. Thiết lập chuỗi kết nối (URI)
	// Lưu ý: Dùng đúng username/password bạn đã đặt trong docker-compose
	uri := "mongodb://admin:password123@localhost:27017"

	// 3. Tạo context với thời gian chờ (timeout) 10 giây
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Đổi sang 10*time.Second
	defer cancel()

	// 4. Kết nối tới MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("❌ Lỗi kết nối ban đầu: %v", err)
	}

	// 5. Ping thử xem Database có phản hồi không
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Không thể phản hồi MongoDB: %v", err)
	}

	fmt.Println("✅ Kết nối MongoDB THÀNH CÔNG!")
	fmt.Println("🚀 Backend đã sẵn sàng phục vụ Sơn tại Linux WSL 2!")

	// 6. Ngắt kết nối khi xong (Cleanup)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}