package repository

import (
	"context"
	"errors"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefundRepository interface {
	Create(ctx context.Context, refund *models.Refund) error
	GetBySubscriptionID(ctx context.Context, subscriptionID string) (*models.Refund, error)
}

type refundRepoImpl struct {
	collection *mongo.Collection
}

func NewRefundRepository(db *mongo.Database) RefundRepository {
	return &refundRepoImpl{
		collection: db.Collection("refunds"),
	}
}

func (r *refundRepoImpl) Create(ctx context.Context, refund *models.Refund) error {
	_, err := r.collection.InsertOne(ctx, refund)
	return err
}

func (r *refundRepoImpl) GetBySubscriptionID(ctx context.Context, subscriptionID string) (*models.Refund, error) {
	objID, err := primitive.ObjectIDFromHex(subscriptionID)
	if err != nil {
		return nil, err
	}

	var refund models.Refund
	err = r.collection.FindOne(ctx, bson.M{"subscription_id": objID}).Decode(&refund)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &refund, nil
}
