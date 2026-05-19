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
	ErrInvalidDiscount               = errors.New("invalid discount")
	ErrSubscriptionCannotRefund      = errors.New("subscription cannot be refunded")
	ErrSubscriptionNoRemaining       = errors.New("subscription has no remaining sessions")
	ErrRefundAlreadyExists           = errors.New("refund already exists")
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) error
	GetSubscriptionByID(ctx context.Context, id string) (*models.Subscription, error)
	ListSubscriptionsByMemberID(ctx context.Context, memberID string) ([]models.Subscription, error)
	ConfirmSubscriptionPayment(ctx context.Context, memberID string, subscriptionID string) error
	SuspendSubscription(ctx context.Context, id string, suspension *models.Suspension) error
	ResumeSubscription(ctx context.Context, id string) error
	ExpireSubscription(ctx context.Context, id string) error
	RefundSubscription(ctx context.Context, id string, reason string) (*models.Refund, error)
}

type subscriptionServiceImpl struct {
	subscriptionRepo repository.SubscriptionRepository
	refundRepo       repository.RefundRepository
	memberRepo       repository.MemberRepository
	courseRepo       repository.CourseRepository
	branchRepo       repository.BranchRepository
}

func NewSubscriptionService(
	subscriptionRepo repository.SubscriptionRepository,
	refundRepo repository.RefundRepository,
	memberRepo repository.MemberRepository,
	courseRepo repository.CourseRepository,
	branchRepo repository.BranchRepository,
) SubscriptionService {
	return &subscriptionServiceImpl{
		subscriptionRepo: subscriptionRepo,
		refundRepo:       refundRepo,
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

	subtotal := course.BasePrice * int64(course.SessionCount)
	discountType := normalizeDiscountType(subscription.DiscountType)
	discountValue := subscription.DiscountValue
	discountAmount, err := calculateDiscountAmount(subtotal, discountType, discountValue)
	if err != nil {
		return err
	}

	// Snapshot course pricing into subscription at creation time.
	subscription.ID = primitive.NewObjectID()
	subscription.Status = "pending"
	subscription.AllowedTags = course.AllowedTags
	subscription.UnitPrice = course.BasePrice
	subscription.TotalSessions = course.SessionCount
	subscription.RemainingSessions = course.SessionCount
	subscription.SubtotalAmount = subtotal
	subscription.DiscountType = discountType
	subscription.DiscountValue = discountValue
	subscription.DiscountAmount = discountAmount
	subscription.TotalAmountPaid = subtotal - discountAmount

	if discountType == "none" {
		subscription.DiscountValue = 0
		subscription.DiscountAmount = 0
		subscription.PromoCode = ""
	}

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

// RefundSubscription calculates refund, atomically closes subscription, and writes audit record.
func (s *subscriptionServiceImpl) RefundSubscription(ctx context.Context, id string, reason string) (*models.Refund, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, ErrInvalidSubscriptionInput
	}

	if s.refundRepo == nil {
		return nil, errors.New("refund repository is not configured")
	}

	if _, err := s.refundRepo.GetBySubscriptionID(ctx, id); err == nil {
		return nil, ErrRefundAlreadyExists
	} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	if subscription.Status != "active" && subscription.Status != "suspended" {
		return nil, ErrSubscriptionCannotRefund
	}
	if subscription.TotalSessions <= 0 {
		return nil, ErrSubscriptionCannotRefund
	}
	if subscription.RemainingSessions <= 0 {
		return nil, ErrSubscriptionNoRemaining
	}

	usedSessions := subscription.TotalSessions - subscription.RemainingSessions
	if usedSessions < 0 {
		return nil, ErrSubscriptionCannotRefund
	}

	refundAmount := subscription.TotalAmountPaid * int64(subscription.RemainingSessions) / int64(subscription.TotalSessions)
	now := time.Now()
	refund := &models.Refund{
		ID:                primitive.NewObjectID(),
		SubscriptionID:    subscription.ID,
		MemberID:          subscription.MemberID,
		UsedSessions:      usedSessions,
		RemainingSessions: subscription.RemainingSessions,
		RefundAmount:      refundAmount,
		Reason:            strings.TrimSpace(reason),
		Status:            models.RefundStatusProcessed,
		CreatedAt:         now,
		ProcessedAt:       now,
	}

	if err := s.subscriptionRepo.RefundSubscription(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSubscriptionCannotRefund
		}
		return nil, err
	}

	if err := s.refundRepo.Create(ctx, refund); err != nil {
		return nil, err
	}

	return refund, nil
}

func normalizeDiscountType(discountType string) string {
	value := strings.TrimSpace(strings.ToLower(discountType))
	if value == "" {
		return "none"
	}
	return value
}

func calculateDiscountAmount(subtotal int64, discountType string, discountValue int64) (int64, error) {
	switch discountType {
	case "none":
		return 0, nil
	case "percent":
		if discountValue < 0 || discountValue > 100 {
			return 0, ErrInvalidDiscount
		}
		return subtotal * discountValue / 100, nil
	case "fixed":
		if discountValue < 0 || discountValue > subtotal {
			return 0, ErrInvalidDiscount
		}
		return discountValue, nil
	default:
		return 0, ErrInvalidDiscount
	}
}
