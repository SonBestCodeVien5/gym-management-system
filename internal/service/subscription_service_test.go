package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type stubSubscriptionRepo struct {
	subscription *models.Subscription
	createErr    error
	refundErr    error
	created      *models.Subscription
	refundedID   string
}

func (r *stubSubscriptionRepo) Create(ctx context.Context, subscription *models.Subscription) error {
	r.created = subscription
	return r.createErr
}

func (r *stubSubscriptionRepo) GetByID(ctx context.Context, id string) (*models.Subscription, error) {
	if r.subscription == nil {
		return nil, repository.ErrNotFound
	}
	return r.subscription, nil
}

func (r *stubSubscriptionRepo) UpdateStatusAndPaymentDate(ctx context.Context, id string, status string, paymentDate time.Time) error {
	return nil
}

func (r *stubSubscriptionRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	return nil
}

func (r *stubSubscriptionRepo) UpdateRemainingSessions(ctx context.Context, id string, remaining int) error {
	return nil
}

func (r *stubSubscriptionRepo) UpdateRemainingSessionsAndStatus(ctx context.Context, id string, remaining int, status string) error {
	return nil
}

func (r *stubSubscriptionRepo) RefundSubscription(ctx context.Context, id string) error {
	r.refundedID = id
	if r.refundErr != nil {
		return r.refundErr
	}
	return nil
}

func (r *stubSubscriptionRepo) UpdateSuspension(ctx context.Context, id string, suspension *models.Suspension, status string) error {
	return nil
}

func (r *stubSubscriptionRepo) ClearSuspension(ctx context.Context, id string, status string) error {
	return nil
}

func (r *stubSubscriptionRepo) ListByMemberID(ctx context.Context, memberID string) ([]models.Subscription, error) {
	return nil, nil
}

type stubRefundRepo struct {
	existing  *models.Refund
	createErr error
	created   *models.Refund
}

func (r *stubRefundRepo) Create(ctx context.Context, refund *models.Refund) error {
	r.created = refund
	return r.createErr
}

func (r *stubRefundRepo) GetBySubscriptionID(ctx context.Context, subscriptionID string) (*models.Refund, error) {
	if r.existing == nil {
		return nil, repository.ErrNotFound
	}
	return r.existing, nil
}

type stubMemberRepo struct {
	member *models.Member
	err    error
}

func (r *stubMemberRepo) Create(ctx context.Context, member *models.Member) error {
	return nil
}

func (r *stubMemberRepo) GetByID(ctx context.Context, id string) (*models.Member, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.member == nil {
		return nil, repository.ErrNotFound
	}
	return r.member, nil
}

func (r *stubMemberRepo) GetByCCID(ctx context.Context, ccid string) (*models.Member, error) {
	return nil, repository.ErrNotFound
}

func (r *stubMemberRepo) UpdateRegistrationStatus(ctx context.Context, id string, isRegistered bool) error {
	return nil
}

func (r *stubMemberRepo) IncrementSessionsAttended(ctx context.Context, id string, delta int) error {
	return nil
}

type stubCourseRepo struct {
	course *models.Course
	err    error
}

func (r *stubCourseRepo) Create(ctx context.Context, course *models.Course) error {
	return nil
}

func (r *stubCourseRepo) GetByID(ctx context.Context, id string) (*models.Course, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.course == nil {
		return nil, repository.ErrNotFound
	}
	return r.course, nil
}

func (r *stubCourseRepo) List(ctx context.Context) ([]models.Course, error) {
	return nil, nil
}

func (r *stubCourseRepo) UpdateByID(ctx context.Context, id string, course *models.Course) error {
	return nil
}

func (r *stubCourseRepo) DeleteByID(ctx context.Context, id string) error {
	return nil
}

type stubBranchRepo struct {
	branch *models.Branch
	err    error
}

func (r *stubBranchRepo) Create(ctx context.Context, branch *models.Branch) error {
	return nil
}

func (r *stubBranchRepo) GetByID(ctx context.Context, id string) (*models.Branch, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.branch == nil {
		return nil, repository.ErrNotFound
	}
	return r.branch, nil
}

func (r *stubBranchRepo) List(ctx context.Context) ([]models.Branch, error) {
	return nil, nil
}

func (r *stubBranchRepo) UpdateByID(ctx context.Context, id string, branch *models.Branch) error {
	return nil
}

func (r *stubBranchRepo) DeleteByID(ctx context.Context, id string) error {
	return nil
}

func (r *stubBranchRepo) Nearby(ctx context.Context, lng float64, lat float64, maxDistance int64, limit int64) ([]models.BranchNearbyResult, error) {
	return nil, nil
}

