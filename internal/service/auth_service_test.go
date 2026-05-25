package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type stubEmployeeAuthRepo struct {
	byID    map[string]*models.Employee
	byEmail map[string]*models.Employee
}

func (r *stubEmployeeAuthRepo) Create(ctx context.Context, employee *models.Employee) error {
	return nil
}

func (r *stubEmployeeAuthRepo) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	employee, ok := r.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return employee, nil
}

func (r *stubEmployeeAuthRepo) GetByNormalizedEmail(ctx context.Context, email string) (*models.Employee, error) {
	employee, ok := r.byEmail[email]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return employee, nil
}

func (r *stubEmployeeAuthRepo) BootstrapAdmin(ctx context.Context, employee *models.Employee) error {
	return nil
}

type stubRefreshTokenAuthRepo struct {
	created []*models.RefreshToken
	byHash  map[string]*models.RefreshToken
}

func (r *stubRefreshTokenAuthRepo) Create(ctx context.Context, token *models.RefreshToken) error {
	if r.byHash == nil {
		r.byHash = map[string]*models.RefreshToken{}
	}
	r.created = append(r.created, token)
	r.byHash[token.TokenHash] = token
	return nil
}

func (r *stubRefreshTokenAuthRepo) FindActiveByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	token, ok := r.byHash[tokenHash]
	if !ok || token.RevokedAt != nil {
		return nil, repository.ErrNotFound
	}
	return token, nil
}

func (r *stubRefreshTokenAuthRepo) RevokeActiveByHash(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	token, ok := r.byHash[tokenHash]
	if !ok || token.RevokedAt != nil {
		return repository.ErrNotFound
	}
	token.RevokedAt = &revokedAt
	token.UpdatedAt = revokedAt
	return nil
}

func TestAuthServiceLoginRefreshAndLogout(t *testing.T) {
	ctx := context.Background()
	fixedNow := time.Date(2026, 5, 25, 10, 0, 0, 0, time.UTC)
	employee := testAuthEmployee(t)
	employeeRepo := &stubEmployeeAuthRepo{
		byID:    map[string]*models.Employee{employee.ID.Hex(): employee},
		byEmail: map[string]*models.Employee{employee.NormalizedEmail: employee},
	}
	refreshRepo := &stubRefreshTokenAuthRepo{byHash: map[string]*models.RefreshToken{}}
	authService := newTestAuthService(t, employeeRepo, refreshRepo, fixedNow)

	pair, err := authService.Login(ctx, " ADMIN@GYM.TEST ", "secret")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("Login() returned empty token pair")
	}
	if pair.Employee == nil || pair.Employee.Email != employee.Email {
		t.Fatalf("Login() employee = %#v, want %s", pair.Employee, employee.Email)
	}
	if len(refreshRepo.created) != 1 {
		t.Fatalf("created refresh tokens = %d, want 1", len(refreshRepo.created))
	}

	refreshed, err := authService.Refresh(ctx, pair.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken == "" {
		t.Fatal("Refresh() returned empty token pair")
	}
	if len(refreshRepo.created) != 2 {
		t.Fatalf("created refresh tokens after refresh = %d, want 2", len(refreshRepo.created))
	}
	if refreshRepo.byHash[HashToken(pair.RefreshToken)].RevokedAt == nil {
		t.Fatal("Refresh() did not revoke old refresh token")
	}

	_, err = authService.Refresh(ctx, pair.RefreshToken)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Refresh() reused token error = %v, want %v", err, ErrInvalidToken)
	}

	if err := authService.Logout(ctx, refreshed.RefreshToken); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if err := authService.Logout(ctx, refreshed.RefreshToken); err != nil {
		t.Fatalf("Logout() repeated error = %v", err)
	}
}

func TestAuthServiceLoginWrongPassword(t *testing.T) {
	employee := testAuthEmployee(t)
	authService := newTestAuthService(
		t,
		&stubEmployeeAuthRepo{
			byID:    map[string]*models.Employee{employee.ID.Hex(): employee},
			byEmail: map[string]*models.Employee{employee.NormalizedEmail: employee},
		},
		&stubRefreshTokenAuthRepo{byHash: map[string]*models.RefreshToken{}},
		time.Date(2026, 5, 25, 10, 0, 0, 0, time.UTC),
	)

	_, err := authService.Login(context.Background(), employee.Email, "wrong")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func newTestAuthService(t *testing.T, employeeRepo *stubEmployeeAuthRepo, refreshRepo *stubRefreshTokenAuthRepo, now time.Time) AuthService {
	t.Helper()

	authService, err := NewAuthService(employeeRepo, refreshRepo, AuthConfig{
		AccessSecret:  "access-secret",
		RefreshSecret: "refresh-secret",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
	})
	if err != nil {
		t.Fatalf("NewAuthService() error = %v", err)
	}
	authService.(*authServiceImpl).now = func() time.Time { return now }
	return authService
}

func testAuthEmployee(t *testing.T) *models.Employee {
	t.Helper()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	email := "admin@gym.test"
	return &models.Employee{
		ID:              primitive.NewObjectID(),
		EmployeeID:      "EMP001",
		FullName:        "Gym Admin",
		Email:           email,
		NormalizedEmail: email,
		PasswordHash:    string(passwordHash),
		Status:          EmployeeStatusActive,
		Role:            []string{RoleAdmin},
		BranchID:        []primitive.ObjectID{},
	}
}
