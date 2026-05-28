package handlers

import (
	"errors"
	"strings"

	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	AuthEmployeeIDKey = "auth_employee_id"
	AuthRolesKey      = "auth_roles"
)

func AuthRequired(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			AbortUnauthorized(c, "missing access token")
			return
		}

		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			AbortUnauthorized(c, "invalid access token")
			return
		}

		claims, err := authService.ValidateAccessToken(c.Request.Context(), parts[1])
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidToken), errors.Is(err, service.ErrInactiveEmployee):
				AbortUnauthorized(c, "invalid access token")
			default:
				AbortInternal(c)
			}
			return
		}

		c.Set(AuthEmployeeIDKey, claims.EmployeeID)
		c.Set(AuthRolesKey, claims.Role)
		c.Next()
	}
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		rawRoles, exists := c.Get(AuthRolesKey)
		if !exists {
			AbortUnauthorized(c, "missing access token")
			return
		}

		roles, ok := rawRoles.([]string)
		if !ok {
			AbortUnauthorized(c, "invalid access token")
			return
		}

		for _, role := range roles {
			if _, ok := allowed[role]; ok {
				c.Next()
				return
			}
		}

		AbortForbidden(c, "forbidden")
	}
}
