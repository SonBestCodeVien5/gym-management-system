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
	ErrInvalidSubscriptionInput      = errors.New("invalid subscription input")
	ErrSubscriptionNotFound          = errors.New("subscription not found")
	ErrSubscriptionReferenceNotFound = errors.New("subscription reference not found")
	ErrSubscriptionAlreadyActive     = errors.New("subscription already active")
	ErrInvalidSubscriptionStatus     = errors.New("invalid subscription status")
	ErrSubscriptionMemberMismatch    = errors.New("subscription does not belong to member")
	ErrSubscriptionAlreadySuspended  = errors.New("subscription already suspended")
	ErrSubscriptionNotActive         = errors.New("subscription is not active")
	ErrSubscriptionExpired           = errors.New("subscription is expired")
	ErrInvalidSuspensionPeriod       = errors.New("invalid suspension period")
	ErrSubscriptionMemberNotFound    = errors.New("member not found")
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) error
	GetSubscriptionByID(ctx context.Context, id string) (*models.Subscription, error)
	ListSubscriptionsByMemberID(ctx context.Context, memberID string) ([]models.Subscription, error)
	ConfirmSubscriptionPayment(ctx context.Context, memberID string, subscriptionID string) error
	SuspendSubscription(ctx context.Context, id string, suspension *models.Suspension) error
	ResumeSubscription(ctx context.Context, id string) error
	ExpireSubscription(ctx context.Context, id string) error
}

type subscriptionServiceImpl struct {
	subscriptionRepo repository.SubscriptionRepository
	memberRepo       repository.MemberRepository
	courseRepo       repository.CourseRepository
	branchRepo       repository.BranchRepository
}

func NewSubscriptionService(
	subscriptionRepo repository.SubscriptionRepository,
	memberRepo repository.MemberRepository,
	courseRepo repository.CourseRepository,
	branchRepo repository.BranchRepository,
) SubscriptionService {
	return &subscriptionServiceImpl{
		subscriptionRepo: subscriptionRepo,
		memberRepo:       memberRepo,
		courseRepo:       courseRepo,
		branchRepo:       branchRepo,
	}
}

// CreateSubscription validates input and snapshots course pricing into subscription.
func (s *subscriptionServiceImpl) CreateSubscription(ctx context.Context, subscription *models.Subscription) error {
	// Basic input validation.
	if subscription == nil || subscription.MemberID.IsZero() || subscription.CourseID.IsZero() || subscription.HomeBranchID.IsZero() {
		return ErrInvalidSubscriptionInput
	}
	if subscription.SessionPerWeek <= 0 {
		return ErrInvalidSubscriptionInput
	}
	if subscription.StartDate.IsZero() || subscription.EndDate.IsZero() || subscription.StartDate.After(subscription.EndDate) {
		return ErrInvalidSubscriptionInput
	}

	// Validate reference IDs exist in database.
	member, err := s.memberRepo.GetByID(ctx, subscription.MemberID.Hex())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionReferenceNotFound
		}
		// Bubble up unexpected storage errors.
		return err
	}
	if member == nil {
		return ErrSubscriptionReferenceNotFound
	}

	course, err := s.courseRepo.GetByID(ctx, subscription.CourseID.Hex())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionReferenceNotFound
		}
		// Bubble up unexpected storage errors.
		return err
	}
	if course == nil {
		return ErrSubscriptionReferenceNotFound
	}

	branch, err := s.branchRepo.GetByID(ctx, subscription.HomeBranchID.Hex())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionReferenceNotFound
		}
		// Bubble up unexpected storage errors.
		return err
	}
	if branch == nil {
		return ErrSubscriptionReferenceNotFound
	}

	// Snapshot course pricing into subscription at creation time.
	subscription.ID = primitive.NewObjectID()
	subscription.Status = "pending"
	subscription.AllowedTags = course.AllowedTags
	subscription.UnitPrice = course.BasePrice
	subscription.TotalSessions = course.SessionCount
	subscription.Total_Amount_Paid = course.BasePrice * int64(course.SessionCount)
	subscription.RemainingSessions = course.SessionCount

	// Persist subscription record.
	return s.subscriptionRepo.Create(ctx, subscription)
}

