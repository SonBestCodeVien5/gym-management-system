package service

import (
	"context"
	"errors"
	"strings"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrInvalidSessionInput = errors.New("invalid session input")
)

type SessionService interface {
	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	ListSessions(ctx context.Context, filter repository.SessionListFilter) ([]models.Session, error)
}

type sessionServiceImpl struct {
	repo repository.SessionRepository
}

func NewSessionService(repo repository.SessionRepository) SessionService {
	return &sessionServiceImpl{repo: repo}
}

func (s *sessionServiceImpl) CreateSession(ctx context.Context, session *models.Session) error {
	if session == nil || session.BranchID.IsZero() || session.TrainerID.IsZero() || session.ScheduledAt.IsZero() {
		return ErrInvalidSessionInput
	}
	if session.DurationMin <= 0 || session.Capacity <= 0 {
		return ErrInvalidSessionInput
	}
	if !isAllowedCourseLevel(session.CourseLevel) {
		return ErrInvalidSessionInput
	}
	if len(session.Tags) == 0 {
		session.Tags = []string{}
	}

	session.ID = primitive.NewObjectID()
	session.EnrolledCount = 0
	return s.repo.Create(ctx, session)
}

func (s *sessionServiceImpl) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return session, nil
}

func (s *sessionServiceImpl) ListSessions(ctx context.Context, filter repository.SessionListFilter) ([]models.Session, error) {
	if filter.Level != "" && !isAllowedCourseLevel(filter.Level) {
		return nil, ErrInvalidSessionInput
	}
	return s.repo.List(ctx, filter)
}

func isAllowedCourseLevel(level string) bool {
	switch strings.ToLower(level) {
	case "basic", "advanced", "professional":
		return true
	default:
		return false
	}
}
