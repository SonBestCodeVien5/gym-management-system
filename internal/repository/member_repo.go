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

// MemberRepository defines member database operations.
type MemberRepository interface {
	Create(ctx context.Context, member *models.Member) error
	GetByID(ctx context.Context, id string) (*models.Member, error)
	GetByCCID(ctx context.Context, ccid string) (*models.Member, error)
	UpdateRegistrationStatus(ctx context.Context, id string, isRegistered bool) error
	IncrementSessionsAttended(ctx context.Context, id string, delta int) error
}

// memberRepoImpl implements the MemberRepository interface.
type memberRepoImpl struct {
	collection *mongo.Collection
}

// NewMemberRepository creates the repo. Indexes are bootstrapped centrally at startup.
func NewMemberRepository(db *mongo.Database) (MemberRepository, error) {
	collection := db.Collection("members")
	return &memberRepoImpl{
		collection: collection,
	}, nil
}

// Create inserts a new member document.
func (r *memberRepoImpl) Create(ctx context.Context, member *models.Member) error {
	_, err := r.collection.InsertOne(ctx, member)
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicate
	}
	return err
}

// GetByID loads a member by ObjectID string.
func (r *memberRepoImpl) GetByID(ctx context.Context, id string) (*models.Member, error) {
	// Convert string ID to ObjectID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var member models.Member
	filter := bson.M{"_id": objID}

	err = r.collection.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &member, nil
}

// GetByCCID loads a member by CCID.
func (r *memberRepoImpl) GetByCCID(ctx context.Context, ccid string) (*models.Member, error) {
	var member models.Member
	filter := bson.M{"ccid": ccid}

	err := r.collection.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &member, nil
}

// UpdateRegistrationStatus toggles is_registered and updated_at.
func (r *memberRepoImpl) UpdateRegistrationStatus(ctx context.Context, id string, isRegistered bool) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"is_registered": isRegistered, "updated_at": time.Now()}},
	)
	if err != nil {
		return err
	}
	// No matched document means the member does not exist.
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// IncrementSessionsAttended updates total_sessions_attended by delta.
func (r *memberRepoImpl) IncrementSessionsAttended(ctx context.Context, id string, delta int) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{"total_sessions_attended": delta}, "$set": bson.M{"updated_at": time.Now()}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}