func TestCreateSubscriptionPricingRules(t *testing.T) {
	memberID := primitive.NewObjectID()
	courseID := primitive.NewObjectID()
	branchID := primitive.NewObjectID()

	tests := []struct {
		name               string
		discountType       string
		discountValue      int64
		wantDiscountType   string
		wantDiscountValue  int64
		wantDiscountAmount int64
		wantTotalPaid      int64
		wantErr            error
	}{
		{
			name:               "no discount",
			discountType:       "",
			discountValue:      50,
			wantDiscountType:   "none",
			wantDiscountValue:  0,
			wantDiscountAmount: 0,
			wantTotalPaid:      1_200_000,
		},
		{
			name:               "percent discount",
			discountType:       "percent",
			discountValue:      25,
			wantDiscountType:   "percent",
			wantDiscountValue:  25,
			wantDiscountAmount: 300_000,
			wantTotalPaid:      900_000,
		},
		{
			name:               "fixed discount",
			discountType:       "fixed",
			discountValue:      200_000,
			wantDiscountType:   "fixed",
			wantDiscountValue:  200_000,
			wantDiscountAmount: 200_000,
			wantTotalPaid:      1_000_000,
		},
		{
			name:          "invalid percent discount",
			discountType:  "percent",
			discountValue: 101,
			wantErr:       ErrInvalidDiscount,
		},
		{
			name:          "invalid fixed discount",
			discountType:  "fixed",
			discountValue: 1_200_001,
			wantErr:       ErrInvalidDiscount,
		},
		{
			name:          "invalid discount type",
			discountType:  "voucher",
			discountValue: 10,
			wantErr:       ErrInvalidDiscount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subRepo := &stubSubscriptionRepo{}
			svc := NewSubscriptionService(
				subRepo,
				&stubRefundRepo{},
				&stubMemberRepo{member: &models.Member{ID: memberID}},
				&stubCourseRepo{course: &models.Course{
					ID:           courseID,
					BasePrice:    100_000,
					SessionCount: 12,
					AllowedTags:  []string{"basic", "yoga"},
				}},
				&stubBranchRepo{branch: &models.Branch{ID: branchID}},
			)

			subscription := &models.Subscription{
				MemberID:       memberID,
				CourseID:       courseID,
				HomeBranchID:   branchID,
				StartDate:      time.Now(),
				EndDate:        time.Now().AddDate(0, 1, 0),
				SessionPerWeek: 3,
				DiscountType:   tt.discountType,
				DiscountValue:  tt.discountValue,
				PromoCode:      "PROMO",
			}

			err := svc.CreateSubscription(context.Background(), subscription)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("CreateSubscription() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("CreateSubscription() unexpected error = %v", err)
			}
			if subRepo.created == nil {
				t.Fatal("CreateSubscription() did not persist subscription")
			}

			got := subRepo.created
			if got.Status != "pending" {
				t.Fatalf("Status = %q, want pending", got.Status)
			}
			if got.UnitPrice != 100_000 {
				t.Fatalf("UnitPrice = %d, want 100000", got.UnitPrice)
			}
			if got.TotalSessions != 12 || got.RemainingSessions != 12 {
				t.Fatalf("sessions = total %d remaining %d, want 12/12", got.TotalSessions, got.RemainingSessions)
			}
			if got.SubtotalAmount != 1_200_000 {
				t.Fatalf("SubtotalAmount = %d, want 1200000", got.SubtotalAmount)
			}
			if got.DiscountType != tt.wantDiscountType {
				t.Fatalf("DiscountType = %q, want %q", got.DiscountType, tt.wantDiscountType)
			}
			if got.DiscountValue != tt.wantDiscountValue {
				t.Fatalf("DiscountValue = %d, want %d", got.DiscountValue, tt.wantDiscountValue)
			}
			if got.DiscountAmount != tt.wantDiscountAmount {
				t.Fatalf("DiscountAmount = %d, want %d", got.DiscountAmount, tt.wantDiscountAmount)
			}
			if got.TotalAmountPaid != tt.wantTotalPaid {
				t.Fatalf("TotalAmountPaid = %d, want %d", got.TotalAmountPaid, tt.wantTotalPaid)
			}
		})
	}
}

