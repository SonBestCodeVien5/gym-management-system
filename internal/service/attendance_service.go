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
	ErrWeeklySessionLimitReached   = errors.New("weekly session limit reached")
	ErrReportedMissedLimitReached  = errors.New("reported missed limit reached within 30 days")
	ErrMakeupReferenceInvalid      = errors.New("invalid makeup reference")
	ErrMakeupReferenceNotFound     = errors.New("makeup reference not found")
	ErrMakeupAlreadyUsed           = errors.New("makeup reference already used")
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

	// 3) Enforce weekly quota for attended/makeup records.
	if attendance.Status == "attended" || attendance.Status == "makeup" {
		weeklyCount, err := s.countWeeklySessions(ctx, attendance.SubID.Hex(), attendance.Date)
		if err != nil {
			return err
		}
		if weeklyCount >= subscription.SessionPerWeek {
			return ErrWeeklySessionLimitReached
		}
	}

	// 4) Only these statuses are accepted for this initial version.
	validStatus := attendance.Status == "attended" || attendance.Status == "makeup" || attendance.Status == "absent" || attendance.Status == "reported_missed"
	if !validStatus {
		return ErrInvalidAttendanceInput
	}

	// 4.1) Enforce report/makeup-specific rules.
	if attendance.Status == "reported_missed" {
		if err := s.validateReportedMissedWindow(ctx, attendance.SubID.Hex(), attendance.Date); err != nil {
			return err
		}
	}
	if attendance.Status == "makeup" {
		if err := s.validateMakeupRequest(ctx, attendance.SubID.Hex(), attendance.Date, attendance.IsMakeupFor); err != nil {
			return err
		}
	}

	// 5) Create attendance record first.
	attendance.ID = primitive.NewObjectID()
	if err := s.attendanceRepo.Create(ctx, attendance); err != nil {
		return err
	}

	// 6) For attended/makeup, decrease remaining_sessions and increase member attended count.
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

// validateReportedMissedWindow enforces one reported_missed record per 30-day sliding window.
func (s *attendanceServiceImpl) validateReportedMissedWindow(ctx context.Context, subscriptionID string, at time.Time) error {
	records, err := s.attendanceRepo.ListBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	windowStart := at.AddDate(0, 0, -30)
	for _, record := range records {
		if record.Status != "reported_missed" {
			continue
		}
		if !record.Date.Before(windowStart) && !record.Date.After(at) {
			return ErrReportedMissedLimitReached
		}
	}

	return nil
}

// validateMakeupRequest checks that makeup references a recent reported_missed record and has not been used.
func (s *attendanceServiceImpl) validateMakeupRequest(ctx context.Context, subscriptionID string, makeupAt time.Time, makeupFor *time.Time) error {
	if makeupFor == nil || makeupFor.IsZero() {
		return ErrMakeupReferenceInvalid
	}
	if makeupFor.After(makeupAt) {
		return ErrMakeupReferenceInvalid
	}
	if makeupAt.Sub(*makeupFor) > 7*24*time.Hour {
		return ErrMakeupReferenceInvalid
	}

	records, err := s.attendanceRepo.ListBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	var sourceReport *models.Attendance
	for _, record := range records {
		if record.Status == "reported_missed" && record.Date.Equal(*makeupFor) {
			recordCopy := record
			sourceReport = &recordCopy
			break
		}
	}
	if sourceReport == nil {
		return ErrMakeupReferenceNotFound
	}

	for _, record := range records {
		if record.Status != "makeup" || record.IsMakeupFor == nil {
			continue
		}
		if record.IsMakeupFor.Equal(*makeupFor) {
			return ErrMakeupAlreadyUsed
		}
	}

	return nil
}

// countWeeklySessions counts attended and makeup records within the current Monday-Sunday window.
func (s *attendanceServiceImpl) countWeeklySessions(ctx context.Context, subscriptionID string, at time.Time) (int, error) {
	records, err := s.attendanceRepo.ListBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return 0, err
	}

	startOfWeek := beginningOfISOWeek(at)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	count := 0
	for _, record := range records {
		if record.Date.Before(startOfWeek) || !record.Date.Before(endOfWeek) {
			continue
		}
		if record.Status == "attended" || record.Status == "makeup" {
			count++
		}
	}

	return count, nil
}

// beginningOfISOWeek returns Monday 00:00:00 in the same location as the input time.
func beginningOfISOWeek(at time.Time) time.Time {
	loc := at.Location()
	dayStart := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, loc)
	offset := (int(dayStart.Weekday()) + 6) % 7
	return dayStart.AddDate(0, 0, -offset)
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
