package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Member struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CCID     string             `bson:"ccid" json:"ccid"`
	FullName string             `bson:"full_name" json:"full_name"`
	Email    string             `bson:"email" json:"email"`
	Phone    string             `bson:"phone" json:"phone"`
	Gender   string             `bson:"gender" json:"gender"`
	Level    string             `bson:"level" json:"level"`

	IsRegistered          bool      `bson:"is_registered" json:"is_registered"`
	// Using bool IsSuspended to indicate if the member is suspended or not. If true, the member is suspended and cannot attend sessions
	// and cannot suspend in the future. If false, the member is not suspended and can attend sessions and can suspend in the future.
	IsSuspended           bool      `bson:"is_suspended" json:"is_suspended"`
	TotalSessionsAttended int       `bson:"total_sessions_attended" json:"total_sessions_attended"`
	CreatedAt             time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time `bson:"updated_at" json:"updated_at"`
}
