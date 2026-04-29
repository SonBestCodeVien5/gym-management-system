package repository

import (
	"context"
	"errors"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseRepository interface {
	GetByID(ctx context.Context, id string) (*models.Course, error)
}

type courseRepoImpl struct {
	collection *mongo.Collection
}

// NewCourseRepository returns a repo bound to courses collection.
func NewCourseRepository(db *mongo.Database) CourseRepository {
	return &courseRepoImpl{
		collection: db.Collection("courses"),
	}
}

// GetByID loads a course by ObjectID string.
func (r *courseRepoImpl) GetByID(ctx context.Context, id string) (*models.Course, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var course models.Course
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&course)
	if err != nil {
		// Normalize Mongo no-document error to ErrNotFound.
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &course, nil
}
