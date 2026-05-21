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
	SessionID      string `json:"session_id"`
	Date           string `json:"date"`
	Status         string `json:"status"`
	IsMakeupFor    string `json:"is_makeup_for"`
}

type attendanceReportRequest struct {
	SubscriptionID string `json:"subscription_id"`
	BranchID       string `json:"branch_id"`
	Date           string `json:"date"`
}

type attendanceMakeupRequest struct {
	SubscriptionID string `json:"subscription_id"`
	BranchID       string `json:"branch_id"`
	Date           string `json:"date"`
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

	var sessionID *primitive.ObjectID
	if req.SessionID != "" {
		parsed, err := primitive.ObjectIDFromHex(req.SessionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid session id"})
			return
		}
		sessionID = &parsed
	}

	// 4) Build attendance model and call service.
	attendance := &models.Attendance{
		SubID:       subID,
		BranchID:    branchID,
		SessionID:   sessionID,
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
		case errors.Is(err, service.ErrWeeklySessionLimitReached):
			c.JSON(http.StatusConflict, gin.H{"message": "weekly session limit reached"})
		case errors.Is(err, service.ErrReportedMissedLimitReached), errors.Is(err, service.ErrMakeupReferenceInvalid), errors.Is(err, service.ErrMakeupReferenceNotFound), errors.Is(err, service.ErrMakeupAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	// 6) Success response.
	c.JSON(http.StatusCreated, gin.H{"message": "attendance check-in recorded successfully", "data": attendance})
}

// ReportMissed handles POST /attendance/report.
func (h *AttendanceHandler) ReportMissed(c *gin.Context) {
	var req attendanceReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

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

	attendanceDate := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid date format"})
			return
		}
		attendanceDate = parsed
	}

	attendance := &models.Attendance{
		SubID:    subID,
		BranchID: branchID,
		Date:     attendanceDate,
		Status:   "reported_missed",
	}
	if err := h.attendanceService.CheckIn(c.Request.Context(), attendance); err != nil {
		h.handleAttendanceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "attendance report recorded successfully", "data": attendance})
}

// Makeup handles POST /attendance/makeup.
func (h *AttendanceHandler) Makeup(c *gin.Context) {
	var req attendanceMakeupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

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

	attendanceDate := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid date format"})
			return
		}
		attendanceDate = parsed
	}

	if req.IsMakeupFor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "is_makeup_for is required"})
		return
	}
	makeupFor, err := time.Parse(time.RFC3339, req.IsMakeupFor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid is_makeup_for format"})
		return
	}

	attendance := &models.Attendance{
		SubID:       subID,
		BranchID:    branchID,
		Date:        attendanceDate,
		Status:      "makeup",
		IsMakeupFor: &makeupFor,
	}
	if err := h.attendanceService.CheckIn(c.Request.Context(), attendance); err != nil {
		h.handleAttendanceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "attendance makeup recorded successfully", "data": attendance})
}

func (h *AttendanceHandler) handleAttendanceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidAttendanceInput):
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	case errors.Is(err, service.ErrSubscriptionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"message": "subscription not found"})
	case errors.Is(err, service.ErrAttendanceCheckInNotAllowed), errors.Is(err, service.ErrSubscriptionExpired), errors.Is(err, service.ErrNoRemainingSessions):
		c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
	case errors.Is(err, service.ErrWeeklySessionLimitReached):
		c.JSON(http.StatusConflict, gin.H{"message": "weekly session limit reached"})
	case errors.Is(err, service.ErrReportedMissedLimitReached), errors.Is(err, service.ErrMakeupReferenceInvalid), errors.Is(err, service.ErrMakeupReferenceNotFound), errors.Is(err, service.ErrMakeupAlreadyUsed):
		c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}
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
