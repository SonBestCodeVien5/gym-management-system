package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/testutil"
	"go.mongodb.org/mongo-driver/bson"
)

func TestIntegrationAuthRoleGuardAndSmoke(t *testing.T) {
	app := testutil.NewTestApp(t)

	ping := app.DoJSON(t, http.MethodGet, "/ping", "", nil)
	testutil.AssertStatus(t, ping, http.StatusOK)

	noToken := app.DoJSON(t, http.MethodGet, "/api/v1/employees", "", nil)
	testutil.AssertError(t, noToken, http.StatusUnauthorized, "UNAUTHORIZED")

	adminList := app.DoJSON(t, http.MethodGet, "/api/v1/employees", app.AdminToken, nil)
	testutil.AssertStatus(t, adminList, http.StatusOK)

	staffEmail := testutil.Unique("receptionist") + "@gym.test"
	staffPassword := "staff-password-123"
	createStaff := app.DoJSON(t, http.MethodPost, "/api/v1/employees", app.AdminToken, map[string]any{
		"employee_id": testutil.Unique("EMP"),
		"full_name":   "Integration Receptionist",
		"email":       staffEmail,
		"password":    staffPassword,
		"role":        []string{"receptionist"},
	})
	testutil.AssertStatus(t, createStaff, http.StatusCreated)

	staffToken, _ := app.Login(t, staffEmail, staffPassword)
	forbidden := app.DoJSON(t, http.MethodGet, "/api/v1/employees", staffToken, nil)
	testutil.AssertError(t, forbidden, http.StatusForbidden, "FORBIDDEN")

	refresh := app.DoJSON(t, http.MethodPost, "/api/v1/auth/refresh", "", map[string]any{
		"refresh_token": app.AdminRefresh,
	})
	testutil.AssertStatus(t, refresh, http.StatusOK)
	refreshed := testutil.DataMap(t, refresh)
	newRefreshToken := testutil.DataString(t, refreshed, "refresh_token")

	logout := app.DoJSON(t, http.MethodPost, "/api/v1/auth/logout", "", map[string]any{
		"refresh_token": newRefreshToken,
	})
	testutil.AssertStatus(t, logout, http.StatusOK)

	reuseLoggedOutToken := app.DoJSON(t, http.MethodPost, "/api/v1/auth/refresh", "", map[string]any{
		"refresh_token": newRefreshToken,
	})
	testutil.AssertError(t, reuseLoggedOutToken, http.StatusUnauthorized, "UNAUTHORIZED")
}

func TestIntegrationMemberSubscriptionAndDataIntegrityConflicts(t *testing.T) {
	app := testutil.NewTestApp(t)

	branchCode := testutil.Unique("BR")
	branchID := app.CreateBranch(t, app.AdminToken, branchCode, 106.7009, 10.7769)
	duplicateBranch := app.DoJSON(t, http.MethodPost, "/api/v1/branches", app.AdminToken, map[string]any{
		"branch_code": branchCode,
		"name":        "Duplicate Branch",
		"address":     "456 Test Street",
		"province":    "Ho Chi Minh",
		"location": map[string]any{
			"type":        "Point",
			"coordinates": []float64{106.7010, 10.7770},
		},
	})
	testutil.AssertError(t, duplicateBranch, http.StatusConflict, "CONFLICT")

	courseID := app.CreateCourse(t, app.AdminToken)
	ccid := testutil.Unique("CCID")
	memberID := app.CreateMember(t, app.AdminToken, ccid)
	duplicateMember := app.DoJSON(t, http.MethodPost, "/api/v1/members", app.AdminToken, map[string]any{
		"ccid":      ccid,
		"full_name": "Duplicate Member",
		"email":     "duplicate-" + ccid + "@gym.test",
		"phone":     "0900000001",
		"gender":    "other",
		"level":     "basic",
	})
	testutil.AssertError(t, duplicateMember, http.StatusConflict, "CONFLICT")

	subscriptionID := app.CreateSubscription(t, app.AdminToken, memberID, courseID, branchID)
	app.ActivateSubscription(t, app.AdminToken, memberID, subscriptionID)

	getSubscription := app.DoJSON(t, http.MethodGet, "/api/v1/subscriptions/"+subscriptionID, app.AdminToken, nil)
	testutil.AssertStatus(t, getSubscription, http.StatusOK)

	memberSubscriptions := app.DoJSON(t, http.MethodGet, "/api/v1/members/"+memberID+"/subscriptions", app.AdminToken, nil)
	testutil.AssertStatus(t, memberSubscriptions, http.StatusOK)
	if got := len(testutil.DataSlice(t, memberSubscriptions)); got != 1 {
		t.Fatalf("member subscriptions count = %d, want 1", got)
	}

	refund := app.DoJSON(t, http.MethodPost, "/api/v1/subscriptions/"+subscriptionID+"/refund", app.AdminToken, map[string]any{
		"reason": "integration refund",
	})
	testutil.AssertStatus(t, refund, http.StatusOK)

	duplicateRefund := app.DoJSON(t, http.MethodPost, "/api/v1/subscriptions/"+subscriptionID+"/refund", app.AdminToken, map[string]any{
		"reason": "duplicate refund",
	})
	testutil.AssertError(t, duplicateRefund, http.StatusConflict, "CONFLICT")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := app.DB.Collection("refunds").CountDocuments(ctx, bson.M{})
	if err != nil {
		t.Fatalf("count refunds: %v", err)
	}
	if count != 1 {
		t.Fatalf("refund document count = %d, want 1", count)
	}
}

