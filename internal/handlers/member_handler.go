package handlers

import (
	"errors"
	"net/http"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService service.MemberService
}

func NewMemberHandler(memberService service.MemberService) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
	}
}

type registerMemberRequest struct {
	CCID     string `json:"ccid"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
	Level    string `json:"level"`
}

func (h *MemberHandler) Register(c *gin.Context) {
	var req registerMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	member := &models.Member{
		CCID:     req.CCID,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Level:    req.Level,
	}

	err := h.memberService.RegisterMember(c.Request.Context(), member)
	if err != nil {
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "member registered successfully",
		"data":    member,
	})
}
