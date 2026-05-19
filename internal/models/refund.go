package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RefundStatusProcessed = "processed"
)

type Refund struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SubscriptionID    primitive.ObjectID `bson:"subscription_id" json:"subscription_id"`
	MemberID          primitive.ObjectID `bson:"member_id" json:"member_id"`
	UsedSessions      int                `bson:"used_sessions" json:"used_sessions"`
	RemainingSessions int                `bson:"remaining_sessions" json:"remaining_sessions"`
	RefundAmount      int64              `bson:"refund_amount" json:"refund_amount"`
	Reason            string             `bson:"reason" json:"reason"`
	Status            string             `bson:"status" json:"status"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	ProcessedAt       time.Time          `bson:"processed_at" json:"processed_at"`
}
