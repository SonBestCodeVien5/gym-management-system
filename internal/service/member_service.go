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
)

type MemberService interface {
	RegisterMember(ctx context.Context, member *models.Member) error
	GetMemberByID(ctx context.Context, id string) (*models.Member, error)
}

type memberServiceImpl struct {
	repo repository.MemberRepository
}

func NewMemberService(repo repository.MemberRepository) MemberService {
	return &memberServiceImpl{
		repo: repo,
	}
}

func (s *memberServiceImpl) RegisterMember(ctx context.Context, member *models.Member) error {
	if member == nil || member.CCID == "" || member.FullName == "" {
		return ErrInvalidMemberInput
	}

	existing, err := s.repo.GetByCCID(ctx, member.CCID)
	if err == nil && existing != nil {
		return ErrMemberCCIDAlreadyExists
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	now := time.Now()
	member.ID = primitive.NewObjectID()
	member.IsRegistered = false
	member.IsSuspended = false
	member.TotalSessionsAttended = 0
	member.CreatedAt = now
	member.UpdatedAt = now

	return s.repo.Create(ctx, member)
}

func (s *memberServiceImpl) GetMemberByID(ctx context.Context, id string) (*models.Member, error) {
	return s.repo.GetByID(ctx, id)
}
