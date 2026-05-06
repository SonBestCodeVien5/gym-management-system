package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID                      primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	BranchID                primitive.ObjectID   `json:"branch_id" bson:"branch_id"`
	TrainerID               primitive.ObjectID   `json:"trainer_id" bson:"trainer_id"`
	CourseLevel             string               `json:"course_level" bson:"course_level"`
	ScheduledAt             time.Time            `json:"scheduled_at" bson:"scheduled_at"`
	DurationMin             int                  `json:"duration_min" bson:"duration_min"`
	Capacity                int                  `json:"capacity" bson:"capacity"`
	EnrolledCount           int                  `json:"enrolled_count" bson:"enrolled_count"`
	EnrolledSubscriptionIDs []primitive.ObjectID `json:"enrolled_subscription_ids,omitempty" bson:"enrolled_subscription_ids,omitempty"`
	Tags                    []string             `json:"tags" bson:"tags"`
}