// GetSubscriptionByID returns subscription by ID or not-found error.
func (s *subscriptionServiceImpl) GetSubscriptionByID(ctx context.Context, id string) (*models.Subscription, error) {
	// Fetch subscription by ID.
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		// Map storage not-found into service-level error.
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return subscription, nil
}

// ListSubscriptionsByMemberID returns all subscriptions that belong to a member.
func (s *subscriptionServiceImpl) ListSubscriptionsByMemberID(ctx context.Context, memberID string) ([]models.Subscription, error) {
	if _, err := primitive.ObjectIDFromHex(memberID); err != nil {
		return nil, ErrInvalidSubscriptionInput
	}

	if _, err := s.memberRepo.GetByID(ctx, memberID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSubscriptionMemberNotFound
		}
		return nil, err
	}

	subscriptions, err := s.subscriptionRepo.ListByMemberID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	if subscriptions == nil {
		return []models.Subscription{}, nil
	}
	return subscriptions, nil
}

// ConfirmSubscriptionPayment activates a pending subscription tied to the member.
func (s *subscriptionServiceImpl) ConfirmSubscriptionPayment(ctx context.Context, memberID string, subscriptionID string) error {
	// Load subscription to verify ownership and status.
	subscription, err := s.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		// Bubble up unexpected storage errors.
		return err
	}

	// Prevent confirming a subscription that belongs to another member.
	if subscription.MemberID.Hex() != memberID {
		return ErrSubscriptionMemberMismatch
	}

	// Only pending subscriptions can be activated.
	if subscription.Status == "active" {
		return ErrSubscriptionAlreadyActive
	}
	if subscription.Status != "pending" {
		return ErrInvalidSubscriptionStatus
	}

	// Update status and payment date in the database.
	now := time.Now()
	if err := s.subscriptionRepo.UpdateStatusAndPaymentDate(ctx, subscriptionID, "active", now); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		// Bubble up unexpected storage errors.
		return err
	}

	// No additional side effects for now.
	return nil
}

// SuspendSubscription sets suspension details and marks subscription as suspended.
func (s *subscriptionServiceImpl) SuspendSubscription(ctx context.Context, id string, suspension *models.Suspension) error {
	// Validate suspension payload.
	if suspension == nil || suspension.StartDate.IsZero() || suspension.EndDate.IsZero() {
		return ErrInvalidSuspensionPeriod
	}
	if suspension.StartDate.After(suspension.EndDate) {
		return ErrInvalidSuspensionPeriod
	}
	if suspension.FrozenSession < 0 {
		return ErrInvalidSuspensionPeriod
	}

	// Load subscription to validate current status.
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}

	if subscription.Status == "suspended" {
		return ErrSubscriptionAlreadySuspended
	}
	if subscription.Status != "active" {
		return ErrSubscriptionNotActive
	}
	if subscription.EndDate.Before(time.Now()) {
		return ErrSubscriptionExpired
	}

	// Persist suspension details and status.
	if err := s.subscriptionRepo.UpdateSuspension(ctx, id, suspension, "suspended"); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}

	return nil
}

// ResumeSubscription clears suspension and sets status back to active.
func (s *subscriptionServiceImpl) ResumeSubscription(ctx context.Context, id string) error {
	// Load subscription to validate current status.
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}

	if subscription.Status != "suspended" {
		return ErrInvalidSubscriptionStatus
	}
	if subscription.EndDate.Before(time.Now()) {
		return ErrSubscriptionExpired
	}

	if err := s.subscriptionRepo.ClearSuspension(ctx, id, "active"); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}

	return nil
}

// ExpireSubscription marks subscription as expired regardless of remaining sessions.
func (s *subscriptionServiceImpl) ExpireSubscription(ctx context.Context, id string) error {
	if err := s.subscriptionRepo.UpdateStatus(ctx, id, "expired"); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}

	return nil
}
