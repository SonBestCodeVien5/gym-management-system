package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes creates all MongoDB indexes required by the current API surface.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	indexesByCollection := map[string][]mongo.IndexModel{
		"members": {
			{
				Keys:    bson.D{{Key: "ccid", Value: 1}},
				Options: options.Index().SetName("ccid_1").SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "created_at", Value: -1}},
				Options: options.Index().SetName("created_at_desc_idx"),
			},
		},
		"branches": {
			{
				Keys:    bson.D{{Key: "branch_code", Value: 1}},
				Options: options.Index().SetName("branch_code_unique").SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "location", Value: "2dsphere"}},
				Options: options.Index().SetName("location_2dsphere"),
			},
		},
		"subscriptions": {
			{
				Keys:    bson.D{{Key: "member_id", Value: 1}},
				Options: options.Index().SetName("member_id_idx"),
			},
			{
				Keys:    bson.D{{Key: "status", Value: 1}},
				Options: options.Index().SetName("status_idx"),
			},
			{
				Keys: bson.D{
					{Key: "member_id", Value: 1},
					{Key: "status", Value: 1},
				},
				Options: options.Index().SetName("member_status_idx"),
			},
			{
				Keys:    bson.D{{Key: "course_id", Value: 1}},
				Options: options.Index().SetName("course_id_idx"),
			},
			{
				Keys:    bson.D{{Key: "home_branch_id", Value: 1}},
				Options: options.Index().SetName("home_branch_id_idx"),
			},
			{
				Keys:    bson.D{{Key: "payment_date", Value: 1}},
				Options: options.Index().SetName("payment_date_idx"),
			},
			{
				Keys: bson.D{
					{Key: "home_branch_id", Value: 1},
					{Key: "payment_date", Value: 1},
				},
				Options: options.Index().SetName("home_branch_payment_date_idx"),
			},
		},
		"attendances": {
			{
				Keys:    bson.D{{Key: "date", Value: 1}},
				Options: options.Index().SetName("date_idx"),
			},
			{
				Keys: bson.D{
					{Key: "branch_id", Value: 1},
					{Key: "date", Value: 1},
				},
				Options: options.Index().SetName("branch_date_idx"),
			},
			{
				Keys: bson.D{
					{Key: "sub_id", Value: 1},
					{Key: "date", Value: -1},
				},
				Options: options.Index().SetName("sub_id_date_desc_idx"),
			},
			{
				Keys:    bson.D{{Key: "session_id", Value: 1}},
				Options: options.Index().SetName("session_id_idx"),
			},
			{
				Keys: bson.D{
					{Key: "session_id", Value: 1},
					{Key: "sub_id", Value: 1},
				},
				Options: options.Index().
					SetName("session_sub_unique").
					SetUnique(true).
					SetPartialFilterExpression(bson.M{"session_id": bson.M{"$exists": true}}),
			},
			{
				Keys: bson.D{
					{Key: "sub_id", Value: 1},
					{Key: "is_makeup_for", Value: 1},
					{Key: "status", Value: 1},
				},
				Options: options.Index().
					SetName("makeup_sub_ref_unique").
					SetUnique(true).
					SetPartialFilterExpression(bson.M{
						"status":        "makeup",
						"is_makeup_for": bson.M{"$exists": true},
					}),
			},
		},
		"sessions": {
			{
				Keys:    bson.D{{Key: "scheduled_at", Value: 1}},
				Options: options.Index().SetName("scheduled_at_idx"),
			},
			{
				Keys: bson.D{
					{Key: "branch_id", Value: 1},
					{Key: "scheduled_at", Value: 1},
				},
				Options: options.Index().SetName("branch_scheduled_at_idx"),
			},
			{
				Keys: bson.D{
					{Key: "course_level", Value: 1},
					{Key: "scheduled_at", Value: 1},
				},
				Options: options.Index().SetName("level_scheduled_at_idx"),
			},
			{
				Keys:    bson.D{{Key: "tags", Value: 1}},
				Options: options.Index().SetName("tags_idx"),
			},
		},
		"refunds": {
			{
				Keys:    bson.D{{Key: "subscription_id", Value: 1}},
				Options: options.Index().SetName("subscription_id_unique").SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "member_id", Value: 1}},
				Options: options.Index().SetName("member_id_idx"),
			},
			{
				Keys:    bson.D{{Key: "processed_at", Value: 1}},
				Options: options.Index().SetName("processed_at_idx"),
			},
		},
		"employees": {
			{
				Keys:    bson.D{{Key: "normalized_email", Value: 1}},
				Options: options.Index().SetName("normalized_email_unique").SetUnique(true).SetSparse(true),
			},
			{
				Keys:    bson.D{{Key: "employee_id", Value: 1}},
				Options: options.Index().SetName("employee_id_unique").SetUnique(true).SetSparse(true),
			},
			{
				Keys: bson.D{
					{Key: "role", Value: 1},
					{Key: "status", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("role_status_created_idx"),
			},
			{
				Keys: bson.D{
					{Key: "branch_id", Value: 1},
					{Key: "status", Value: 1},
				},
				Options: options.Index().SetName("branch_status_idx"),
			},
		},
		"refresh_tokens": {
			{
				Keys:    bson.D{{Key: "token_hash", Value: 1}},
				Options: options.Index().SetName("token_hash_unique").SetUnique(true),
			},
			{
				Keys: bson.D{
					{Key: "employee_id", Value: 1},
					{Key: "revoked_at", Value: 1},
				},
				Options: options.Index().SetName("employee_revoked_idx"),
			},
			{
				Keys:    bson.D{{Key: "expires_at", Value: 1}},
				Options: options.Index().SetName("expires_at_ttl").SetExpireAfterSeconds(0),
			},
		},
	}

	for collectionName, indexes := range indexesByCollection {
		if _, err := db.Collection(collectionName).Indexes().CreateMany(ctx, indexes); err != nil {
			return fmt.Errorf("create indexes for %s: %w", collectionName, err)
		}
	}

	return nil
}
