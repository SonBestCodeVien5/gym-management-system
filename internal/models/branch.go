package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GeoLocation struct {
	Type        string    `json:"type" bson:"type"` // e.g., "Point"
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type Branch struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BranchCode string             `json:"branch_code" bson:"branch_code"`
	Name       string             `json:"name" bson:"name"`
	Address    string             `json:"address" bson:"address"`
	Province   string             `json:"province" bson:"province"`
	Location   GeoLocation        `json:"location" bson:"location"`
	ManagerID  primitive.ObjectID `json:"manager_id" bson:"manager_id"`
}
