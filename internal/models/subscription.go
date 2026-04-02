package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Suspension struct {
	StartDate     time.Time `bson:"start_date" json:"start_date"`
	EndDate       time.Time `bson:"end_date" json:"end_date"`
	FrozenSession int       `bson:"frozen_session" json:"frozen_session"`
	Reason        string    `bson:"reason" json:"reason"`
}

type Subscription struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	MemberID     primitive.ObjectID `bson:"member_id" json:"member_id"`
	CourseID     primitive.ObjectID `bson:"course_id" json:"course_id"`
	HomeBranchID primitive.ObjectID `bson:"home_branch_id" json:"home_branch_id"`
	Status       string             `bson:"status" json:"status"` // active, suspended, expired, refunded

	// fiancial details
	PaymentDate       time.Time `bson:"payment_date" json:"payment_date"`
	Total_Amount_Paid int64     `bson:"total_amount_paid" json:"total_amount_paid"`
	UnitPrice         int64     `bson:"unit_price" json:"unit_price"` // price per session
	TotalSessions     int       `bson:"total_sessions" json:"total_sessions"`
	RemainingSessions int       `bson:"remaining_sessions" json:"remaining_sessions"`

	// subscription details
	StartDate      time.Time   `bson:"start_date" json:"start_date"`
	EndDate        time.Time   `bson:"end_date" json:"end_date"`
	SessionPerWeek int         `bson:"session_per_week" json:"session_per_week"`
	Suspension     *Suspension `bson:"suspension,omitempty" json:"suspension,omitempty"`
}
