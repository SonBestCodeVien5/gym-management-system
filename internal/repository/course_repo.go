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
	Create(ctx context.Context, course *models.Course) error
	GetByID(ctx context.Context, id string) (*models.Course, error)
	List(ctx context.Context) ([]models.Course, error)
	UpdateByID(ctx context.Context, id string, course *models.Course) error
	DeleteByID(ctx context.Context, id string) error
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

// Create inserts a course document.
func (r *courseRepoImpl) Create(ctx context.Context, course *models.Course) error {
	_, err := r.collection.InsertOne(ctx, course)
	return err
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

// List returns all courses in the collection.
func (r *courseRepoImpl) List(ctx context.Context) ([]models.Course, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var courses []models.Course
	if err := cursor.All(ctx, &courses); err != nil {
		return nil, err
	}

	return courses, nil
}

// UpdateByID updates mutable fields for a course.
func (r *courseRepoImpl) UpdateByID(ctx context.Context, id string, course *models.Course) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{
		"title":         course.Title,
		"level":         course.Level,
		"allowed_tags":  course.AllowedTags,
		"base_price":    course.BasePrice,
		"session_count": course.SessionCount,
		"description":   course.Description,
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

// DeleteByID removes a course by ID.
func (r *courseRepoImpl) DeleteByID(ctx context.Context, id string) error {
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
