package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionHandler struct {
	sessionService service.SessionService
}

func NewSessionHandler(sessionService service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

type sessionRequest struct {
	BranchID    string   `json:"branch_id"`
	TrainerID   string   `json:"trainer_id"`
	CourseLevel string   `json:"course_level"`
	ScheduledAt string   `json:"scheduled_at"`
	DurationMin int      `json:"duration_min"`
	Capacity    int      `json:"capacity"`
	Tags        []string `json:"tags"`
}

type sessionEnrollRequest struct {
	SubscriptionID string `json:"subscription_id"`
}

func (h *SessionHandler) Create(c *gin.Context) {
	var req sessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	branchID, err := primitive.ObjectIDFromHex(req.BranchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch id"})
		return
	}

	trainerID, err := primitive.ObjectIDFromHex(req.TrainerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid trainer id"})
		return
	}

	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid scheduled_at format"})
		return
	}

	session := &models.Session{
		BranchID:    branchID,
		TrainerID:   trainerID,
		CourseLevel: req.CourseLevel,
		ScheduledAt: scheduledAt,
		DurationMin: req.DurationMin,
		Capacity:    req.Capacity,
		Tags:        req.Tags,
	}

	if err := h.sessionService.CreateSession(c.Request.Context(), session); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSessionInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "session created successfully", "data": session})
}

func (h *SessionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid session id"})
		return
	}

	session, err := h.sessionService.GetSessionByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSessionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "session not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session fetched successfully", "data": session})
}

func (h *SessionHandler) List(c *gin.Context) {
	var filter repository.SessionListFilter
	filter.BranchID = c.Query("branchId")
	filter.Level = c.Query("level")
	if dateParam := c.Query("date"); dateParam != "" {
		parsed, err := time.Parse(time.RFC3339, dateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid date format"})
			return
		}
		filter.Date = &parsed
	}

	sessions, err := h.sessionService.ListSessions(c.Request.Context(), filter)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSessionInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sessions fetched successfully", "data": sessions})
}

// Enroll handles POST /sessions/:id/enroll.
func (h *SessionHandler) Enroll(c *gin.Context) {
	sessionID := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(sessionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid session id"})
		return
	}

	var req sessionEnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	session, err := h.sessionService.EnrollSubscription(c.Request.Context(), sessionID, req.SubscriptionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSessionInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case errors.Is(err, service.ErrSessionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "session not found"})
		case errors.Is(err, service.ErrSessionTagNotAllowed), errors.Is(err, service.ErrSessionAlreadyEnrolled), errors.Is(err, service.ErrSessionAlreadyFull):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session enrolled successfully", "data": session})
}

// CheckIn handles POST /sessions/:id/checkin.
func (h *SessionHandler) CheckIn(c *gin.Context) {
	sessionID := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(sessionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid session id"})
		return
	}

	var req sessionEnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	attendance, err := h.sessionService.CheckInSubscription(c.Request.Context(), sessionID, req.SubscriptionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSessionInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case errors.Is(err, service.ErrSessionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "session not found"})
		case errors.Is(err, service.ErrSessionNotEnrolled), errors.Is(err, service.ErrSessionCheckInClosed):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		case errors.Is(err, service.ErrAttendanceCheckInNotAllowed), errors.Is(err, service.ErrSubscriptionExpired), errors.Is(err, service.ErrNoRemainingSessions), errors.Is(err, service.ErrWeeklySessionLimitReached):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "session check-in recorded successfully", "data": attendance})
}
