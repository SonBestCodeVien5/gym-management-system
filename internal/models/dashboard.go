package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type DashboardSummary struct {
	ActiveMembers        int64          `json:"active_members"`
	ActiveMembersDelta   int64          `json:"active_members_delta"`
	MonthlyRevenue       int64          `json:"monthly_revenue"`
	MonthlyRevenueDelta  int64          `json:"monthly_revenue_delta"`
	TodayCheckins        int64          `json:"today_checkins"`
	TodayCheckinsDelta   int64          `json:"today_checkins_delta"`
	ClassesThisWeek      int64          `json:"classes_this_week"`
	ClassesThisWeekDelta int64          `json:"classes_this_week_delta"`
	Range                DashboardRange `json:"range"`
}

type DashboardRevenueResponse struct {
	Bucket string                   `json:"bucket"`
	From   time.Time                `json:"from"`
	To     time.Time                `json:"to"`
	Items  []DashboardRevenueBucket `json:"items"`
}

type DashboardRevenueBucket struct {
	Label        string `json:"label" bson:"label"`
	GrossAmount  int64  `json:"gross_amount" bson:"gross_amount"`
	RefundAmount int64  `json:"refund_amount" bson:"refund_amount"`
	NetAmount    int64  `json:"net_amount" bson:"net_amount"`
}

type DashboardPlanDistributionResponse struct {
	Items []DashboardPlanDistributionItem `json:"items"`
}

type DashboardPlanDistributionItem struct {
	CourseID primitive.ObjectID `json:"course_id" bson:"course_id"`
	Label    string             `json:"label" bson:"label"`
	Count    int64              `json:"count" bson:"count"`
}

type DashboardRecentMembersResponse struct {
	Items []DashboardRecentMember `json:"items"`
}

type DashboardRecentMember struct {
	ID           primitive.ObjectID `json:"id"`
	FullName     string             `json:"full_name"`
	Phone        string             `json:"phone"`
	Level        string             `json:"level"`
	IsRegistered bool               `json:"is_registered"`
	CreatedAt    time.Time          `json:"created_at"`
}

type DashboardTodaySessionsResponse struct {
	Items []Session `json:"items"`
}
