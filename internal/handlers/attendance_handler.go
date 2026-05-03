package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AttendanceHandler exposes HTTP endpoints for check-in and attendance history.
type AttendanceHandler struct {
	attendanceService service.AttendanceService
}

// NewAttendanceHandler wires attendance service into handler.
func NewAttendanceHandler(attendanceService service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{attendanceService: attendanceService}
}

// checkInRequest is JSON body for attendance check-in.
type checkInRequest struct {
	SubscriptionID string `json:"subscription_id"`
	BranchID       string `json:"branch_id"`
	Date           string `json:"date"`
	Status         string `json:"status"`
	IsMakeupFor    string `json:"is_makeup_for"`
}

// CheckIn handles POST /attendance/checkin.
func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	// 1) Parse JSON body.
	var req checkInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	// 2) Validate and convert required IDs.
	subID, err := primitive.ObjectIDFromHex(req.SubscriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid subscription id"})
		return
	}
	branchID, err := primitive.ObjectIDFromHex(req.BranchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch id"})
		return
	}

	// 3) Parse optional date and makeup date.
	attendanceDate := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid date format"})
			return
		}
		attendanceDate = parsed
	}

	var makeupFor *time.Time
	if req.IsMakeupFor != "" {
		parsed, err := time.Parse(time.RFC3339, req.IsMakeupFor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid is_makeup_for format"})
			return
		}
		makeupFor = &parsed
	}

	// 4) Build attendance model and call service.
	attendance := &models.Attendance{
		SubID:       subID,
		BranchID:    branchID,
		Date:        attendanceDate,
		Status:      req.Status,
		IsMakeupFor: makeupFor,
	}
	if err := h.attendanceService.CheckIn(c.Request.Context(), attendance); err != nil {
		// 5) Map domain errors to HTTP status codes.
		switch {
		case errors.Is(err, service.ErrInvalidAttendanceInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case errors.Is(err, service.ErrSubscriptionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "subscription not found"})
		case errors.Is(err, service.ErrAttendanceCheckInNotAllowed), errors.Is(err, service.ErrSubscriptionExpired), errors.Is(err, service.ErrNoRemainingSessions):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	// 6) Success response.
	c.JSON(http.StatusCreated, gin.H{"message": "attendance check-in recorded successfully", "data": attendance})
}

// ListBySubscription handles GET /subscriptions/:id/attendance.
func (h *AttendanceHandler) ListBySubscription(c *gin.Context) {
	// 1) Validate subscription ID in path.
	subscriptionID := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(subscriptionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid subscription id"})
		return
	}

	// 2) Delegate to service and map errors.
	records, err := h.attendanceService.ListBySubscriptionID(c.Request.Context(), subscriptionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAttendanceInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	// 3) Success response.
	c.JSON(http.StatusOK, gin.H{"message": "attendance records fetched successfully", "data": records})
}
