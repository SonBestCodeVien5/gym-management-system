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
	ErrInvalidAttendanceInput      = errors.New("invalid attendance input")
	ErrAttendanceCheckInNotAllowed = errors.New("attendance check-in is not allowed for current subscription status")
	ErrNoRemainingSessions         = errors.New("no remaining sessions")
)

// AttendanceService defines check-in and attendance history operations.
type AttendanceService interface {
	CheckIn(ctx context.Context, attendance *models.Attendance) error
	ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]models.Attendance, error)
}

type attendanceServiceImpl struct {
	attendanceRepo   repository.AttendanceRepository
	subscriptionRepo repository.SubscriptionRepository
	memberRepo       repository.MemberRepository
}

// NewAttendanceService builds attendance service with required repositories.
func NewAttendanceService(attendanceRepo repository.AttendanceRepository, subscriptionRepo repository.SubscriptionRepository, memberRepo repository.MemberRepository) AttendanceService {
	return &attendanceServiceImpl{
		attendanceRepo:   attendanceRepo,
		subscriptionRepo: subscriptionRepo,
		memberRepo:       memberRepo,
	}
}

// CheckIn validates attendance and updates related counters.
func (s *attendanceServiceImpl) CheckIn(ctx context.Context, attendance *models.Attendance) error {
	// 1) Validate required fields.
	if attendance == nil || attendance.SubID.IsZero() || attendance.BranchID.IsZero() {
		return ErrInvalidAttendanceInput
	}
	if attendance.Status == "" {
		return ErrInvalidAttendanceInput
	}
	if attendance.Date.IsZero() {
		attendance.Date = time.Now()
	}

	// 2) Load subscription and validate business state.
	subscription, err := s.subscriptionRepo.GetByID(ctx, attendance.SubID.Hex())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}
	if subscription.Status != "active" {
		return ErrAttendanceCheckInNotAllowed
	}
	if subscription.EndDate.Before(attendance.Date) {
		return ErrSubscriptionExpired
	}

	// 3) Only these statuses are accepted for this initial version.
	validStatus := attendance.Status == "attended" || attendance.Status == "makeup" || attendance.Status == "absent" || attendance.Status == "reported_missed"
	if !validStatus {
		return ErrInvalidAttendanceInput
	}

	// 4) Create attendance record first.
	attendance.ID = primitive.NewObjectID()
	if err := s.attendanceRepo.Create(ctx, attendance); err != nil {
		return err
	}

	// 5) For attended/makeup, decrease remaining_sessions and increase member attended count.
	if attendance.Status == "attended" || attendance.Status == "makeup" {
		if subscription.RemainingSessions <= 0 {
			return ErrNoRemainingSessions
		}

		newRemaining := subscription.RemainingSessions - 1
		if newRemaining <= 0 {
			if err := s.subscriptionRepo.UpdateRemainingSessionsAndStatus(ctx, attendance.SubID.Hex(), 0, "expired"); err != nil {
				return err
			}
		} else {
			if err := s.subscriptionRepo.UpdateRemainingSessions(ctx, attendance.SubID.Hex(), newRemaining); err != nil {
				return err
			}
		}

		if err := s.memberRepo.IncrementSessionsAttended(ctx, subscription.MemberID.Hex(), 1); err != nil {
			return err
		}
	}

	return nil
}

// ListBySubscriptionID returns attendance history for a subscription.
func (s *attendanceServiceImpl) ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]models.Attendance, error) {
	if _, err := primitive.ObjectIDFromHex(subscriptionID); err != nil {
		return nil, ErrInvalidAttendanceInput
	}

	records, err := s.attendanceRepo.ListBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	return records, nil
}