func TestIntegrationAttendanceMakeupReuseConflict(t *testing.T) {
	app := testutil.NewTestApp(t)
	fixture := app.CreateActiveCoreFixture(t)

	reportedMissedAt := "2026-05-12T08:00:00Z"
	makeupAt := "2026-05-14T08:00:00Z"
	report := app.DoJSON(t, http.MethodPost, "/api/v1/attendance/report", app.AdminToken, map[string]any{
		"subscription_id": fixture.SubscriptionID,
		"branch_id":       fixture.BranchID,
		"date":            reportedMissedAt,
	})
	testutil.AssertStatus(t, report, http.StatusCreated)

	makeup := app.DoJSON(t, http.MethodPost, "/api/v1/attendance/makeup", app.AdminToken, map[string]any{
		"subscription_id": fixture.SubscriptionID,
		"branch_id":       fixture.BranchID,
		"date":            makeupAt,
		"is_makeup_for":   reportedMissedAt,
	})
	testutil.AssertStatus(t, makeup, http.StatusCreated)

	reusedMakeup := app.DoJSON(t, http.MethodPost, "/api/v1/attendance/makeup", app.AdminToken, map[string]any{
		"subscription_id": fixture.SubscriptionID,
		"branch_id":       fixture.BranchID,
		"date":            "2026-05-15T08:00:00Z",
		"is_makeup_for":   reportedMissedAt,
	})
	testutil.AssertError(t, reusedMakeup, http.StatusConflict, "CONFLICT")
}

func TestIntegrationBranchNearby(t *testing.T) {
	app := testutil.NewTestApp(t)

	app.CreateBranch(t, app.AdminToken, testutil.Unique("HCM"), 106.7009, 10.7769)
	app.CreateBranch(t, app.AdminToken, testutil.Unique("HCM"), 106.7020, 10.7780)

	nearby := app.DoJSON(t, http.MethodGet, "/api/v1/branches/nearby?lng=106.7009&lat=10.7769&max_distance=5000&limit=10", app.AdminToken, nil)
	testutil.AssertStatus(t, nearby, http.StatusOK)
	results := testutil.DataSlice(t, nearby)
	if len(results) == 0 {
		t.Fatalf("nearby branches count = 0, want non-empty")
	}
	if _, ok := results[0]["distance_meters"]; !ok {
		t.Fatalf("nearby branch does not include distance_meters: %#v", results[0])
	}

	invalidQuery := app.DoJSON(t, http.MethodGet, "/api/v1/branches/nearby?lng=200&lat=10.7769", app.AdminToken, nil)
	testutil.AssertError(t, invalidQuery, http.StatusBadRequest, "INVALID_INPUT")
}
