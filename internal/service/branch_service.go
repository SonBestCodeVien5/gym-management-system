package service

import (
	"context"
	"errors"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrBranchNotFound     = errors.New("branch not found")
	ErrInvalidBranchInput = errors.New("invalid branch input")
)

// BranchService defines business operations for branch management.
type BranchService interface {
	CreateBranch(ctx context.Context, branch *models.Branch) error
	GetBranchByID(ctx context.Context, id string) (*models.Branch, error)
	ListBranches(ctx context.Context) ([]models.Branch, error)
	UpdateBranch(ctx context.Context, id string, branch *models.Branch) error
	DeleteBranch(ctx context.Context, id string) error
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
	if branch.Location.Type == "" || len(branch.Location.Coordinates) != 2 {
		return ErrInvalidBranchInput
	}

	branch.ID = primitive.NewObjectID()
	return s.repo.Create(ctx, branch)
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
	if branch.Location.Type == "" || len(branch.Location.Coordinates) != 2 {
		return ErrInvalidBranchInput
	}

	if err := s.repo.UpdateByID(ctx, id, branch); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBranchNotFound
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
