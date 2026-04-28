package repository

import (
	"context"
	"errors"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	GetByID(ctx context.Context, id string) (*models.Subscription, error)
}

type subscriptionRepoImpl struct {
	collection *mongo.Collection
}

func NewSubscriptionRepository(db *mongo.Database) SubscriptionRepository {
	return &subscriptionRepoImpl{
		collection: db.Collection("subscriptions"),
	}
}

func (r *subscriptionRepoImpl) Create(ctx context.Context, subscription *models.Subscription) error {
	_, err := r.collection.InsertOne(ctx, subscription)
	return err
}

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
