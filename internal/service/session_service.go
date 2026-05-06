package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrSessionNotFound        = errors.New("session not found")
	ErrInvalidSessionInput    = errors.New("invalid session input")
	ErrSessionAlreadyFull     = errors.New("session is full")
	ErrSessionAlreadyEnrolled = errors.New("subscription already enrolled in session")
	ErrSessionNotEnrolled     = errors.New("subscription is not enrolled in session")
	ErrSessionCheckInClosed   = errors.New("session check-in is closed")
)

type SessionService interface {
	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	ListSessions(ctx context.Context, filter repository.SessionListFilter) ([]models.Session, error)
	EnrollSubscription(ctx context.Context, sessionID string, subscriptionID string) (*models.Session, error)
	CheckInSubscription(ctx context.Context, sessionID string, subscriptionID string) (*models.Attendance, error)
}

type sessionServiceImpl struct {
	repo              repository.SessionRepository
	subscriptionRepo  repository.SubscriptionRepository
	attendanceRepo    repository.AttendanceRepository
	attendanceService AttendanceService
}

func NewSessionService(
	repo repository.SessionRepository,
	subscriptionRepo repository.SubscriptionRepository,
	attendanceRepo repository.AttendanceRepository,
	attendanceService AttendanceService,
) SessionService {
	return &sessionServiceImpl{
		repo:              repo,
		subscriptionRepo:  subscriptionRepo,
		attendanceRepo:    attendanceRepo,
		attendanceService: attendanceService,
	}
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

// EnrollSubscription reserves a slot in a session for a subscription.
func (s *sessionServiceImpl) EnrollSubscription(ctx context.Context, sessionID string, subscriptionID string) (*models.Session, error) {
	if _, err := primitive.ObjectIDFromHex(subscriptionID); err != nil {
		return nil, ErrInvalidSessionInput
	}

	session, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidSessionInput
		}
		return nil, err
	}
	if subscription.Status != "active" || subscription.EndDate.Before(time.Now()) {
		return nil, ErrInvalidSessionInput
	}

	for _, enrolledID := range session.EnrolledSubscriptionIDs {
		if enrolledID == subscription.ID {
			return nil, ErrSessionAlreadyEnrolled
		}
	}
	if session.EnrolledCount >= session.Capacity {
		return nil, ErrSessionAlreadyFull
	}

	session.EnrolledSubscriptionIDs = append(session.EnrolledSubscriptionIDs, subscription.ID)
	session.EnrolledCount++
	if err := s.repo.UpdateByID(ctx, sessionID, session); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	return session, nil
}

// CheckInSubscription records an attendance against an enrolled session.
func (s *sessionServiceImpl) CheckInSubscription(ctx context.Context, sessionID string, subscriptionID string) (*models.Attendance, error) {
	if _, err := primitive.ObjectIDFromHex(subscriptionID); err != nil {
		return nil, ErrInvalidSessionInput
	}

	session, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidSessionInput
		}
		return nil, err
	}
	if subscription.Status != "active" || subscription.EndDate.Before(time.Now()) {
		return nil, ErrInvalidSessionInput
	}

	enrolled := false
	for _, enrolledID := range session.EnrolledSubscriptionIDs {
		if enrolledID == subscription.ID {
			enrolled = true
			break
		}
	}
	if !enrolled {
		return nil, ErrSessionNotEnrolled
	}

	records, err := s.attendanceRepo.ListBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		if record.SessionID != nil && *record.SessionID == session.ID {
			return nil, ErrSessionCheckInClosed
		}
	}

	attendance := &models.Attendance{
		SubID:     subscription.ID,
		BranchID:  session.BranchID,
		SessionID: &session.ID,
		Date:      time.Now(),
		Status:    "attended",
	}
	if err := s.attendanceService.CheckIn(ctx, attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func isAllowedCourseLevel(level string) bool {
	switch strings.ToLower(level) {
	case "basic", "advanced", "professional":
		return true
	default:
		return false
	}
}
