package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Attendance struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SubID    primitive.ObjectID `json:"sub_id" bson:"sub_id"`
	BranchID primitive.ObjectID `json:"branch_id" bson:"branch_id"`
	Date     time.Time          `json:"date" bson:"date"`

	// atended, absent, reported_missed, makeup
	Status string `json:"status" bson:"status"`
	// If the status is "makeup", this field will indicate which date it is making up for
	//else it will be null
	IsMakeupFor *time.Time `json:"is_makeup_for" bson:"is_makeup_for"`
}
