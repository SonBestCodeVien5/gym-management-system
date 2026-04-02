package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// e.g., "Yoga for Beginners", "Advanced Weightlifting"
	Title string `json:"title" bson:"title"`

	// basic, advanced, professional
	Level string `json:"level" bson:"level"`

	BasePrice    int64  `json:"base_price" bson:"base_price"` // price per session without discount
	SessionCount int    `json:"session_count" bson:"session_count"`
	Description  string `json:"description" bson:"description"`
}
