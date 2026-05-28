package service

import (
	"context"
	"errors"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrBranchNotFound          = errors.New("branch not found")
	ErrInvalidBranchInput      = errors.New("invalid branch input")
	ErrBranchCodeAlreadyExists = errors.New("branch code already exists")
)

// BranchService defines business operations for branch management.
type BranchService interface {
	CreateBranch(ctx context.Context, branch *models.Branch) error
	GetBranchByID(ctx context.Context, id string) (*models.Branch, error)
	ListBranches(ctx context.Context) ([]models.Branch, error)
	UpdateBranch(ctx context.Context, id string, branch *models.Branch) error
	DeleteBranch(ctx context.Context, id string) error
	NearbyBranches(ctx context.Context, lng float64, lat float64, maxDistance int64, limit int64) ([]models.BranchNearbyResult, error)
}

type branchServiceImpl struct {
	repo repository.BranchRepository
}

// NewBranchService builds the branch service with repository dependency.
func NewBranchService(repo repository.BranchRepository) BranchService {
	return &branchServiceImpl{repo: repo}
}

// CreateBranch validates input and creates a branch.
func (s *branchServiceImpl) CreateBranch(ctx context.Context, branch *models.Branch) error {
	if branch == nil || branch.BranchCode == "" || branch.Name == "" || branch.Address == "" || branch.Province == "" {
		return ErrInvalidBranchInput
	}
	if !isValidBranchLocation(branch.Location) {
		return ErrInvalidBranchInput
	}

	branch.ID = primitive.NewObjectID()
	if err := s.repo.Create(ctx, branch); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return ErrBranchCodeAlreadyExists
		}
		return err
	}

	return nil
}

// GetBranchByID returns a branch by ID.
func (s *branchServiceImpl) GetBranchByID(ctx context.Context, id string) (*models.Branch, error) {
	branch, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrBranchNotFound
		}
		return nil, err
	}
	return branch, nil
}

// ListBranches returns all branches.
func (s *branchServiceImpl) ListBranches(ctx context.Context) ([]models.Branch, error) {
	return s.repo.List(ctx)
}

// UpdateBranch validates input and updates the given branch by ID.
func (s *branchServiceImpl) UpdateBranch(ctx context.Context, id string, branch *models.Branch) error {
	if branch == nil || branch.BranchCode == "" || branch.Name == "" || branch.Address == "" || branch.Province == "" {
		return ErrInvalidBranchInput
	}
	if !isValidBranchLocation(branch.Location) {
		return ErrInvalidBranchInput
	}

	if err := s.repo.UpdateByID(ctx, id, branch); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBranchNotFound
		}
		if errors.Is(err, repository.ErrDuplicate) {
			return ErrBranchCodeAlreadyExists
		}
		return err
	}

	return nil
}

// DeleteBranch removes a branch by ID.
func (s *branchServiceImpl) DeleteBranch(ctx context.Context, id string) error {
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBranchNotFound
		}
		return err
	}
	return nil
}

// NearbyBranches validates query values and returns nearby branches.
func (s *branchServiceImpl) NearbyBranches(ctx context.Context, lng float64, lat float64, maxDistance int64, limit int64) ([]models.BranchNearbyResult, error) {
	if !isValidLngLat(lng, lat) {
		return nil, ErrInvalidBranchInput
	}

	if maxDistance == 0 {
		maxDistance = 5000
	}
	if maxDistance < 1 {
		return nil, ErrInvalidBranchInput
	}

	if limit == 0 {
		limit = 10
	}
	if limit < 1 || limit > 100 {
		return nil, ErrInvalidBranchInput
	}

	return s.repo.Nearby(ctx, lng, lat, maxDistance, limit)
}

func isValidBranchLocation(location models.GeoLocation) bool {
	if location.Type != "Point" || len(location.Coordinates) != 2 {
		return false
	}

	return isValidLngLat(location.Coordinates[0], location.Coordinates[1])
}

func isValidLngLat(lng float64, lat float64) bool {
	return lng >= -180 && lng <= 180 && lat >= -90 && lat <= 90
}
