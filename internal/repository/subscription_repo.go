package repository

import (
	"context"
	"errors"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetByID(ctx context.Context, id string) (*models.Subscription, error)
	UpdateStatusAndPaymentDate(ctx context.Context, id string, status string, paymentDate time.Time) error
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateRemainingSessions(ctx context.Context, id string, remaining int) error
	UpdateRemainingSessionsAndStatus(ctx context.Context, id string, remaining int, status string) error
	RefundSubscription(ctx context.Context, id string) error
	UpdateSuspension(ctx context.Context, id string, suspension *models.Suspension, status string) error
	ClearSuspension(ctx context.Context, id string, status string) error
	ListByMemberID(ctx context.Context, memberID string) ([]models.Subscription, error)
}

type subscriptionRepoImpl struct {
	collection *mongo.Collection
}

// NewSubscriptionRepository returns a repo bound to subscriptions collection.
func NewSubscriptionRepository(db *mongo.Database) SubscriptionRepository {
	return &subscriptionRepoImpl{
		collection: db.Collection("subscriptions"),
	}
}

// Create inserts a subscription document.
func (r *subscriptionRepoImpl) Create(ctx context.Context, subscription *models.Subscription) error {
	_, err := r.collection.InsertOne(ctx, subscription)
	return err
}

// GetByID loads a subscription by ObjectID string.
func (r *subscriptionRepoImpl) GetByID(ctx context.Context, id string) (*models.Subscription, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var subscription models.Subscription
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&subscription)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &subscription, nil
}

// ListByMemberID returns all subscriptions belonging to a member.
func (r *subscriptionRepoImpl) ListByMemberID(ctx context.Context, memberID string) ([]models.Subscription, error) {
	objID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"member_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subs []models.Subscription
	if err := cursor.All(ctx, &subs); err != nil {
		return nil, err
	}

	return subs, nil
}

// UpdateStatusAndPaymentDate sets status and payment_date for a subscription.
func (r *subscriptionRepoImpl) UpdateStatusAndPaymentDate(ctx context.Context, id string, status string, paymentDate time.Time) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": status, "payment_date": paymentDate}},
	)
	if err != nil {
		return err
	}
	// No matched document means subscription does not exist.
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateStatus sets subscription status only.
func (r *subscriptionRepoImpl) UpdateStatus(ctx context.Context, id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateRemainingSessions sets remaining_sessions only.
func (r *subscriptionRepoImpl) UpdateRemainingSessions(ctx context.Context, id string, remaining int) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"remaining_sessions": remaining}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateRemainingSessionsAndStatus updates remaining_sessions and status together.
func (r *subscriptionRepoImpl) UpdateRemainingSessionsAndStatus(ctx context.Context, id string, remaining int, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"remaining_sessions": remaining, "status": status}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// RefundSubscription atomically marks an active subscription as refunded.
func (r *subscriptionRepoImpl) RefundSubscription(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":                objID,
			"status":             "active",
			"remaining_sessions": bson.M{"$gt": 0},
		},
		bson.M{"$set": bson.M{"status": "refunded", "remaining_sessions": 0}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateSuspension sets suspension details and status.
func (r *subscriptionRepoImpl) UpdateSuspension(ctx context.Context, id string, suspension *models.Suspension, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"suspension": suspension, "status": status}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// ClearSuspension removes suspension details and updates status.
func (r *subscriptionRepoImpl) ClearSuspension(ctx context.Context, id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$unset": bson.M{"suspension": ""}, "$set": bson.M{"status": status}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}
