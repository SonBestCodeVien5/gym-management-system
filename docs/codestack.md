# Code Stack & Project Basics

## 1. Môi trường phát triển (Environment)
* **OS:** Linux (Ubuntu trên WSL 2)
* **Containerization:** Docker & Docker Compose (Quản lý Database và các services phụ trợ)
* **Version Control:** Git & GitHub

## 2. Backend Stack
* **Ngôn ngữ:** Go (Golang) v1.26+
* **Web Framework:** Gin (`gin-gonic/gin`) - Xử lý routing, middleware và HTTP requests nhanh chóng.
* **Architecture:** Clean Architecture (Domain Driven Design cơ bản).
  * `models/`: Định nghĩa cấu trúc dữ liệu (Structs & Tags).
  * `repository/`: Tương tác trực tiếp với Database.
  * `service/`: Xử lý logic nghiệp vụ (Tính toán, kiểm tra luật).
  * `handlers/`: Tiếp nhận HTTP Request và trả về HTTP Response (JSON).

## 3. Database
* **Hệ quản trị:** MongoDB (NoSQL)
* **Driver:** `go.mongodb.org/mongo-driver`
* **Lý do chọn:** Lược đồ động (Dynamic Schema) phù hợp với các cấu trúc lồng nhau (Embedded Documents) như lịch sử bảo lưu (`SuspensionHistory`).

## 4. Các module chuẩn bị tích hợp (Upcoming)
* **Authentication:** JWT (JSON Web Token) - Tái sử dụng logic từ project Auth cũ để phân quyền Admin/Trainer/Member.