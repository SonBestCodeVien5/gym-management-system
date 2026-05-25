package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshToken struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EmployeeID primitive.ObjectID `json:"employee_id" bson:"employee_id"`
	TokenHash  string             `json:"-" bson:"token_hash"`
	ExpiresAt  time.Time          `json:"expires_at" bson:"expires_at"`
	RevokedAt  *time.Time         `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}
