package app

import (
	"context"
	"net/http"

	"github.com/SonBestCodeVien5/gym-management-system/internal/handlers"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Auth           service.AuthConfig
	BootstrapAdmin service.BootstrapAdminConfig
}

func NewRouter(ctx context.Context, db *mongo.Database, cfg Config) (*gin.Engine, error) {
	memberRepo, err := repository.NewMemberRepository(db)
	if err != nil {
		return nil, err
	}
	courseRepo := repository.NewCourseRepository(db)
	branchRepo, err := repository.NewBranchRepository(db)
	if err != nil {
		return nil, err
	}
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	refundRepo := repository.NewRefundRepository(db)
	attendanceRepo := repository.NewAttendanceRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	employeeRepo, err := repository.NewEmployeeRepository(db)
	if err != nil {
		return nil, err
	}
	refreshTokenRepo, err := repository.NewRefreshTokenRepository(db)
	if err != nil {
		return nil, err
	}

	subscriptionService := service.NewSubscriptionService(subscriptionRepo, refundRepo, memberRepo, courseRepo, branchRepo)
	memberService := service.NewMemberService(memberRepo)
	courseService := service.NewCourseService(courseRepo)
	branchService := service.NewBranchService(branchRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo, subscriptionRepo, memberRepo)
	sessionService := service.NewSessionService(sessionRepo, subscriptionRepo, attendanceRepo, attendanceService)
	employeeService := service.NewEmployeeService(employeeRepo, branchRepo, refreshTokenRepo)
	authService, err := service.NewAuthService(employeeRepo, refreshTokenRepo, cfg.Auth)
	if err != nil {
		return nil, err
	}
	if err := authService.BootstrapAdmin(ctx, cfg.BootstrapAdmin); err != nil {
		return nil, err
	}

	memberHandler := handlers.NewMemberHandler(memberService, subscriptionService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	courseHandler := handlers.NewCourseHandler(courseService)
	branchHandler := handlers.NewBranchHandler(branchService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	sessionHandler := handlers.NewSessionHandler(sessionService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()
	RegisterRoutes(r, Handlers{
		Member:       memberHandler,
		Subscription: subscriptionHandler,
		Course:       courseHandler,
		Branch:       branchHandler,
		Attendance:   attendanceHandler,
		Session:      sessionHandler,
		Employee:     employeeHandler,
		Auth:         authHandler,
		AuthService:  authService,
	})
	return r, nil
}

type Handlers struct {
	Member       *handlers.MemberHandler
	Subscription *handlers.SubscriptionHandler
	Course       *handlers.CourseHandler
	Branch       *handlers.BranchHandler
	Attendance   *handlers.AttendanceHandler
	Session      *handlers.SessionHandler
	Employee     *handlers.EmployeeHandler
	Auth         *handlers.AuthHandler
	AuthService  service.AuthService
}

func RegisterRoutes(r *gin.Engine, h Handlers) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "Backend Go + MongoDB đã sẵn sàng và đang chờ lệnh!",
		})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", h.Auth.Login)
		api.POST("/auth/refresh", h.Auth.Refresh)
		api.POST("/auth/logout", h.Auth.Logout)

		protected := api.Group("")
		protected.Use(handlers.AuthRequired(h.AuthService))

		employeeRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin))
		employeeRoutes.POST("/employees", h.Employee.Create)
		employeeRoutes.GET("/employees", h.Employee.List)
		employeeRoutes.GET("/employees/:id", h.Employee.GetByID)
		employeeRoutes.PATCH("/employees/:id/password", h.Employee.UpdatePassword)
		employeeRoutes.PATCH("/employees/:id", h.Employee.Update)

		memberRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleReceptionist))
		memberRoutes.POST("/members", h.Member.Register)
		memberRoutes.GET("/members/:id", h.Member.GetByID)
		memberRoutes.GET("/members/:id/subscriptions", h.Member.ListSubscriptions)
		memberRoutes.PATCH("/members/:id/activate", h.Member.Activate)

		courseRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager))
		courseRoutes.POST("/courses", h.Course.Create)
		courseRoutes.GET("/courses", h.Course.List)
		courseRoutes.GET("/courses/:id", h.Course.GetByID)
		courseRoutes.PATCH("/courses/:id", h.Course.Update)
		courseRoutes.DELETE("/courses/:id", h.Course.Delete)

		branchRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager))
		branchRoutes.POST("/branches", h.Branch.Create)
		branchRoutes.GET("/branches", h.Branch.List)
		branchRoutes.GET("/branches/nearby", h.Branch.Nearby)
		branchRoutes.GET("/branches/:id", h.Branch.GetByID)
		branchRoutes.PATCH("/branches/:id", h.Branch.Update)
		branchRoutes.DELETE("/branches/:id", h.Branch.Delete)

		subscriptionRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleReceptionist))
		subscriptionRoutes.POST("/subscriptions", h.Subscription.Create)
		subscriptionRoutes.POST("/subscriptions/:id/refund", h.Subscription.Refund)
		subscriptionRoutes.GET("/subscriptions/:id", h.Subscription.GetByID)
		subscriptionRoutes.PATCH("/subscriptions/:id/suspend", h.Subscription.Suspend)
		subscriptionRoutes.PATCH("/subscriptions/:id/unsuspend", h.Subscription.Resume)
		subscriptionRoutes.PATCH("/subscriptions/:id/expire", h.Subscription.Expire)

		attendanceRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleReceptionist))
		attendanceRoutes.POST("/attendance/checkin", h.Attendance.CheckIn)
		attendanceRoutes.POST("/attendance/report", h.Attendance.ReportMissed)
		attendanceRoutes.POST("/attendance/makeup", h.Attendance.Makeup)
		attendanceRoutes.GET("/subscriptions/:id/attendance", h.Attendance.ListBySubscription)

		sessionCreateRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleTrainer))
		sessionCreateRoutes.POST("/sessions", h.Session.Create)
		sessionRoutes := protected.Group("", handlers.RequireRoles(service.RoleAdmin, service.RoleManager, service.RoleTrainer))
		sessionRoutes.GET("/sessions", h.Session.List)
		sessionRoutes.GET("/sessions/:id", h.Session.GetByID)
		sessionRoutes.POST("/sessions/:id/enroll", h.Session.Enroll)
		sessionRoutes.POST("/sessions/:id/checkin", h.Session.CheckIn)
	}
}
