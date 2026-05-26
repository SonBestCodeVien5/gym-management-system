package repository

import (
	"context"
	"errors"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	FindActiveByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeActiveByHash(ctx context.Context, tokenHash string, revokedAt time.Time) error
	RevokeActiveByEmployeeID(ctx context.Context, employeeID primitive.ObjectID, revokedAt time.Time) error
}

type refreshTokenRepoImpl struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) (RefreshTokenRepository, error) {
	collection := db.Collection("refresh_tokens")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "token_hash", Value: 1}},
		Options: options.Index().SetName("token_hash_unique").SetUnique(true),
	}
	if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return nil, err
	}

	return &refreshTokenRepoImpl{collection: collection}, nil
}

func (r *refreshTokenRepoImpl) Create(ctx context.Context, token *models.RefreshToken) error {
	if token.ID.IsZero() {
		token.ID = primitive.NewObjectID()
	}
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *refreshTokenRepoImpl) FindActiveByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.collection.FindOne(ctx, bson.M{
		"token_hash": tokenHash,
		"revoked_at": bson.M{"$exists": false},
	}).Decode(&token)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepoImpl) RevokeActiveByHash(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"token_hash": tokenHash, "revoked_at": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"revoked_at": revokedAt, "updated_at": revokedAt}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *refreshTokenRepoImpl) RevokeActiveByEmployeeID(ctx context.Context, employeeID primitive.ObjectID, revokedAt time.Time) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"employee_id": employeeID, "revoked_at": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"revoked_at": revokedAt, "updated_at": revokedAt}},
	)
	return err
}
