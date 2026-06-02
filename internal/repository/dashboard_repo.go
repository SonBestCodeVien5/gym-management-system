package repository

import (
	"context"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DashboardRepository interface {
	CountActiveMembers(ctx context.Context) (int64, error)
	CountRegisteredMembersCreated(ctx context.Context, from time.Time, to time.Time) (int64, error)
	NetRevenueTotal(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error)
	RevenueBuckets(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) ([]models.DashboardRevenueBucket, error)
	CountCheckins(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error)
	CountSessions(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error)
	PlanDistribution(ctx context.Context, from *time.Time, to *time.Time, branchID *primitive.ObjectID) ([]models.DashboardPlanDistributionItem, error)
	RecentMembers(ctx context.Context, limit int) ([]models.DashboardRecentMember, error)
	TodaySessions(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) ([]models.Session, error)
}

type dashboardRepoImpl struct {
	members       *mongo.Collection
	subscriptions *mongo.Collection
	refunds       *mongo.Collection
	attendances   *mongo.Collection
	sessions      *mongo.Collection
}

func NewDashboardRepository(db *mongo.Database) DashboardRepository {
	return &dashboardRepoImpl{
		members:       db.Collection("members"),
		subscriptions: db.Collection("subscriptions"),
		refunds:       db.Collection("refunds"),
		attendances:   db.Collection("attendances"),
		sessions:      db.Collection("sessions"),
	}
}

func (r *dashboardRepoImpl) CountActiveMembers(ctx context.Context) (int64, error) {
	return r.members.CountDocuments(ctx, bson.M{
		"is_registered": true,
		"is_suspended":  false,
	})
}

func (r *dashboardRepoImpl) CountRegisteredMembersCreated(ctx context.Context, from time.Time, to time.Time) (int64, error) {
	return r.members.CountDocuments(ctx, bson.M{
		"is_registered": true,
		"created_at":    bson.M{"$gte": from, "$lt": to},
	})
}

func (r *dashboardRepoImpl) NetRevenueTotal(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error) {
	gross, err := r.sumSubscriptionPayments(ctx, from, to, branchID)
	if err != nil {
		return 0, err
	}
	refunds, err := r.sumRefunds(ctx, from, to, branchID)
	if err != nil {
		return 0, err
	}
	return gross - refunds, nil
}

func (r *dashboardRepoImpl) RevenueBuckets(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) ([]models.DashboardRevenueBucket, error) {
	grossBuckets, err := r.subscriptionPaymentBuckets(ctx, from, to, branchID)
	if err != nil {
		return nil, err
	}
	refundBuckets, err := r.refundBuckets(ctx, from, to, branchID)
	if err != nil {
		return nil, err
	}

	byLabel := map[string]models.DashboardRevenueBucket{}
	for _, bucket := range grossBuckets {
		byLabel[bucket.Label] = bucket
	}
	for _, bucket := range refundBuckets {
		current := byLabel[bucket.Label]
		current.Label = bucket.Label
		current.RefundAmount = bucket.RefundAmount
		current.NetAmount = current.GrossAmount - current.RefundAmount
		byLabel[bucket.Label] = current
	}

	items := make([]models.DashboardRevenueBucket, 0, len(byLabel))
	for _, bucket := range byLabel {
		bucket.NetAmount = bucket.GrossAmount - bucket.RefundAmount
		items = append(items, bucket)
	}
	return items, nil
}

func (r *dashboardRepoImpl) CountCheckins(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"date":   bson.M{"$gte": from, "$lt": to},
		"status": bson.M{"$in": []string{"attended", "makeup"}},
	}
	if branchID != nil {
		filter["branch_id"] = *branchID
	}
	return r.attendances.CountDocuments(ctx, filter)
}

func (r *dashboardRepoImpl) CountSessions(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error) {
	filter := bson.M{"scheduled_at": bson.M{"$gte": from, "$lt": to}}
	if branchID != nil {
		filter["branch_id"] = *branchID
	}
	return r.sessions.CountDocuments(ctx, filter)
}

func (r *dashboardRepoImpl) PlanDistribution(ctx context.Context, from *time.Time, to *time.Time, branchID *primitive.ObjectID) ([]models.DashboardPlanDistributionItem, error) {
	match := bson.M{
		"status": bson.M{"$in": []string{"active", "suspended", "expired"}},
	}
	if branchID != nil {
		match["home_branch_id"] = *branchID
	}
	if from != nil && to != nil {
		match["payment_date"] = bson.M{"$gte": *from, "$lt": *to}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$course_id",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "courses",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "course",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$course",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":       0,
			"course_id": "$_id",
			"label":     bson.M{"$ifNull": []any{"$course.title", ""}},
			"count":     1,
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1, "label": 1}}},
	}

	cursor, err := r.subscriptions.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []models.DashboardPlanDistributionItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	if items == nil {
		return []models.DashboardPlanDistributionItem{}, nil
	}
	return items, nil
}

