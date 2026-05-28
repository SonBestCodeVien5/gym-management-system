package handlers

import (
	"errors"
	"net/http"

	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondInvalidRequestBody(c)
		return
	}

	tokenPair, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAuthInput):
			RespondInvalidInput(c, err.Error())
		case errors.Is(err, service.ErrInvalidCredentials), errors.Is(err, service.ErrInactiveEmployee):
			RespondUnauthorized(c, "invalid credentials")
		default:
			RespondInternal(c)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "data": tokenPair})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondInvalidRequestBody(c)
		return
	}
	if req.RefreshToken == "" {
		RespondInvalidInput(c, "refresh_token is required")
		return
	}

	tokenPair, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidToken), errors.Is(err, service.ErrInactiveEmployee):
			RespondUnauthorized(c, "invalid refresh token")
		default:
			RespondInternal(c)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token refreshed successfully", "data": tokenPair})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondInvalidRequestBody(c)
		return
	}
	if req.RefreshToken == "" {
		RespondInvalidInput(c, "refresh_token is required")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidToken):
			RespondUnauthorized(c, "invalid refresh token")
		default:
			RespondInternal(c)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
