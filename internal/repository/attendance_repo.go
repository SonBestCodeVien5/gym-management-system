package repository

import (
	"context"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AttendanceRepository defines storage operations for attendance records.
type AttendanceRepository interface {
	Create(ctx context.Context, attendance *models.Attendance) error
	ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]models.Attendance, error)
}

type attendanceRepoImpl struct {
	collection *mongo.Collection
}

// NewAttendanceRepository returns a repo bound to attendances collection.
func NewAttendanceRepository(db *mongo.Database) AttendanceRepository {
	return &attendanceRepoImpl{collection: db.Collection("attendances")}
}

// Create inserts a new attendance record.
func (r *attendanceRepoImpl) Create(ctx context.Context, attendance *models.Attendance) error {
	_, err := r.collection.InsertOne(ctx, attendance)
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicate
	}
	return err
}

// ListBySubscriptionID returns attendance history sorted by date descending.
func (r *attendanceRepoImpl) ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]models.Attendance, error) {
	objID, err := primitive.ObjectIDFromHex(subscriptionID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"sub_id": objID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []models.Attendance
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	return records, nil
}
