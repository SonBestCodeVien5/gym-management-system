package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	EmployeeID string             `json:"employee_id" bson:"employee_id"` // unique employee identifier
	FullName   string             `json:"full_name" bson:"full_name"`
	//  "Manager", "Trainer", "Receptionist"
	Role []string `json:"role" bson:"role"`
	// basic, advanced, professional, higher levels can teach lower levels but not vice versa
	Level    string               `json:"level" bson:"level"`
	Phone    string               `json:"phone" bson:"phone"`
	Email    string               `json:"email" bson:"email"`
	BranchID []primitive.ObjectID `json:"branch_id" bson:"branch_id"`
}
