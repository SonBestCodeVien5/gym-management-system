package handlers

import (
	"errors"
	"net/http"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *MemberHandler) GetByID(c *gin.Context) {
    id := c.Param("id")
    if _, err := primitive.ObjectIDFromHex(id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid member id"})
        return
    }

    member, err := h.memberService.GetMemberByID(c.Request.Context(), id)
    if err != nil {
        switch {
        case errors.Is(err, repository.ErrNotFound):
            c.JSON(http.StatusNotFound, gin.H{"message": "member not found"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "member fetched successfully",
        "data":    member,
    })
}