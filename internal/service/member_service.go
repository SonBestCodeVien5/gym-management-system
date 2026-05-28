package service

import (
	"context"
	"errors"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrMemberCCIDAlreadyExists = errors.New("member with this CCID already exists")
	ErrInvalidMemberInput      = errors.New("invalid member input")
	ErrMemberNotFound          = errors.New("member not found")
)

type MemberService interface {
	RegisterMember(ctx context.Context, member *models.Member) error
	GetMemberByID(ctx context.Context, id string) (*models.Member, error)
	ActivateMember(ctx context.Context, id string) error
}

type memberServiceImpl struct {
	repo repository.MemberRepository
}

func NewMemberService(repo repository.MemberRepository) MemberService {
	return &memberServiceImpl{
		repo: repo,
	}
}

// RegisterMember validates input, checks CCID uniqueness, and creates member.
func (s *memberServiceImpl) RegisterMember(ctx context.Context, member *models.Member) error {
	// Validate required fields before touching the database.
	if member == nil || member.CCID == "" || member.FullName == "" {
		return ErrInvalidMemberInput
	}

	// Ensure CCID is unique.
	existing, err := s.repo.GetByCCID(ctx, member.CCID)
	if err == nil && existing != nil {
		return ErrMemberCCIDAlreadyExists
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		// Unexpected storage error should bubble up.
		return err
	}

	// Initialize server-side defaults.
	now := time.Now()
	member.ID = primitive.NewObjectID()
	member.IsRegistered = false
	member.IsSuspended = false
	member.TotalSessionsAttended = 0
	member.CreatedAt = now
	member.UpdatedAt = now

	// Persist member record.
	if err := s.repo.Create(ctx, member); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return ErrMemberCCIDAlreadyExists
		}
		return err
	}

	return nil
}

// GetMemberByID loads a member or returns not-found error.
func (s *memberServiceImpl) GetMemberByID(ctx context.Context, id string) (*models.Member, error) {
	// Fetch member from repository.
	member, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Map storage not-found into service-level error.
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	return member, nil
}

// ActivateMember sets is_registered=true for the member.
func (s *memberServiceImpl) ActivateMember(ctx context.Context, id string) error {
	// Flip is_registered to true in storage.
	err := s.repo.UpdateRegistrationStatus(ctx, id, true)
	if err != nil {
		// Map storage not-found into service-level error.
		if errors.Is(err, repository.ErrNotFound) {
			return ErrMemberNotFound
		}
		return err
	}

	// No additional side effects for now.
	return nil
}
