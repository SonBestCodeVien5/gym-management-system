package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EmployeeID      string             `json:"employee_id" bson:"employee_id"` // unique employee identifier
	FullName        string             `json:"full_name" bson:"full_name"`
	Email           string             `json:"email" bson:"email"`
	NormalizedEmail string             `json:"-" bson:"normalized_email,omitempty"`
	PasswordHash    string             `json:"-" bson:"password_hash,omitempty"`
	Status          string             `json:"status" bson:"status"`
	//  "Manager", "Trainer", "Receptionist"
	Role []string `json:"role" bson:"role"`
	// basic, advanced, professional, higher levels can teach lower levels but not vice versa
	Level     string               `json:"level" bson:"level"`
	Phone     string               `json:"phone" bson:"phone"`
	BranchID  []primitive.ObjectID `json:"branch_id" bson:"branch_id"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
}
