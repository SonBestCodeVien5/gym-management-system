package repository

import (
	"context"
	"errors"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BranchRepository interface {
	Create(ctx context.Context, branch *models.Branch) error
	GetByID(ctx context.Context, id string) (*models.Branch, error)
	List(ctx context.Context) ([]models.Branch, error)
	UpdateByID(ctx context.Context, id string, branch *models.Branch) error
	DeleteByID(ctx context.Context, id string) error
}

type branchRepoImpl struct {
	collection *mongo.Collection
}

// NewBranchRepository returns a repo bound to branches collection.
func NewBranchRepository(db *mongo.Database) BranchRepository {
	return &branchRepoImpl{
		collection: db.Collection("branches"),
	}
}

// Create inserts a branch document.
func (r *branchRepoImpl) Create(ctx context.Context, branch *models.Branch) error {
	_, err := r.collection.InsertOne(ctx, branch)
	return err
}

// GetByID loads a branch by ObjectID string.
func (r *branchRepoImpl) GetByID(ctx context.Context, id string) (*models.Branch, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var branch models.Branch
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&branch)
	if err != nil {
		// Normalize Mongo no-document error to ErrNotFound.
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &branch, nil
}

// List returns all branches in the collection.
func (r *branchRepoImpl) List(ctx context.Context) ([]models.Branch, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var branches []models.Branch
	if err := cursor.All(ctx, &branches); err != nil {
		return nil, err
	}

	return branches, nil
}

// UpdateByID updates mutable fields for a branch.
func (r *branchRepoImpl) UpdateByID(ctx context.Context, id string, branch *models.Branch) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{
		"branch_code": branch.BranchCode,
		"name":        branch.Name,
		"address":     branch.Address,
		"province":    branch.Province,
		"location":    branch.Location,
		"manager_id":  branch.ManagerID,
	}}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteByID removes a branch by ID.
func (r *branchRepoImpl) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
