package handlers

import (
	"errors"
	"net/http"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberHandler struct {
	memberService       service.MemberService
	subscriptionService service.SubscriptionService
}

// NewMemberHandler wires member and subscription services for activate flow.
func NewMemberHandler(memberService service.MemberService, subscriptionService service.SubscriptionService) *MemberHandler {
	return &MemberHandler{
		memberService:       memberService,
		subscriptionService: subscriptionService,
	}
}

// registerMemberRequest is the JSON body for member registration.
type registerMemberRequest struct {
	CCID     string `json:"ccid"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
	Level    string `json:"level"`
}

// activateMemberRequest is the JSON body for offline payment confirmation.
type activateMemberRequest struct {
	SubscriptionID string `json:"subscription_id"`
}

// Register creates a new member record.
func (h *MemberHandler) Register(c *gin.Context) {
	// 1) Parse JSON body into request struct.
	var req registerMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// 2) Map request fields into domain model (service will fill defaults).
	member := &models.Member{
		CCID:     req.CCID,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Level:    req.Level,
	}

	// 3) Delegate to service for validation + persistence.
	err := h.memberService.RegisterMember(c.Request.Context(), member)
	if err != nil {
		// 4) Translate service errors into HTTP status codes.
		switch {
		case errors.Is(err, service.ErrInvalidMemberInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case errors.Is(err, service.ErrMemberCCIDAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	// 5) Success response includes created member payload.
	c.JSON(http.StatusCreated, gin.H{
		"message": "member registered successfully",
		"data":    member,
	})
}

// GetByID fetches a member by ID path param.
func (h *MemberHandler) GetByID(c *gin.Context) {
	// 1) Validate member ID from URL.
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid member id"})
		return
	}

	// 2) Delegate to service for lookup and domain-level errors.
	member, err := h.memberService.GetMemberByID(c.Request.Context(), id)
	if err != nil {
		// 3) Map not-found to 404, everything else to 500.
		switch {
		case errors.Is(err, service.ErrMemberNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "member not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	// 4) Success response with member payload.
	c.JSON(http.StatusOK, gin.H{
		"message": "member fetched successfully",
		"data":    member,
	})
}

// Activate confirms subscription payment then marks member as registered.
func (h *MemberHandler) Activate(c *gin.Context) {
	// Validate member ID from path param.
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid member id"})
		return
	}

	// 1) Body must include subscription_id to link payment to this member.
	var req activateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}
	// 2) Validate subscription_id presence and format.
	if req.SubscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "subscription_id is required"})
		return
	}
	if _, err := primitive.ObjectIDFromHex(req.SubscriptionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid subscription id"})
		return
	}

	// 3) Confirm subscription payment before activating the member.
	if err := h.subscriptionService.ConfirmSubscriptionPayment(c.Request.Context(), id, req.SubscriptionID); err != nil {
		// 4) Subscription errors are reported as not-found or conflict.
		switch {
		case errors.Is(err, service.ErrSubscriptionNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "subscription not found"})
		case errors.Is(err, service.ErrSubscriptionAlreadyActive), errors.Is(err, service.ErrInvalidSubscriptionStatus):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		case errors.Is(err, service.ErrSubscriptionMemberMismatch):
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	// 5) Mark member as registered after payment is confirmed.
	if err := h.memberService.ActivateMember(c.Request.Context(), id); err != nil {
		// 6) Not-found means member ID does not exist.
		switch {
		case errors.Is(err, service.ErrMemberNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "member not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	// 7) Final success response.
	c.JSON(http.StatusOK, gin.H{
		"message": "member activated successfully",
	})
}
