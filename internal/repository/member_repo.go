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

// MemberRepository defines the interface for member-related database operations
type MemberRepository interface {
	Create(ctx context.Context, member *models.Member) error
	GetByID(ctx context.Context, id string) (*models.Member, error)
	GetByCCID(ctx context.Context, ccid string) (*models.Member, error)
	UpdateRegistrationStatus(ctx context.Context, id string, isRegistered bool) error
}

// memberRepoImpl implements the MemberRepository interface
type memberRepoImpl struct {
	collection *mongo.Collection
}

// NewMemberRepository creates a new instance of MemberRepository
func NewMemberRepository(db *mongo.Database) (MemberRepository, error) {
	collection := db.Collection("members")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "ccid", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return nil, err
	}

	return &memberRepoImpl{
		collection: collection,
	}, nil
}

// --Start using database operations--
// 1. Create a new member
func (r *memberRepoImpl) Create(ctx context.Context, member *models.Member) error {
	_, err := r.collection.InsertOne(ctx, member)
	return err
}

// 2. Get a member by ID
func (r *memberRepoImpl) GetByID(ctx context.Context, id string) (*models.Member, error) {
	// Convert string ID to ObjectID
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

// 3. Get a member by CCID
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

// 4. Update member registration status
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
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}