func (r *dashboardRepoImpl) RecentMembers(ctx context.Context, limit int) ([]models.DashboardRecentMember, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	cursor, err := r.members.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []models.DashboardRecentMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	if members == nil {
		return []models.DashboardRecentMember{}, nil
	}
	return members, nil
}

func (r *dashboardRepoImpl) TodaySessions(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) ([]models.Session, error) {
	filter := bson.M{"scheduled_at": bson.M{"$gte": from, "$lt": to}}
	if branchID != nil {
		filter["branch_id"] = *branchID
	}

	opts := options.Find().SetSort(bson.D{{Key: "scheduled_at", Value: 1}})
	cursor, err := r.sessions.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.Session
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	if sessions == nil {
		return []models.Session{}, nil
	}
	return sessions, nil
}

func (r *dashboardRepoImpl) sumSubscriptionPayments(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error) {
	match := bson.M{"payment_date": bson.M{"$gte": from, "$lt": to}}
	if branchID != nil {
		match["home_branch_id"] = *branchID
	}

	var result struct {
		Amount int64 `bson:"amount"`
	}
	if err := aggregateOne(ctx, r.subscriptions, mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{"_id": nil, "amount": bson.M{"$sum": "$total_amount_paid"}}}},
	}, &result); err != nil {
		return 0, err
	}
	return result.Amount, nil
}

func (r *dashboardRepoImpl) sumRefunds(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error) {
	var result struct {
		Amount int64 `bson:"amount"`
	}
	if err := aggregateOne(ctx, r.refunds, refundAmountPipeline(from, to, branchID, false), &result); err != nil {
		return 0, err
	}
	return result.Amount, nil
}

func (r *dashboardRepoImpl) subscriptionPaymentBuckets(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) ([]models.DashboardRevenueBucket, error) {
	match := bson.M{"payment_date": bson.M{"$gte": from, "$lt": to}}
	if branchID != nil {
		match["home_branch_id"] = *branchID
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{"$dateToString": bson.M{
				"format":   "%Y-%m-%d",
				"date":     "$payment_date",
				"timezone": "UTC",
			}},
			"gross_amount": bson.M{"$sum": "$total_amount_paid"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":          0,
			"label":        "$_id",
			"gross_amount": 1,
		}}},
	}

	cursor, err := r.subscriptions.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var buckets []models.DashboardRevenueBucket
	if err := cursor.All(ctx, &buckets); err != nil {
		return nil, err
	}
	return buckets, nil
}

func (r *dashboardRepoImpl) refundBuckets(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) ([]models.DashboardRevenueBucket, error) {
	cursor, err := r.refunds.Aggregate(ctx, refundAmountPipeline(from, to, branchID, true))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var buckets []models.DashboardRevenueBucket
	if err := cursor.All(ctx, &buckets); err != nil {
		return nil, err
	}
	return buckets, nil
}

func refundAmountPipeline(from time.Time, to time.Time, branchID *primitive.ObjectID, bucket bool) mongo.Pipeline {
	match := bson.M{
		"status":       models.RefundStatusProcessed,
		"processed_at": bson.M{"$gte": from, "$lt": to},
	}
	pipeline := mongo.Pipeline{{{Key: "$match", Value: match}}}

	if branchID != nil {
		pipeline = append(pipeline,
			bson.D{{Key: "$lookup", Value: bson.M{
				"from":         "subscriptions",
				"localField":   "subscription_id",
				"foreignField": "_id",
				"as":           "subscription",
			}}},
			bson.D{{Key: "$unwind", Value: "$subscription"}},
			bson.D{{Key: "$match", Value: bson.M{"subscription.home_branch_id": *branchID}}},
		)
	}

	if bucket {
		pipeline = append(pipeline,
			bson.D{{Key: "$group", Value: bson.M{
				"_id": bson.M{"$dateToString": bson.M{
					"format":   "%Y-%m-%d",
					"date":     "$processed_at",
					"timezone": "UTC",
				}},
				"refund_amount": bson.M{"$sum": "$refund_amount"},
			}}},
			bson.D{{Key: "$project", Value: bson.M{
				"_id":           0,
				"label":         "$_id",
				"refund_amount": 1,
			}}},
		)
		return pipeline
	}

	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.M{"_id": nil, "amount": bson.M{"$sum": "$refund_amount"}}}})
	return pipeline
}

func aggregateOne(ctx context.Context, collection *mongo.Collection, pipeline mongo.Pipeline, out any) error {
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		return cursor.Decode(out)
	}
	return cursor.Err()
}
