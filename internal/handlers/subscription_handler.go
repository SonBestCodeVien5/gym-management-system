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

// NewSubscriptionHandler wires subscription service into HTTP handlers.
func NewSubscriptionHandler(subscriptionService service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// createSubscriptionRequest is the JSON body for creating a subscription.
type createSubscriptionRequest struct {
	MemberID       string `json:"member_id"`
	CourseID       string `json:"course_id"`
	HomeBranchID   string `json:"home_branch_id"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	SessionPerWeek int    `json:"session_per_week"`
	DiscountType   string `json:"discount_type"`
	DiscountValue  int64  `json:"discount_value"`
	PromoCode      string `json:"promo_code"`
}

// suspendSubscriptionRequest is the JSON body for suspending a subscription.
type suspendSubscriptionRequest struct {
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	FrozenSession int    `json:"frozen_session"`
	Reason        string `json:"reason"`
}

// refundSubscriptionRequest is the JSON body for refunding a subscription.
type refundSubscriptionRequest struct {
	Reason string `json:"reason"`
}

// Create validates input and delegates to subscription service.
func (h *SubscriptionHandler) Create(c *gin.Context) {
	// 1) Parse JSON body into request struct.
	var req createSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondInvalidRequestBody(c)
		return
	}

	// 2) Validate and convert IDs from request body.
	memberID, err := primitive.ObjectIDFromHex(req.MemberID)
	if err != nil {
		RespondInvalidID(c, "invalid member id")
		return
	}

	courseID, err := primitive.ObjectIDFromHex(req.CourseID)
	if err != nil {
		RespondInvalidID(c, "invalid course id")
		return
	}

	branchID, err := primitive.ObjectIDFromHex(req.HomeBranchID)
	if err != nil {
		RespondInvalidID(c, "invalid branch id")
		return
	}

	// 3) Parse RFC3339 dates from body.
	startDate, err := parseTimeValue(req.StartDate)
	if err != nil {
		RespondInvalidDate(c, "invalid start_date format")
		return
	}

	endDate, err := parseTimeValue(req.EndDate)
	if err != nil {
		RespondInvalidDate(c, "invalid end_date format")
		return
	}

	// 4) Build subscription model for service layer validation and persistence.
	subscription := &models.Subscription{
		MemberID:       memberID,
		CourseID:       courseID,
		HomeBranchID:   branchID,
		StartDate:      startDate,
		EndDate:        endDate,
		SessionPerWeek: req.SessionPerWeek,
		DiscountType:   req.DiscountType,
		DiscountValue:  req.DiscountValue,
		PromoCode:      req.PromoCode,
	}

	// 5) Delegate create logic to service.
	err = h.subscriptionService.CreateSubscription(c.Request.Context(), subscription)
	if err != nil {
		// 6) Map service validation/reference errors to HTTP status codes.
		switch {
		case errors.Is(err, service.ErrInvalidSubscriptionInput), errors.Is(err, service.ErrInvalidDiscount):
			RespondInvalidInput(c, err.Error())
		case errors.Is(err, service.ErrSubscriptionReferenceNotFound):
			RespondNotFound(c, err.Error())
		default:
			RespondInternal(c)
		}
		return
	}

	// 7) Success response includes created subscription payload.
	c.JSON(http.StatusCreated, gin.H{
		"message": "subscription created successfully",
		"data":    subscription,
	})
}

// GetByID fetches subscription by ID path param.
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	// 1) Validate subscription ID from URL.
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		RespondInvalidID(c, "invalid subscription id")
		return
	}

	// 2) Delegate lookup to service layer.
	subscription, err := h.subscriptionService.GetSubscriptionByID(c.Request.Context(), id)
	if err != nil {
		// 3) Not-found maps to 404, all other errors to 500.
		switch {
		case errors.Is(err, service.ErrSubscriptionNotFound):
			RespondNotFound(c, "subscription not found")
		default:
			RespondInternal(c)
		}
		return
	}

	// 4) Success response with subscription payload.
	c.JSON(http.StatusOK, gin.H{
		"message": "subscription fetched successfully",
		"data":    subscription,
	})
}

// Suspend moves subscription to suspended status with suspension details.
func (h *SubscriptionHandler) Suspend(c *gin.Context) {
	// 1) Validate subscription ID from URL.
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		RespondInvalidID(c, "invalid subscription id")
		return
	}

	// 2) Parse JSON body for suspension details.
	var req suspendSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondInvalidRequestBody(c)
		return
	}

	startDate, err := parseTimeValue(req.StartDate)
	if err != nil {
		RespondInvalidDate(c, "invalid start_date format")
		return
	}
	endDate, err := parseTimeValue(req.EndDate)
	if err != nil {
		RespondInvalidDate(c, "invalid end_date format")
		return
	}

	// 3) Build suspension model and delegate to service.
	suspension := &models.Suspension{
		StartDate:     startDate,
		EndDate:       endDate,
		FrozenSession: req.FrozenSession,
		Reason:        req.Reason,
	}

	if err := h.subscriptionService.SuspendSubscription(c.Request.Context(), id, suspension); err != nil {
		// 4) Map service errors to HTTP status codes.
		switch {
		case errors.Is(err, service.ErrSubscriptionNotFound):
			RespondNotFound(c, "subscription not found")
		case errors.Is(err, service.ErrInvalidSuspensionPeriod):
			RespondInvalidInput(c, err.Error())
		case errors.Is(err, service.ErrSubscriptionAlreadySuspended), errors.Is(err, service.ErrSubscriptionNotActive), errors.Is(err, service.ErrSubscriptionExpired):
			RespondConflict(c, err.Error())
		default:
			RespondInternal(c)
		}
		return
	}

	// 5) Success response.
	c.JSON(http.StatusOK, gin.H{"message": "subscription suspended successfully"})
}

// Resume clears suspension and sets status back to active.
func (h *SubscriptionHandler) Resume(c *gin.Context) {
	// 1) Validate subscription ID from URL.
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		RespondInvalidID(c, "invalid subscription id")
		return
	}

	// 2) Delegate to service.
	if err := h.subscriptionService.ResumeSubscription(c.Request.Context(), id); err != nil {
		// 3) Map service errors to HTTP status codes.
		switch {
		case errors.Is(err, service.ErrSubscriptionNotFound):
			RespondNotFound(c, "subscription not found")
		case errors.Is(err, service.ErrInvalidSubscriptionStatus), errors.Is(err, service.ErrSubscriptionExpired):
			RespondConflict(c, err.Error())
		default:
			RespondInternal(c)
		}
		return
	}

	// 4) Success response.
	c.JSON(http.StatusOK, gin.H{"message": "subscription resumed successfully"})
}

// Expire marks subscription as expired manually.
func (h *SubscriptionHandler) Expire(c *gin.Context) {
	// 1) Validate subscription ID from URL.
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		RespondInvalidID(c, "invalid subscription id")
		return
	}

	// 2) Delegate to service.
	if err := h.subscriptionService.ExpireSubscription(c.Request.Context(), id); err != nil {
		// 3) Map service errors to HTTP status codes.
		switch {
		case errors.Is(err, service.ErrSubscriptionNotFound):
			RespondNotFound(c, "subscription not found")
		default:
			RespondInternal(c)
		}
		return
	}

	// 4) Success response.
	c.JSON(http.StatusOK, gin.H{"message": "subscription expired successfully"})
}

// Refund refunds an active or suspended subscription by remaining sessions.
func (h *SubscriptionHandler) Refund(c *gin.Context) {
	// 1) Validate subscription ID from URL.
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		RespondInvalidID(c, "invalid subscription id")
		return
	}

	// 2) Parse optional JSON body.
	var req refundSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondInvalidRequestBody(c)
		return
	}

	refund, err := h.subscriptionService.RefundSubscription(c.Request.Context(), id, req.Reason)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSubscriptionInput):
			RespondInvalidInput(c, err.Error())
		case errors.Is(err, service.ErrSubscriptionNotFound):
			RespondNotFound(c, "subscription not found")
		case errors.Is(err, service.ErrSubscriptionCannotRefund), errors.Is(err, service.ErrSubscriptionNoRemaining), errors.Is(err, service.ErrRefundAlreadyExists):
			RespondConflict(c, err.Error())
		default:
			RespondInternal(c)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subscription refunded successfully",
		"refund":  refund,
	})
}

// parseTimeValue parses RFC3339 date input.
func parseTimeValue(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}
