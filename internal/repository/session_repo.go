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

type SessionListFilter struct {
	BranchID string
	Level    string
	Date     *time.Time
}

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id string) (*models.Session, error)
	UpdateByID(ctx context.Context, id string, session *models.Session) error
	ReserveEnrollment(ctx context.Context, id string, subscriptionID primitive.ObjectID) (*models.Session, error)
	List(ctx context.Context, filter SessionListFilter) ([]models.Session, error)
}

type sessionRepoImpl struct {
	collection *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) SessionRepository {
	return &sessionRepoImpl{collection: db.Collection("sessions")}
}

func (r *sessionRepoImpl) Create(ctx context.Context, session *models.Session) error {
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

func (r *sessionRepoImpl) GetByID(ctx context.Context, id string) (*models.Session, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var session models.Session
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &session, nil
}

// UpdateByID updates mutable fields for a session.
func (r *sessionRepoImpl) UpdateByID(ctx context.Context, id string, session *models.Session) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"branch_id":                 session.BranchID,
		"trainer_id":                session.TrainerID,
		"course_level":              session.CourseLevel,
		"scheduled_at":              session.ScheduledAt,
		"duration_min":              session.DurationMin,
		"capacity":                  session.Capacity,
		"enrolled_count":            session.EnrolledCount,
		"enrolled_subscription_ids": session.EnrolledSubscriptionIDs,
		"tags":                      session.Tags,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// ReserveEnrollment atomically adds a subscription to a session if a slot is still available.
func (r *sessionRepoImpl) ReserveEnrollment(ctx context.Context, id string, subscriptionID primitive.ObjectID) (*models.Session, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":                       objID,
		"enrolled_subscription_ids": bson.M{"$ne": subscriptionID},
		"$expr":                     bson.M{"$lt": []any{"$enrolled_count", "$capacity"}},
	}
	update := bson.M{
		"$addToSet": bson.M{"enrolled_subscription_ids": subscriptionID},
		"$inc":      bson.M{"enrolled_count": 1},
	}

	var session models.Session
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepoImpl) List(ctx context.Context, filter SessionListFilter) ([]models.Session, error) {
	query := bson.M{}
	if filter.BranchID != "" {
		branchID, err := primitive.ObjectIDFromHex(filter.BranchID)
		if err != nil {
			return nil, err
		}
		query["branch_id"] = branchID
	}
	if filter.Level != "" {
		query["course_level"] = filter.Level
	}
	if filter.Date != nil {
		start := time.Date(filter.Date.Year(), filter.Date.Month(), filter.Date.Day(), 0, 0, 0, 0, filter.Date.Location())
		end := start.AddDate(0, 0, 1)
		query["scheduled_at"] = bson.M{"$gte": start, "$lt": end}
	}

	opts := options.Find().SetSort(bson.D{{Key: "scheduled_at", Value: 1}})
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.Session
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}
