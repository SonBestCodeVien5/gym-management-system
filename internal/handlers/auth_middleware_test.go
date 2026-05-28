package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
)

type stubAuthService struct {
	claims *service.AuthClaims
	err    error
}

func (s *stubAuthService) BootstrapAdmin(ctx context.Context, cfg service.BootstrapAdminConfig) error {
	return nil
}

func (s *stubAuthService) Login(ctx context.Context, email string, password string) (*service.AuthTokenPair, error) {
	return nil, nil
}

func (s *stubAuthService) Refresh(ctx context.Context, refreshToken string) (*service.AuthTokenPair, error) {
	return nil, nil
}

func (s *stubAuthService) Logout(ctx context.Context, refreshToken string) error {
	return nil
}

func (s *stubAuthService) CurrentEmployee(ctx context.Context, employeeID string) (*service.AuthEmployeeResponse, error) {
	return nil, nil
}

func (s *stubAuthService) ValidateAccessToken(ctx context.Context, accessToken string) (*service.AuthClaims, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.claims, nil
}

func TestAuthRequiredMissingTokenReturnsUnauthorized(t *testing.T) {
	router := authMiddlewareTestRouter(&stubAuthService{
		claims: &service.AuthClaims{EmployeeID: "employee-id", Role: []string{service.RoleAdmin}},
	}, service.RoleAdmin)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
	assertErrorResponse(t, res, string(ErrorCodeUnauthorized), "missing access token")
}

func TestAuthRequiredInvalidTokenReturnsUnauthorized(t *testing.T) {
	router := authMiddlewareTestRouter(&stubAuthService{err: service.ErrInvalidToken}, service.RoleAdmin)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
	assertErrorResponse(t, res, string(ErrorCodeUnauthorized), "invalid access token")
}

func TestRequireRolesAllowsMatchingRole(t *testing.T) {
	router := authMiddlewareTestRouter(&stubAuthService{
		claims: &service.AuthClaims{EmployeeID: "employee-id", Role: []string{service.RoleManager}},
	}, service.RoleAdmin, service.RoleManager)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
}

func TestRequireRolesRejectsNonMatchingRole(t *testing.T) {
	router := authMiddlewareTestRouter(&stubAuthService{
		claims: &service.AuthClaims{EmployeeID: "employee-id", Role: []string{service.RoleReceptionist}},
	}, service.RoleAdmin, service.RoleManager)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusForbidden)
	}
	assertErrorResponse(t, res, string(ErrorCodeForbidden), "forbidden")
}

func TestAuthRequiredUnexpectedServiceErrorReturnsServerError(t *testing.T) {
	router := authMiddlewareTestRouter(&stubAuthService{err: errors.New("storage failed")}, service.RoleAdmin)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusInternalServerError)
	}
	assertErrorResponse(t, res, string(ErrorCodeInternalError), "internal server error")
}

func authMiddlewareTestRouter(authService service.AuthService, allowedRoles ...string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", AuthRequired(authService), RequireRoles(allowedRoles...), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	return router
}

func assertErrorResponse(t *testing.T, res *httptest.ResponseRecorder, code string, message string) {
	t.Helper()

	var body struct {
		Error struct {
			Code    string         `json:"code"`
			Message string         `json:"message"`
			Details map[string]any `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not JSON: %v", err)
	}
	if body.Error.Code != code {
		t.Fatalf("error.code = %q, want %q", body.Error.Code, code)
	}
	if body.Error.Message != message {
		t.Fatalf("error.message = %q, want %q", body.Error.Message, message)
	}
	if body.Error.Details == nil {
		t.Fatalf("error.details is nil, want empty object")
	}
}