func TestRefundSubscription(t *testing.T) {
	subscriptionID := primitive.NewObjectID()
	memberID := primitive.NewObjectID()

	tests := []struct {
		name         string
		id           string
		subscription *models.Subscription
		existing     *models.Refund
		refundErr    error
		wantAmount   int64
		wantErr      error
	}{
		{
			name: "pending cannot refund",
			id:   subscriptionID.Hex(),
			subscription: &models.Subscription{
				ID:                subscriptionID,
				MemberID:          memberID,
				Status:            "pending",
				TotalSessions:     12,
				RemainingSessions: 8,
				TotalAmountPaid:   900_000,
			},
			wantErr: ErrSubscriptionCannotRefund,
		},
		{
			name: "expired cannot refund",
			id:   subscriptionID.Hex(),
			subscription: &models.Subscription{
				ID:                subscriptionID,
				MemberID:          memberID,
				Status:            "expired",
				TotalSessions:     12,
				RemainingSessions: 8,
				TotalAmountPaid:   900_000,
			},
			wantErr: ErrSubscriptionCannotRefund,
		},
		{
			name: "refunded already by audit record",
			id:   subscriptionID.Hex(),
			subscription: &models.Subscription{
				ID:                subscriptionID,
				MemberID:          memberID,
				Status:            "active",
				TotalSessions:     12,
				RemainingSessions: 8,
				TotalAmountPaid:   900_000,
			},
			existing: &models.Refund{SubscriptionID: subscriptionID},
			wantErr:  ErrRefundAlreadyExists,
		},
		{
			name: "no remaining sessions",
			id:   subscriptionID.Hex(),
			subscription: &models.Subscription{
				ID:                subscriptionID,
				MemberID:          memberID,
				Status:            "active",
				TotalSessions:     12,
				RemainingSessions: 0,
				TotalAmountPaid:   900_000,
			},
			wantErr: ErrSubscriptionNoRemaining,
		},
		{
			name: "remaining greater than total is data conflict",
			id:   subscriptionID.Hex(),
			subscription: &models.Subscription{
				ID:                subscriptionID,
				MemberID:          memberID,
				Status:            "active",
				TotalSessions:     12,
				RemainingSessions: 13,
				TotalAmountPaid:   900_000,
			},
			wantErr: ErrSubscriptionCannotRefund,
		},
		{
			name: "atomic update conflict maps to cannot refund",
			id:   subscriptionID.Hex(),
			subscription: &models.Subscription{
				ID:                subscriptionID,
				MemberID:          memberID,
				Status:            "active",
				TotalSessions:     12,
				RemainingSessions: 8,
				TotalAmountPaid:   900_000,
			},
			refundErr: repository.ErrNotFound,
			wantErr:   ErrSubscriptionCannotRefund,
		},
		{
			name: "active refund success",
			id:   subscriptionID.Hex(),
			subscription: &models.Subscription{
				ID:                subscriptionID,
				MemberID:          memberID,
				Status:            "active",
				TotalSessions:     12,
				RemainingSessions: 8,
				TotalAmountPaid:   900_000,
			},
			wantAmount: 600_000,
		},
		{
			name: "suspended cannot refund",
			id:   subscriptionID.Hex(),
			subscription: &models.Subscription{
				ID:                subscriptionID,
				MemberID:          memberID,
				Status:            "suspended",
				TotalSessions:     10,
				RemainingSessions: 3,
				TotalAmountPaid:   1_000_000,
			},
			wantErr: ErrSubscriptionCannotRefund,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subRepo := &stubSubscriptionRepo{subscription: tt.subscription, refundErr: tt.refundErr}
			refundRepo := &stubRefundRepo{existing: tt.existing}
			svc := NewSubscriptionService(
				subRepo,
				refundRepo,
				&stubMemberRepo{},
				&stubCourseRepo{},
				&stubBranchRepo{},
			)

			refund, err := svc.RefundSubscription(context.Background(), tt.id, "  member requested cancellation  ")
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("RefundSubscription() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("RefundSubscription() unexpected error = %v", err)
			}
			if refund == nil {
				t.Fatal("RefundSubscription() refund is nil")
			}
			if refund.RefundAmount != tt.wantAmount {
				t.Fatalf("RefundAmount = %d, want %d", refund.RefundAmount, tt.wantAmount)
			}
			if refund.UsedSessions != tt.subscription.TotalSessions-tt.subscription.RemainingSessions {
				t.Fatalf("UsedSessions = %d, want %d", refund.UsedSessions, tt.subscription.TotalSessions-tt.subscription.RemainingSessions)
			}
			if refund.Reason != "member requested cancellation" {
				t.Fatalf("Reason = %q, want trimmed reason", refund.Reason)
			}
			if refund.Status != models.RefundStatusProcessed {
				t.Fatalf("Status = %q, want %q", refund.Status, models.RefundStatusProcessed)
			}
			if refundRepo.created == nil {
				t.Fatal("RefundSubscription() did not create refund audit")
			}
			if subRepo.refundedID != tt.id {
				t.Fatalf("refundedID = %q, want %q", subRepo.refundedID, tt.id)
			}
		})
	}
}
