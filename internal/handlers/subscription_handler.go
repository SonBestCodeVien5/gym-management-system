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

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

type createSubscriptionRequest struct {
	MemberID       string `json:"member_id"`
	CourseID       string `json:"course_id"`
	HomeBranchID   string `json:"home_branch_id"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	SessionPerWeek int    `json:"session_per_week"`
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req createSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	memberID, err := primitive.ObjectIDFromHex(req.MemberID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid member id"})
		return
	}

	courseID, err := primitive.ObjectIDFromHex(req.CourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid course id"})
		return
	}

	branchID, err := primitive.ObjectIDFromHex(req.HomeBranchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch id"})
		return
	}

	startDate, err := parseTimeValue(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid start_date format"})
		return
	}

	endDate, err := parseTimeValue(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid end_date format"})
		return
	}

	subscription := &models.Subscription{
		MemberID:       memberID,
		CourseID:       courseID,
		HomeBranchID:   branchID,
		StartDate:      startDate,
		EndDate:        endDate,
		SessionPerWeek: req.SessionPerWeek,
	}

	err = h.subscriptionService.CreateSubscription(c.Request.Context(), subscription)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSubscriptionInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case errors.Is(err, service.ErrSubscriptionReferenceNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "subscription created successfully",
		"data":    subscription,
	})
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid subscription id"})
		return
	}

	subscription, err := h.subscriptionService.GetSubscriptionByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSubscriptionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "subscription not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subscription fetched successfully",
		"data":    subscription,
	})
}

func parseTimeValue(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}
