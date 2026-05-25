package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	EmployeeStatusActive = "active"
	RoleAdmin            = "admin"
	RoleManager          = "manager"
	RoleTrainer          = "trainer"
	RoleReceptionist     = "receptionist"
)

var (
	ErrInvalidAuthInput   = errors.New("invalid auth input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInactiveEmployee   = errors.New("employee is inactive")
	ErrAuthConfig         = errors.New("invalid auth config")
)

type AuthConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type BootstrapAdminConfig struct {
	EmployeeID string
	FullName   string
	Email      string
	Password   string
	Phone      string
	Level      string
}

type AuthEmployeeResponse struct {
	ID         primitive.ObjectID   `json:"id"`
	EmployeeID string               `json:"employee_id"`
	Email      string               `json:"email"`
	FullName   string               `json:"full_name"`
	Role       []string             `json:"role"`
	BranchID   []primitive.ObjectID `json:"branch_id"`
}

type AuthTokenPair struct {
	AccessToken  string                `json:"access_token"`
	RefreshToken string                `json:"refresh_token"`
	Employee     *AuthEmployeeResponse `json:"employee,omitempty"`
}

type AuthClaims struct {
	EmployeeID string   `json:"employee_id"`
	Role       []string `json:"role"`
	TokenType  string   `json:"token_type"`
	TokenID    string   `json:"jti"`
	ExpiresAt  int64    `json:"exp"`
	IssuedAt   int64    `json:"iat"`
}

type AuthService interface {
	BootstrapAdmin(ctx context.Context, cfg BootstrapAdminConfig) error
	Login(ctx context.Context, email string, password string) (*AuthTokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthTokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateAccessToken(ctx context.Context, accessToken string) (*AuthClaims, error)
}

type authServiceImpl struct {
	employeeRepo     repository.EmployeeRepository
	refreshTokenRepo repository.RefreshTokenRepository
	config           AuthConfig
	now              func() time.Time
}

func NewAuthService(employeeRepo repository.EmployeeRepository, refreshTokenRepo repository.RefreshTokenRepository, config AuthConfig) (AuthService, error) {
	if strings.TrimSpace(config.AccessSecret) == "" || strings.TrimSpace(config.RefreshSecret) == "" {
		return nil, ErrAuthConfig
	}
	if config.AccessTTL <= 0 || config.RefreshTTL <= 0 {
		return nil, ErrAuthConfig
	}

	return &authServiceImpl{
		employeeRepo:     employeeRepo,
		refreshTokenRepo: refreshTokenRepo,
		config:           config,
		now:              time.Now,
	}, nil
}

func (s *authServiceImpl) BootstrapAdmin(ctx context.Context, cfg BootstrapAdminConfig) error {
	if strings.TrimSpace(cfg.Email) == "" || strings.TrimSpace(cfg.Password) == "" {
		return nil
	}

	employeeID := strings.TrimSpace(cfg.EmployeeID)
	if employeeID == "" {
		employeeID = "ADMIN001"
	}
	fullName := strings.TrimSpace(cfg.FullName)
	if fullName == "" {
		fullName = "Gym Admin"
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	normalizedEmail := NormalizeEmail(cfg.Email)
	admin := &models.Employee{
		ID:              primitive.NewObjectID(),
		EmployeeID:      employeeID,
		FullName:        fullName,
		Email:           normalizedEmail,
		NormalizedEmail: normalizedEmail,
		PasswordHash:    string(passwordHash),
		Status:          EmployeeStatusActive,
		Role:            []string{RoleAdmin},
		Level:           strings.TrimSpace(cfg.Level),
		Phone:           strings.TrimSpace(cfg.Phone),
		BranchID:        []primitive.ObjectID{},
	}

	return s.employeeRepo.BootstrapAdmin(ctx, admin)
}

func (s *authServiceImpl) Login(ctx context.Context, email string, password string) (*AuthTokenPair, error) {
	normalizedEmail := NormalizeEmail(email)
	if normalizedEmail == "" || password == "" {
		return nil, ErrInvalidAuthInput
	}

	employee, err := s.employeeRepo.GetByNormalizedEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if employee.Status != EmployeeStatusActive {
		return nil, ErrInactiveEmployee
	}
	if bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokenPair(ctx, employee, true)
}

func (s *authServiceImpl) Refresh(ctx context.Context, refreshToken string) (*AuthTokenPair, error) {
	claims, err := s.parseToken(refreshToken, s.config.RefreshSecret, "refresh", true)
	if err != nil {
		return nil, ErrInvalidToken
	}

	tokenHash := HashToken(refreshToken)
	storedToken, err := s.refreshTokenRepo.FindActiveByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	if storedToken.ExpiresAt.Before(s.now()) {
		return nil, ErrInvalidToken
	}
	if storedToken.EmployeeID.Hex() != claims.EmployeeID {
		return nil, ErrInvalidToken
	}

	employee, err := s.employeeRepo.GetByID(ctx, claims.EmployeeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	if employee.Status != EmployeeStatusActive {
		return nil, ErrInactiveEmployee
	}

	if err := s.refreshTokenRepo.RevokeActiveByHash(ctx, tokenHash, s.now()); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	return s.issueTokenPair(ctx, employee, false)
}

func (s *authServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	if _, err := s.parseToken(refreshToken, s.config.RefreshSecret, "refresh", false); err != nil {
		return ErrInvalidToken
	}

	if err := s.refreshTokenRepo.RevokeActiveByHash(ctx, HashToken(refreshToken), s.now()); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (s *authServiceImpl) ValidateAccessToken(ctx context.Context, accessToken string) (*AuthClaims, error) {
	claims, err := s.parseToken(accessToken, s.config.AccessSecret, "access", true)
	if err != nil {
		return nil, ErrInvalidToken
	}

	employee, err := s.employeeRepo.GetByID(ctx, claims.EmployeeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	if employee.Status != EmployeeStatusActive {
		return nil, ErrInactiveEmployee
	}

	claims.Role = employee.Role
	return claims, nil
}

func (s *authServiceImpl) issueTokenPair(ctx context.Context, employee *models.Employee, includeEmployee bool) (*AuthTokenPair, error) {
	accessToken, err := s.createToken(employee, "access", s.config.AccessSecret, s.config.AccessTTL)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.createToken(employee, "refresh", s.config.RefreshSecret, s.config.RefreshTTL)
	if err != nil {
		return nil, err
	}

	now := s.now()
	refreshRecord := &models.RefreshToken{
		ID:         primitive.NewObjectID(),
		EmployeeID: employee.ID,
		TokenHash:  HashToken(refreshToken),
		ExpiresAt:  now.Add(s.config.RefreshTTL),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.refreshTokenRepo.Create(ctx, refreshRecord); err != nil {
		return nil, err
	}

	pair := &AuthTokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	if includeEmployee {
		pair.Employee = employeeResponse(employee)
	}
	return pair, nil
}

func (s *authServiceImpl) createToken(employee *models.Employee, tokenType string, secret string, ttl time.Duration) (string, error) {
	now := s.now()
	tokenID, err := randomTokenID()
	if err != nil {
		return "", err
	}
	claims := AuthClaims{
		EmployeeID: employee.ID.Hex(),
		Role:       employee.Role,
		TokenType:  tokenType,
		TokenID:    tokenID,
		ExpiresAt:  now.Add(ttl).Unix(),
		IssuedAt:   now.Unix(),
	}

	headerBytes, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsBytes)
	unsigned := encodedHeader + "." + encodedClaims
	signature := signJWT(unsigned, secret)
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

func (s *authServiceImpl) parseToken(token string, secret string, expectedType string, checkExpiry bool) (*AuthClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}
	expectedSignature := signJWT(unsigned, secret)
	if !hmac.Equal(signature, expectedSignature) {
		return nil, ErrInvalidToken
	}

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, ErrInvalidToken
	}
	if header.Alg != "HS256" || header.Typ != "JWT" {
		return nil, ErrInvalidToken
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}
	var claims AuthClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, ErrInvalidToken
	}
	if claims.TokenType != expectedType || claims.EmployeeID == "" {
		return nil, ErrInvalidToken
	}
	if checkExpiry && claims.ExpiresAt <= s.now().Unix() {
		return nil, ErrInvalidToken
	}
	return &claims, nil
}

func signJWT(unsigned string, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsigned))
	return mac.Sum(nil)
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func employeeResponse(employee *models.Employee) *AuthEmployeeResponse {
	return &AuthEmployeeResponse{
		ID:         employee.ID,
		EmployeeID: employee.EmployeeID,
		Email:      employee.Email,
		FullName:   employee.FullName,
		Role:       employee.Role,
		BranchID:   employee.BranchID,
	}
}

func RandomSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate secret: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func randomTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate token id: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
