package repository

import(
	"context"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// MemberRepository defines the interface for member-related database operations
type MemberRepository interface {
	Create(ctx context.Context, member *models.Member) error
	GetByID(ctx context.Context, id string) (*models.Member, error)
	GetByCCID(ctx context.Context, ccid string) (*models.Member, error)
}

//memberRepoImpl implements the MemberRepository interface
type memberRepoImpl struct {
	collection *mongo.Collection
}

// NewMemberRepository creates a new instance of MemberRepository
func NewMemberRepository(db *mongo.Database, collectionName string) MemberRepository {
	return &memberRepoImpl{
		collection: db.Collection("members"),
	}
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
		return nil, err
	}
	return &member, nil
}