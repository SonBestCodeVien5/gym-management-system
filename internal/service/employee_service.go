package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	EmployeeStatusInactive = "inactive"
)

var (
	ErrEmployeeNotFound     = errors.New("employee not found")
	ErrInvalidEmployeeInput = errors.New("invalid employee input")
	ErrEmployeeConflict     = errors.New("employee conflict")
)

type EmployeeCreateInput struct {
	EmployeeID string
	FullName   string
	Email      string
	Password   string
	Role       []string
	Level      string
	Phone      string
	BranchID   []primitive.ObjectID
	Status     string
}

type EmployeeUpdateInput struct {
	EmployeeID *string
	FullName   *string
	Email      *string
	Role       *[]string
	Level      *string
	Phone      *string
	BranchID   *[]primitive.ObjectID
	Status     *string
}

type EmployeeListInput struct {
	Role     string
	Status   string
	BranchID primitive.ObjectID
}

type EmployeeResponse struct {
	ID         primitive.ObjectID   `json:"id"`
	EmployeeID string               `json:"employee_id"`
	FullName   string               `json:"full_name"`
	Email      string               `json:"email"`
	Status     string               `json:"status"`
	Role       []string             `json:"role"`
	Level      string               `json:"level"`
	Phone      string               `json:"phone"`
	BranchID   []primitive.ObjectID `json:"branch_id"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

type EmployeeService interface {
	CreateEmployee(ctx context.Context, input EmployeeCreateInput) (*EmployeeResponse, error)
	ListEmployees(ctx context.Context, input EmployeeListInput) ([]EmployeeResponse, error)
	GetEmployeeByID(ctx context.Context, id string) (*EmployeeResponse, error)
	UpdateEmployee(ctx context.Context, actorID string, id string, input EmployeeUpdateInput) (*EmployeeResponse, error)
	UpdatePassword(ctx context.Context, id string, password string) error
}

type employeeServiceImpl struct {
	employeeRepo     repository.EmployeeRepository
	branchRepo       repository.BranchRepository
	refreshTokenRepo repository.RefreshTokenRepository
	now              func() time.Time
}

func NewEmployeeService(
	employeeRepo repository.EmployeeRepository,
	branchRepo repository.BranchRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) EmployeeService {
	return &employeeServiceImpl{
		employeeRepo:     employeeRepo,
		branchRepo:       branchRepo,
		refreshTokenRepo: refreshTokenRepo,
		now:              time.Now,
	}
}

func (s *employeeServiceImpl) CreateEmployee(ctx context.Context, input EmployeeCreateInput) (*EmployeeResponse, error) {
	employeeID := strings.TrimSpace(input.EmployeeID)
	fullName := strings.TrimSpace(input.FullName)
	normalizedEmail := NormalizeEmail(input.Email)
	if employeeID == "" || fullName == "" || normalizedEmail == "" {
		return nil, ErrInvalidEmployeeInput
	}

	roles, err := normalizeEmployeeRoles(input.Role)
	if err != nil {
		return nil, err
	}
	status, err := normalizeEmployeeStatus(input.Status, EmployeeStatusActive)
	if err != nil {
		return nil, err
	}
	level, err := normalizeEmployeeLevel(input.Level, hasEmployeeRole(roles, RoleTrainer))
	if err != nil {
		return nil, err
	}
	branches, err := s.validateBranches(ctx, input.BranchID)
	if err != nil {
		return nil, err
	}
	passwordHash, err := hashEmployeePassword(input.Password)
	if err != nil {
		return nil, err
	}

	now := s.now()
	employee := &models.Employee{
		ID:              primitive.NewObjectID(),
		EmployeeID:      employeeID,
		FullName:        fullName,
		Email:           normalizedEmail,
		NormalizedEmail: normalizedEmail,
		PasswordHash:    passwordHash,
		Status:          status,
		Role:            roles,
		Level:           level,
		Phone:           strings.TrimSpace(input.Phone),
		BranchID:        branches,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.employeeRepo.Create(ctx, employee); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, ErrEmployeeConflict
		}
		return nil, err
	}
	return employeeResponseForManagement(employee), nil
}

func (s *employeeServiceImpl) ListEmployees(ctx context.Context, input EmployeeListInput) ([]EmployeeResponse, error) {
	role := ""
	if strings.TrimSpace(input.Role) != "" {
		roles, err := normalizeEmployeeRoles([]string{input.Role})
		if err != nil {
			return nil, err
		}
		role = roles[0]
	}

	status := ""
	if strings.TrimSpace(input.Status) != "" {
		normalized, err := normalizeEmployeeStatus(input.Status, "")
		if err != nil {
			return nil, err
		}
		status = normalized
	}

	employees, err := s.employeeRepo.List(ctx, repository.EmployeeListFilter{
		Role:     role,
		Status:   status,
		BranchID: input.BranchID,
	})
	if err != nil {
		return nil, err
	}

	responses := make([]EmployeeResponse, 0, len(employees))
	for i := range employees {
		responses = append(responses, *employeeResponseForManagement(&employees[i]))
	}
	return responses, nil
}

func (s *employeeServiceImpl) GetEmployeeByID(ctx context.Context, id string) (*EmployeeResponse, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, ErrInvalidEmployeeInput
	}

	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}
	return employeeResponseForManagement(employee), nil
}

func (s *employeeServiceImpl) UpdateEmployee(ctx context.Context, actorID string, id string, input EmployeeUpdateInput) (*EmployeeResponse, error) {
	targetID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidEmployeeInput
	}

	current, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	update := repository.EmployeeUpdate{UpdatedAt: s.now()}
	mutableFieldSupplied := false

	mergedEmployeeID := current.EmployeeID
	if input.EmployeeID != nil {
		value := strings.TrimSpace(*input.EmployeeID)
		if value == "" {
			return nil, ErrInvalidEmployeeInput
		}
		update.EmployeeID = &value
		mergedEmployeeID = value
		mutableFieldSupplied = true
	}

	mergedFullName := current.FullName
	if input.FullName != nil {
		value := strings.TrimSpace(*input.FullName)
		if value == "" {
			return nil, ErrInvalidEmployeeInput
		}
		update.FullName = &value
		mergedFullName = value
		mutableFieldSupplied = true
	}

	mergedEmail := current.Email
	if input.Email != nil {
		normalizedEmail := NormalizeEmail(*input.Email)
		if normalizedEmail == "" {
			return nil, ErrInvalidEmployeeInput
		}
		update.Email = &normalizedEmail
		update.NormalizedEmail = &normalizedEmail
		mergedEmail = normalizedEmail
		mutableFieldSupplied = true
	}

	mergedRoles := append([]string(nil), current.Role...)
	if input.Role != nil {
		roles, err := normalizeEmployeeRoles(*input.Role)
		if err != nil {
			return nil, err
		}
		update.Role = &roles
		mergedRoles = roles
		mutableFieldSupplied = true
	}

	mergedLevel := current.Level
	if input.Level != nil {
		level, err := normalizeEmployeeLevel(*input.Level, hasEmployeeRole(mergedRoles, RoleTrainer))
		if err != nil {
			return nil, err
		}
		update.Level = &level
		mergedLevel = level
		mutableFieldSupplied = true
	}

	if hasEmployeeRole(mergedRoles, RoleTrainer) && strings.TrimSpace(mergedLevel) == "" {
		return nil, ErrInvalidEmployeeInput
	}

	mergedPhone := current.Phone
	if input.Phone != nil {
		value := strings.TrimSpace(*input.Phone)
		update.Phone = &value
		mergedPhone = value
		mutableFieldSupplied = true
	}

	mergedBranches := append([]primitive.ObjectID(nil), current.BranchID...)
	if input.BranchID != nil {
		branches, err := s.validateBranches(ctx, *input.BranchID)
		if err != nil {
			return nil, err
		}
		update.BranchID = &branches
		mergedBranches = branches
		mutableFieldSupplied = true
	}

	mergedStatus := current.Status
	if input.Status != nil {
		status, err := normalizeEmployeeStatus(*input.Status, "")
		if err != nil {
			return nil, err
		}
		update.Status = &status
		mergedStatus = status
		mutableFieldSupplied = true
	}

	if !mutableFieldSupplied {
		return nil, ErrInvalidEmployeeInput
	}
	if actorID == targetID.Hex() {
		if mergedStatus == EmployeeStatusInactive {
			return nil, ErrEmployeeConflict
		}
		if hasEmployeeRole(current.Role, RoleAdmin) && !hasEmployeeRole(mergedRoles, RoleAdmin) {
			return nil, ErrEmployeeConflict
		}
	}

	updated, err := s.employeeRepo.UpdateByID(ctx, id, update)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrEmployeeNotFound
		}
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, ErrEmployeeConflict
		}
		return nil, err
	}

	if current.Status == EmployeeStatusActive && mergedStatus == EmployeeStatusInactive {
		if err := s.refreshTokenRepo.RevokeActiveByEmployeeID(ctx, targetID, s.now()); err != nil {
			return nil, err
		}
	}

	updated.EmployeeID = mergedEmployeeID
	updated.FullName = mergedFullName
	updated.Email = mergedEmail
	updated.Role = mergedRoles
	updated.Level = mergedLevel
	updated.Phone = mergedPhone
	updated.BranchID = mergedBranches
	updated.Status = mergedStatus
	return employeeResponseForManagement(updated), nil
}

func (s *employeeServiceImpl) UpdatePassword(ctx context.Context, id string, password string) error {
	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidEmployeeInput
	}

	passwordHash, err := hashEmployeePassword(password)
	if err != nil {
		return err
	}

	now := s.now()
	if err := s.employeeRepo.UpdatePasswordByID(ctx, id, passwordHash, now); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrEmployeeNotFound
		}
		return err
	}
	if err := s.refreshTokenRepo.RevokeActiveByEmployeeID(ctx, employeeID, now); err != nil {
		return err
	}
	return nil
}

func (s *employeeServiceImpl) validateBranches(ctx context.Context, branchIDs []primitive.ObjectID) ([]primitive.ObjectID, error) {
	if len(branchIDs) == 0 {
		return []primitive.ObjectID{}, nil
	}

	seen := make(map[primitive.ObjectID]struct{}, len(branchIDs))
	branches := make([]primitive.ObjectID, 0, len(branchIDs))
	for _, branchID := range branchIDs {
		if branchID.IsZero() {
			return nil, ErrInvalidEmployeeInput
		}
		if _, ok := seen[branchID]; ok {
			continue
		}
		if _, err := s.branchRepo.GetByID(ctx, branchID.Hex()); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrInvalidEmployeeInput
			}
			return nil, err
		}
		seen[branchID] = struct{}{}
		branches = append(branches, branchID)
	}
	return branches, nil
}

func normalizeEmployeeRoles(input []string) ([]string, error) {
	if len(input) == 0 {
		return nil, ErrInvalidEmployeeInput
	}

	seen := map[string]struct{}{}
	roles := make([]string, 0, len(input))
	for _, rawRole := range input {
		role := strings.ToLower(strings.TrimSpace(rawRole))
		if !isAllowedEmployeeRole(role) {
			return nil, ErrInvalidEmployeeInput
		}
		if _, ok := seen[role]; ok {
			continue
		}
		seen[role] = struct{}{}
		roles = append(roles, role)
	}
	if len(roles) == 0 {
		return nil, ErrInvalidEmployeeInput
	}
	return roles, nil
}

func isAllowedEmployeeRole(role string) bool {
	switch role {
	case RoleAdmin, RoleManager, RoleTrainer, RoleReceptionist:
		return true
	default:
		return false
	}
}

func hasEmployeeRole(roles []string, expected string) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}

func normalizeEmployeeStatus(input string, fallback string) (string, error) {
	status := strings.ToLower(strings.TrimSpace(input))
	if status == "" {
		if fallback != "" {
			return fallback, nil
		}
		return "", ErrInvalidEmployeeInput
	}

	switch status {
	case EmployeeStatusActive, EmployeeStatusInactive:
		return status, nil
	default:
		return "", ErrInvalidEmployeeInput
	}
}

func normalizeEmployeeLevel(input string, required bool) (string, error) {
	level := strings.ToLower(strings.TrimSpace(input))
	if level == "" {
		if required {
			return "", ErrInvalidEmployeeInput
		}
		return "", nil
	}
	if !isAllowedCourseLevel(level) {
		return "", ErrInvalidEmployeeInput
	}
	return level, nil
}

func hashEmployeePassword(password string) (string, error) {
	if len(password) < 8 || strings.TrimSpace(password) == "" {
		return "", ErrInvalidEmployeeInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func employeeResponseForManagement(employee *models.Employee) *EmployeeResponse {
	return &EmployeeResponse{
		ID:         employee.ID,
		EmployeeID: employee.EmployeeID,
		FullName:   employee.FullName,
		Email:      employee.Email,
		Status:     employee.Status,
		Role:       employee.Role,
		Level:      employee.Level,
		Phone:      employee.Phone,
		BranchID:   employee.BranchID,
		CreatedAt:  employee.CreatedAt,
		UpdatedAt:  employee.UpdatedAt,
	}
}
