package testutil

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
)

type CoreFixture struct {
	BranchID       string
	CourseID       string
	MemberID       string
	SubscriptionID string
}

func Unique(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, nextUnique())
}

var uniqueCounter uint64

func nextUnique() uint64 {
	return atomic.AddUint64(&uniqueCounter, 1)
}

func (a *TestApp) CreateBranch(t *testing.T, token string, code string, lng float64, lat float64) string {
	t.Helper()

	res := a.DoJSON(t, http.MethodPost, "/api/v1/branches", token, map[string]any{
		"branch_code": code,
		"name":        "Integration Branch " + code,
		"address":     "123 Test Street",
		"province":    "Ho Chi Minh",
		"location": map[string]any{
			"type":        "Point",
			"coordinates": []float64{lng, lat},
		},
	})
	AssertStatus(t, res, http.StatusCreated)
	return DataString(t, DataMap(t, res), "id")
}

func (a *TestApp) CreateCourse(t *testing.T, token string) string {
	t.Helper()

	res := a.DoJSON(t, http.MethodPost, "/api/v1/courses", token, map[string]any{
		"title":         "Integration Basic Course",
		"level":         "basic",
		"allowed_tags":  []string{"basic", "yoga"},
		"base_price":    100000,
		"session_count": 12,
		"description":   "Integration test course",
	})
	AssertStatus(t, res, http.StatusCreated)
	return DataString(t, DataMap(t, res), "id")
}

func (a *TestApp) CreateMember(t *testing.T, token string, ccid string) string {
	t.Helper()

	res := a.DoJSON(t, http.MethodPost, "/api/v1/members", token, map[string]any{
		"ccid":      ccid,
		"full_name": "Integration Member",
		"email":     ccid + "@gym.test",
		"phone":     "0900000000",
		"gender":    "other",
		"level":     "basic",
	})
	AssertStatus(t, res, http.StatusCreated)
	return DataString(t, DataMap(t, res), "id")
}

func (a *TestApp) CreateSubscription(t *testing.T, token string, memberID string, courseID string, branchID string) string {
	t.Helper()

	res := a.DoJSON(t, http.MethodPost, "/api/v1/subscriptions", token, map[string]any{
		"member_id":        memberID,
		"course_id":        courseID,
		"home_branch_id":   branchID,
		"start_date":       "2026-05-01T00:00:00Z",
		"end_date":         "2026-12-31T00:00:00Z",
		"session_per_week": 3,
		"discount_type":    "",
		"discount_value":   0,
	})
	AssertStatus(t, res, http.StatusCreated)
	return DataString(t, DataMap(t, res), "id")
}

func (a *TestApp) ActivateSubscription(t *testing.T, token string, memberID string, subscriptionID string) {
	t.Helper()

	res := a.DoJSON(t, http.MethodPatch, "/api/v1/members/"+memberID+"/activate", token, map[string]any{
		"subscription_id": subscriptionID,
	})
	AssertStatus(t, res, http.StatusOK)
}

func (a *TestApp) CreateActiveCoreFixture(t *testing.T) CoreFixture {
	t.Helper()

	branchID := a.CreateBranch(t, a.AdminToken, "BR"+fmt.Sprint(nextUnique()), 106.7009, 10.7769)
	courseID := a.CreateCourse(t, a.AdminToken)
	memberID := a.CreateMember(t, a.AdminToken, "CCID"+fmt.Sprint(nextUnique()))
	subscriptionID := a.CreateSubscription(t, a.AdminToken, memberID, courseID, branchID)
	a.ActivateSubscription(t, a.AdminToken, memberID, subscriptionID)
	return CoreFixture{
		BranchID:       branchID,
		CourseID:       courseID,
		MemberID:       memberID,
		SubscriptionID: subscriptionID,
	}
}
