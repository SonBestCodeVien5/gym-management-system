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

type stubEmployeeManagementRepo struct {
	byID       map[string]*models.Employee
	created    *models.Employee
	createErr  error
	updateErr  error
	passwordID string
}

func (r *stubEmployeeManagementRepo) Create(ctx context.Context, employee *models.Employee) error {
	r.created = employee
	if r.byID == nil {
		r.byID = map[string]*models.Employee{}
	}
	r.byID[employee.ID.Hex()] = employee
	return r.createErr
}

func (r *stubEmployeeManagementRepo) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	employee, ok := r.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return employee, nil
}

func (r *stubEmployeeManagementRepo) GetByNormalizedEmail(ctx context.Context, email string) (*models.Employee, error) {
	for _, employee := range r.byID {
		if employee.NormalizedEmail == email {
			return employee, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (r *stubEmployeeManagementRepo) List(ctx context.Context, filter repository.EmployeeListFilter) ([]models.Employee, error) {
	employees := make([]models.Employee, 0, len(r.byID))
	for _, employee := range r.byID {
		employees = append(employees, *employee)
	}
	return employees, nil
}

func (r *stubEmployeeManagementRepo) UpdateByID(ctx context.Context, id string, update repository.EmployeeUpdate) (*models.Employee, error) {
	if r.updateErr != nil {
		return nil, r.updateErr
	}
	employee, ok := r.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	if update.EmployeeID != nil {
		employee.EmployeeID = *update.EmployeeID
	}
	if update.FullName != nil {
		employee.FullName = *update.FullName
	}
	if update.Email != nil {
		employee.Email = *update.Email
	}
	if update.NormalizedEmail != nil {
		employee.NormalizedEmail = *update.NormalizedEmail
	}
	if update.Status != nil {
		employee.Status = *update.Status
	}
	if update.Role != nil {
		employee.Role = *update.Role
	}
	if update.Level != nil {
		employee.Level = *update.Level
	}
	if update.Phone != nil {
		employee.Phone = *update.Phone
	}
	if update.BranchID != nil {
		employee.BranchID = *update.BranchID
	}
	employee.UpdatedAt = update.UpdatedAt
	return employee, nil
}

func (r *stubEmployeeManagementRepo) UpdatePasswordByID(ctx context.Context, id string, passwordHash string, updatedAt time.Time) error {
	employee, ok := r.byID[id]
	if !ok {
		return repository.ErrNotFound
	}
	employee.PasswordHash = passwordHash
	employee.UpdatedAt = updatedAt
	r.passwordID = id
	return nil
}

func (r *stubEmployeeManagementRepo) BootstrapAdmin(ctx context.Context, employee *models.Employee) error {
	return nil
}

type stubRefreshTokenManagementRepo struct {
	revokedEmployeeIDs []primitive.ObjectID
}

func (r *stubRefreshTokenManagementRepo) Create(ctx context.Context, token *models.RefreshToken) error {
	return nil
}

func (r *stubRefreshTokenManagementRepo) FindActiveByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	return nil, repository.ErrNotFound
}

func (r *stubRefreshTokenManagementRepo) RevokeActiveByHash(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	return nil
}

func (r *stubRefreshTokenManagementRepo) RevokeActiveByEmployeeID(ctx context.Context, employeeID primitive.ObjectID, revokedAt time.Time) error {
	r.revokedEmployeeIDs = append(r.revokedEmployeeIDs, employeeID)
	return nil
}

func TestEmployeeServiceCreateEmployee(t *testing.T) {
	branchID := primitive.NewObjectID()
	employeeRepo := &stubEmployeeManagementRepo{}
	service := NewEmployeeService(
		employeeRepo,
		&stubBranchRepo{branch: &models.Branch{ID: branchID}},
		&stubRefreshTokenManagementRepo{},
	)
	service.(*employeeServiceImpl).now = func() time.Time {
		return time.Date(2026, 5, 26, 8, 0, 0, 0, time.UTC)
	}

	response, err := service.CreateEmployee(context.Background(), EmployeeCreateInput{
		EmployeeID: " EMP001 ",
		FullName:   " Tran Van Trainer ",
		Email:      " Trainer@Gym.Test ",
		Password:   "strong-password-123",
		Role:       []string{" Trainer ", "trainer"},
		Level:      " Advanced ",
		Phone:      " 0900000002 ",
		BranchID:   []primitive.ObjectID{branchID},
	})
	if err != nil {
		t.Fatalf("CreateEmployee() error = %v", err)
	}
	if response.Email != "trainer@gym.test" {
		t.Fatalf("Email = %q, want normalized email", response.Email)
	}
	if response.EmployeeID != "EMP001" || response.FullName != "Tran Van Trainer" {
		t.Fatalf("response = %#v, want trimmed identity fields", response)
	}
	if len(response.Role) != 1 || response.Role[0] != RoleTrainer {
		t.Fatalf("Role = %#v, want deduplicated trainer", response.Role)
	}
	if employeeRepo.created == nil {
		t.Fatal("CreateEmployee() did not persist employee")
	}
	if employeeRepo.created.PasswordHash == "" {
		t.Fatal("CreateEmployee() did not hash password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(employeeRepo.created.PasswordHash), []byte("strong-password-123")); err != nil {
		t.Fatalf("stored password hash does not match password: %v", err)
	}
	if employeeRepo.created.NormalizedEmail != "trainer@gym.test" {
		t.Fatalf("NormalizedEmail = %q, want trainer@gym.test", employeeRepo.created.NormalizedEmail)
	}
}

func TestEmployeeServiceCreateTrainerRequiresLevel(t *testing.T) {
	service := NewEmployeeService(
		&stubEmployeeManagementRepo{},
		&stubBranchRepo{},
		&stubRefreshTokenManagementRepo{},
	)

	_, err := service.CreateEmployee(context.Background(), EmployeeCreateInput{
		EmployeeID: "EMP001",
		FullName:   "Tran Van Trainer",
		Email:      "trainer@gym.test",
		Password:   "strong-password-123",
		Role:       []string{RoleTrainer},
	})
	if !errors.Is(err, ErrInvalidEmployeeInput) {
		t.Fatalf("CreateEmployee() error = %v, want %v", err, ErrInvalidEmployeeInput)
	}
}

func TestEmployeeServiceUpdateSelfLockout(t *testing.T) {
	employee := testManagedEmployee([]string{RoleAdmin}, EmployeeStatusActive)
	employeeRepo := &stubEmployeeManagementRepo{
		byID: map[string]*models.Employee{employee.ID.Hex(): employee},
	}
	service := NewEmployeeService(
		employeeRepo,
		&stubBranchRepo{},
		&stubRefreshTokenManagementRepo{},
	)

	inactive := EmployeeStatusInactive
	_, err := service.UpdateEmployee(context.Background(), employee.ID.Hex(), employee.ID.Hex(), EmployeeUpdateInput{
		Status: &inactive,
	})
	if !errors.Is(err, ErrEmployeeConflict) {
		t.Fatalf("UpdateEmployee() self-deactivation error = %v, want %v", err, ErrEmployeeConflict)
	}

	roles := []string{RoleManager}
	_, err = service.UpdateEmployee(context.Background(), employee.ID.Hex(), employee.ID.Hex(), EmployeeUpdateInput{
		Role: &roles,
	})
	if !errors.Is(err, ErrEmployeeConflict) {
		t.Fatalf("UpdateEmployee() self-admin removal error = %v, want %v", err, ErrEmployeeConflict)
	}
}

func TestEmployeeServiceUpdatePasswordRevokesRefreshTokens(t *testing.T) {
	employee := testManagedEmployee([]string{RoleReceptionist}, EmployeeStatusActive)
	employeeRepo := &stubEmployeeManagementRepo{
		byID: map[string]*models.Employee{employee.ID.Hex(): employee},
	}
	refreshRepo := &stubRefreshTokenManagementRepo{}
	service := NewEmployeeService(
		employeeRepo,
		&stubBranchRepo{},
		refreshRepo,
	)

	if err := service.UpdatePassword(context.Background(), employee.ID.Hex(), "new-strong-password-123"); err != nil {
		t.Fatalf("UpdatePassword() error = %v", err)
	}
	if employeeRepo.passwordID != employee.ID.Hex() {
		t.Fatalf("passwordID = %q, want %s", employeeRepo.passwordID, employee.ID.Hex())
	}
	if err := bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte("new-strong-password-123")); err != nil {
		t.Fatalf("updated password hash does not match: %v", err)
	}
	if len(refreshRepo.revokedEmployeeIDs) != 1 || refreshRepo.revokedEmployeeIDs[0] != employee.ID {
		t.Fatalf("revokedEmployeeIDs = %#v, want employee ID", refreshRepo.revokedEmployeeIDs)
	}
}

func testManagedEmployee(roles []string, status string) *models.Employee {
	email := "staff@gym.test"
	return &models.Employee{
		ID:              primitive.NewObjectID(),
		EmployeeID:      "EMP001",
		FullName:        "Staff User",
		Email:           email,
		NormalizedEmail: email,
		PasswordHash:    "$2a$10$placeholderplaceholderplaceholderplaceholder",
		Status:          status,
		Role:            roles,
		Level:           "advanced",
		BranchID:        []primitive.ObjectID{},
	}
}
