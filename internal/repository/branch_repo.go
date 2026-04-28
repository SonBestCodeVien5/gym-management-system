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
	GetByID(ctx context.Context, id string) (*models.Branch, error)
}

type branchRepoImpl struct {
	collection *mongo.Collection
}

func NewBranchRepository(db *mongo.Database) BranchRepository {
	return &branchRepoImpl{
		collection: db.Collection("branches"),
	}
}

func (r *branchRepoImpl) GetByID(ctx context.Context, id string) (*models.Branch, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var branch models.Branch
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&branch)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &branch, nil
}
