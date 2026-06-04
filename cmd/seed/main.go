package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/SonBestCodeVien5/gym-management-system/pkg/database"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const demoPassword = "demo123456"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables instead")
	}

	mongoURI := strings.TrimSpace(os.Getenv("MONGODB_URI"))
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is required")
	}

	client, err := database.ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbName := databaseNameFromEnv()
	db := client.Database(dbName)
	if err := database.EnsureIndexes(ctx, db); err != nil {
		log.Fatalf("Failed to ensure indexes: %v", err)
	}

	result, err := seedDemoData(ctx, db)
	if err != nil {
		log.Fatalf("Failed to seed demo data: %v", err)
	}

	fmt.Printf("Seeded demo data into database %q.\n", dbName)
	for label, count := range result {
		fmt.Printf("- %s: %d\n", label, count)
	}
	fmt.Println()
	fmt.Println("Demo login accounts:")
	fmt.Printf("- admin: admin@gym.test / %s\n", demoPassword)
	fmt.Printf("- manager: manager@gym.test / %s\n", demoPassword)
	fmt.Printf("- receptionist: receptionist@gym.test / %s\n", demoPassword)
	fmt.Printf("- trainer: trainer@gym.test / %s\n", demoPassword)
}

func seedDemoData(ctx context.Context, db *mongo.Database) (map[string]int, error) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	ids := demoIDs{
		admin:        oid("665000000000000000000001"),
		manager:      oid("665000000000000000000002"),
		receptionist: oid("665000000000000000000003"),
		trainer:      oid("665000000000000000000004"),

		branchCentral: oid("665000000000000000000101"),
		branchWest:    oid("665000000000000000000102"),
		branchEast:    oid("665000000000000000000103"),

		courseBasic:        oid("665000000000000000000201"),
		courseIntermediate: oid("665000000000000000000202"),
		courseAdvanced:     oid("665000000000000000000203"),

		memberActive:       oid("665000000000000000000301"),
		memberSuspended:    oid("665000000000000000000302"),
		memberPending:      oid("665000000000000000000303"),
		memberExpired:      oid("665000000000000000000304"),
		memberRefunded:     oid("665000000000000000000305"),
		memberSessionReady: oid("665000000000000000000306"),

		subActive:       oid("665000000000000000000401"),
		subSuspended:    oid("665000000000000000000402"),
		subPending:      oid("665000000000000000000403"),
		subExpired:      oid("665000000000000000000404"),
		subRefunded:     oid("665000000000000000000405"),
		subSessionReady: oid("665000000000000000000406"),

		attendanceToday:   oid("665000000000000000000501"),
		attendanceReport:  oid("665000000000000000000502"),
		attendanceMakeup:  oid("665000000000000000000503"),
		attendanceSession: oid("665000000000000000000504"),

		sessionMorning: oid("665000000000000000000601"),
		sessionEvening: oid("665000000000000000000602"),
		sessionFuture:  oid("665000000000000000000603"),

		refund: oid("665000000000000000000701"),
	}

	employees, err := demoEmployees(ids, now, false)
	if err != nil {
		return nil, err
	}

	counts := map[string]int{}
	employeeCollection := db.Collection("employees")
	if err := upsertManyWithFilter(ctx, employeeCollection, employees, func(doc bson.M) bson.M {
		return bson.M{"normalized_email": doc["normalized_email"]}
	}); err != nil {
		return nil, fmt.Errorf("employees: %w", err)
	}
	counts["employees"] = len(employees)
	if err := captureEmployeeIDs(ctx, employeeCollection, &ids); err != nil {
		return nil, err
	}

	branches := []models.Branch{
		{
			ID:         ids.branchCentral,
			BranchCode: "IFG-CEN",
			Name:       "Iron Forge Central",
			Address:    "72 Nguyen Hue, Ben Nghe Ward",
			Province:   "Ho Chi Minh City",
			Location:   models.GeoLocation{Type: "Point", Coordinates: []float64{106.7020, 10.7758}},
			ManagerID:  ids.manager,
		},
		{
			ID:         ids.branchWest,
			BranchCode: "IFG-WST",
			Name:       "Iron Forge West",
			Address:    "18 Quang Trung, Go Vap District",
			Province:   "Ho Chi Minh City",
			Location:   models.GeoLocation{Type: "Point", Coordinates: []float64{106.6695, 10.8231}},
			ManagerID:  ids.manager,
		},
		{
			ID:         ids.branchEast,
			BranchCode: "IFG-EST",
			Name:       "Iron Forge East",
			Address:    "25 Mai Chi Tho, Thu Duc City",
			Province:   "Ho Chi Minh City",
			Location:   models.GeoLocation{Type: "Point", Coordinates: []float64{106.7592, 10.7726}},
			ManagerID:  ids.manager,
		},
	}
	branchCollection := db.Collection("branches")
	if err := upsertManyWithFilter(ctx, branchCollection, branches, func(doc bson.M) bson.M {
		return bson.M{"branch_code": doc["branch_code"]}
	}); err != nil {
		return nil, fmt.Errorf("branches: %w", err)
	}
	counts["branches"] = len(branches)
	if err := captureBranchIDs(ctx, branchCollection, &ids); err != nil {
		return nil, err
	}

	employees, err = demoEmployees(ids, now, true)
	if err != nil {
		return nil, err
	}
	if err := upsertManyWithFilter(ctx, employeeCollection, employees, func(doc bson.M) bson.M {
		return bson.M{"normalized_email": doc["normalized_email"]}
	}); err != nil {
		return nil, fmt.Errorf("employees branch assignment: %w", err)
	}

	courses := []models.Course{
		{
			ID:           ids.courseBasic,
			Title:        "Foundation Strength",
			Level:        "basic",
			AllowedTags:  []string{"strength", "mobility", "beginner"},
			BasePrice:    180000,
			SessionCount: 12,
			Description:  "Beginner-friendly strength and movement fundamentals.",
		},
		{
			ID:           ids.courseIntermediate,
			Title:        "Hybrid Conditioning",
			Level:        "advanced",
			AllowedTags:  []string{"conditioning", "strength", "cardio"},
			BasePrice:    220000,
			SessionCount: 16,
			Description:  "Mixed strength and conditioning program for regular members.",
		},
		{
			ID:           ids.courseAdvanced,
			Title:        "Performance Barbell",
			Level:        "professional",
			AllowedTags:  []string{"barbell", "power", "advanced"},
			BasePrice:    280000,
			SessionCount: 20,
			Description:  "Advanced coached barbell sessions with performance tracking.",
		},
	}
	courseCollection := db.Collection("courses")
	if err := upsertManyWithFilter(ctx, courseCollection, courses, func(doc bson.M) bson.M {
		return bson.M{"title": doc["title"]}
	}); err != nil {
		return nil, fmt.Errorf("courses: %w", err)
	}
	counts["courses"] = len(courses)
	if err := captureCourseIDs(ctx, courseCollection, &ids); err != nil {
		return nil, err
	}
	courses[0].ID = ids.courseBasic
	courses[1].ID = ids.courseIntermediate
	courses[2].ID = ids.courseAdvanced

	members := []models.Member{
		member(ids.memberActive, "0911000001", "CCID-DEMO-001", "Nguyen Minh Khoa", "khoa.demo@gym.test", "male", "advanced", true, 5, now.AddDate(0, 0, -20), now),
		member(ids.memberSuspended, "0911000002", "CCID-DEMO-002", "Tran Ha Linh", "linh.demo@gym.test", "female", "basic", true, 2, now.AddDate(0, 0, -18), now),
		member(ids.memberPending, "0911000003", "CCID-DEMO-003", "Pham Gia Bao", "bao.demo@gym.test", "male", "basic", false, 0, now.AddDate(0, 0, -2), now),
		member(ids.memberExpired, "0911000004", "CCID-DEMO-004", "Le Anh Thu", "thu.demo@gym.test", "female", "advanced", true, 12, now.AddDate(0, -2, 0), now),
		member(ids.memberRefunded, "0911000005", "CCID-DEMO-005", "Vo Quoc Huy", "huy.demo@gym.test", "male", "professional", true, 3, now.AddDate(0, 0, -12), now),
		member(ids.memberSessionReady, "0911000006", "CCID-DEMO-006", "Dang Minh Anh", "anh.demo@gym.test", "female", "advanced", true, 1, now.AddDate(0, 0, -6), now),
	}
	memberCollection := db.Collection("members")
	if err := upsertManyWithFilter(ctx, memberCollection, members, func(doc bson.M) bson.M {
		return bson.M{"ccid": doc["ccid"]}
	}); err != nil {
		return nil, fmt.Errorf("members: %w", err)
	}
	counts["members"] = len(members)
	if err := captureMemberIDs(ctx, memberCollection, &ids); err != nil {
		return nil, err
	}

	subscriptions := []models.Subscription{
		subscription(ids.subActive, ids.memberActive, ids.courseIntermediate, ids.branchCentral, courses[1], "active", ptrTime(now.AddDate(0, 0, -15)), nil, now.AddDate(0, 0, -15), now.AddDate(0, 1, 15), 3, 11),
		subscription(ids.subSuspended, ids.memberSuspended, ids.courseBasic, ids.branchWest, courses[0], "suspended", ptrTime(now.AddDate(0, 0, -14)), &models.Suspension{
			StartDate:     today.AddDate(0, 0, -1),
			EndDate:       today.AddDate(0, 0, 6),
			FrozenSession: 2,
			Reason:        "Demo travel hold",
		}, now.AddDate(0, 0, -14), now.AddDate(0, 1, 14), 2, 9),
		subscription(ids.subPending, ids.memberPending, ids.courseBasic, ids.branchCentral, courses[0], "pending", nil, nil, now.AddDate(0, 0, -2), now.AddDate(0, 1, -2), 2, courses[0].SessionCount),
		subscription(ids.subExpired, ids.memberExpired, ids.courseIntermediate, ids.branchEast, courses[1], "expired", ptrTime(now.AddDate(0, -2, 3)), nil, now.AddDate(0, -2, 0), now.AddDate(0, 0, -3), 3, 0),
		subscription(ids.subRefunded, ids.memberRefunded, ids.courseAdvanced, ids.branchEast, courses[2], "refunded", ptrTime(now.AddDate(0, 0, -10)), nil, now.AddDate(0, 0, -10), now.AddDate(0, 1, 20), 3, 0),
		subscription(ids.subSessionReady, ids.memberSessionReady, ids.courseAdvanced, ids.branchCentral, courses[2], "active", ptrTime(now.AddDate(0, 0, -5)), nil, now.AddDate(0, 0, -5), now.AddDate(0, 1, 25), 3, 18),
	}
	if err := upsertMany(ctx, db.Collection("subscriptions"), subscriptions); err != nil {
		return nil, fmt.Errorf("subscriptions: %w", err)
	}
	counts["subscriptions"] = len(subscriptions)

	reportedMissedDate := today.AddDate(0, 0, -3).Add(8 * time.Hour)
	makeupDate := today.Add(9 * time.Hour)
	attendances := []models.Attendance{
		{
			ID:       ids.attendanceToday,
			SubID:    ids.subActive,
			BranchID: ids.branchCentral,
			Date:     today.Add(7 * time.Hour),
			Status:   "attended",
		},
		{
			ID:          ids.attendanceReport,
			SubID:       ids.subActive,
			BranchID:    ids.branchCentral,
			Date:        reportedMissedDate,
			Status:      "reported_missed",
			IsMakeupFor: nil,
		},
		{
			ID:          ids.attendanceMakeup,
			SubID:       ids.subActive,
			BranchID:    ids.branchCentral,
			Date:        makeupDate,
			Status:      "makeup",
			IsMakeupFor: &reportedMissedDate,
		},
		{
			ID:        ids.attendanceSession,
			SubID:     ids.subSessionReady,
			BranchID:  ids.branchCentral,
			SessionID: &ids.sessionMorning,
			Date:      today.Add(10 * time.Hour),
			Status:    "attended",
		},
	}
	if err := upsertMany(ctx, db.Collection("attendances"), attendances); err != nil {
		return nil, fmt.Errorf("attendances: %w", err)
	}
	counts["attendances"] = len(attendances)

	sessions := []models.Session{
		{
			ID:                      ids.sessionMorning,
			BranchID:                ids.branchCentral,
			TrainerID:               ids.trainer,
			CourseLevel:             "professional",
			ScheduledAt:             today.Add(10 * time.Hour),
			DurationMin:             75,
			Capacity:                12,
			EnrolledCount:           1,
			EnrolledSubscriptionIDs: []primitive.ObjectID{ids.subSessionReady},
			Tags:                    []string{"barbell", "power"},
		},
		{
			ID:                      ids.sessionEvening,
			BranchID:                ids.branchWest,
			TrainerID:               ids.trainer,
			CourseLevel:             "advanced",
			ScheduledAt:             today.Add(18 * time.Hour),
			DurationMin:             60,
			Capacity:                16,
			EnrolledCount:           1,
			EnrolledSubscriptionIDs: []primitive.ObjectID{ids.subActive},
			Tags:                    []string{"conditioning", "strength"},
		},
		{
			ID:                      ids.sessionFuture,
			BranchID:                ids.branchEast,
			TrainerID:               ids.trainer,
			CourseLevel:             "basic",
			ScheduledAt:             today.AddDate(0, 0, 2).Add(8 * time.Hour),
			DurationMin:             45,
			Capacity:                20,
			EnrolledCount:           0,
			EnrolledSubscriptionIDs: []primitive.ObjectID{},
			Tags:                    []string{"mobility", "beginner"},
		},
	}
	if err := upsertMany(ctx, db.Collection("sessions"), sessions); err != nil {
		return nil, fmt.Errorf("sessions: %w", err)
	}
	counts["sessions"] = len(sessions)

	refunds := []models.Refund{
		{
			ID:                ids.refund,
			SubscriptionID:    ids.subRefunded,
			MemberID:          ids.memberRefunded,
			UsedSessions:      3,
			RemainingSessions: courses[2].SessionCount - 3,
			RefundAmount:      int64(courses[2].SessionCount-3) * courses[2].BasePrice,
			Reason:            "Demo member relocation",
			Status:            models.RefundStatusProcessed,
			CreatedAt:         now.AddDate(0, 0, -7),
			ProcessedAt:       now.AddDate(0, 0, -7),
		},
	}
	if err := upsertMany(ctx, db.Collection("refunds"), refunds); err != nil {
		return nil, fmt.Errorf("refunds: %w", err)
	}
	counts["refunds"] = len(refunds)

	return counts, nil
}

func demoEmployees(ids demoIDs, now time.Time, includeBranches bool) ([]models.Employee, error) {
	managerBranches := []primitive.ObjectID{}
	receptionistBranches := []primitive.ObjectID{}
	trainerBranches := []primitive.ObjectID{}
	if includeBranches {
		managerBranches = []primitive.ObjectID{ids.branchCentral, ids.branchWest, ids.branchEast}
		receptionistBranches = []primitive.ObjectID{ids.branchCentral}
		trainerBranches = []primitive.ObjectID{ids.branchCentral, ids.branchWest, ids.branchEast}
	}

	employees := []models.Employee{
		employee(ids.admin, "ADMIN001", "Gym Admin", "admin@gym.test", []string{service.RoleAdmin}, "", "", []primitive.ObjectID{}, now),
		employee(ids.manager, "MANAGER001", "Central Manager", "manager@gym.test", []string{service.RoleManager}, "", "0901000001", managerBranches, now),
		employee(ids.receptionist, "REC001", "Front Desk Demo", "receptionist@gym.test", []string{service.RoleReceptionist}, "", "0901000002", receptionistBranches, now),
		employee(ids.trainer, "TRAINER001", "Coach Demo", "trainer@gym.test", []string{service.RoleTrainer}, "professional", "0901000003", trainerBranches, now),
	}

	for i := range employees {
		hash, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		employees[i].PasswordHash = string(hash)
	}
	return employees, nil
}

func employee(id primitive.ObjectID, employeeID string, fullName string, email string, roles []string, level string, phone string, branches []primitive.ObjectID, now time.Time) models.Employee {
	normalizedEmail := service.NormalizeEmail(email)
	return models.Employee{
		ID:              id,
		EmployeeID:      employeeID,
		FullName:        fullName,
		Email:           normalizedEmail,
		NormalizedEmail: normalizedEmail,
		Status:          service.EmployeeStatusActive,
		Role:            roles,
		Level:           level,
		Phone:           phone,
		BranchID:        branches,
		CreatedAt:       now.AddDate(0, 0, -30),
		UpdatedAt:       now,
	}
}

func member(id primitive.ObjectID, phone string, ccid string, fullName string, email string, gender string, level string, registered bool, attended int, createdAt time.Time, updatedAt time.Time) models.Member {
	return models.Member{
		ID:                    id,
		CCID:                  ccid,
		FullName:              fullName,
		Email:                 email,
		Phone:                 phone,
		Gender:                gender,
		Level:                 level,
		IsRegistered:          registered,
		IsSuspended:           false,
		TotalSessionsAttended: attended,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}
}

func subscription(id primitive.ObjectID, memberID primitive.ObjectID, courseID primitive.ObjectID, branchID primitive.ObjectID, course models.Course, status string, paymentDate *time.Time, suspension *models.Suspension, startDate time.Time, endDate time.Time, sessionPerWeek int, remaining int) models.Subscription {
	subtotal := course.BasePrice * int64(course.SessionCount)
	return models.Subscription{
		ID:                id,
		MemberID:          memberID,
		CourseID:          courseID,
		HomeBranchID:      branchID,
		AllowedTags:       course.AllowedTags,
		Status:            status,
		PaymentDate:       paymentDate,
		SubtotalAmount:    subtotal,
		DiscountType:      "none",
		DiscountValue:     0,
		DiscountAmount:    0,
		TotalAmountPaid:   subtotal,
		UnitPrice:         course.BasePrice,
		TotalSessions:     course.SessionCount,
		RemainingSessions: remaining,
		StartDate:         startDate,
		EndDate:           endDate,
		SessionPerWeek:    sessionPerWeek,
		Suspension:        suspension,
	}
}

func upsertMany[T any](ctx context.Context, collection *mongo.Collection, docs []T) error {
	for _, doc := range docs {
		raw, err := bson.Marshal(doc)
		if err != nil {
			return err
		}
		var decoded bson.M
		if err := bson.Unmarshal(raw, &decoded); err != nil {
			return err
		}

		id, ok := decoded["_id"]
		if !ok {
			return fmt.Errorf("document for %s has no _id", collection.Name())
		}

		_, err = collection.ReplaceOne(
			ctx,
			bson.M{"_id": id},
			doc,
			options.Replace().SetUpsert(true),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func upsertManyWithFilter[T any](ctx context.Context, collection *mongo.Collection, docs []T, filterFor func(bson.M) bson.M) error {
	for _, doc := range docs {
		raw, err := bson.Marshal(doc)
		if err != nil {
			return err
		}
		var decoded bson.M
		if err := bson.Unmarshal(raw, &decoded); err != nil {
			return err
		}

		id, ok := decoded["_id"]
		if !ok {
			return fmt.Errorf("document for %s has no _id", collection.Name())
		}
		delete(decoded, "_id")

		_, err = collection.UpdateOne(
			ctx,
			filterFor(decoded),
			bson.M{
				"$set":         decoded,
				"$setOnInsert": bson.M{"_id": id},
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func captureEmployeeIDs(ctx context.Context, collection *mongo.Collection, ids *demoIDs) error {
	values := map[string]*primitive.ObjectID{
		"admin@gym.test":        &ids.admin,
		"manager@gym.test":      &ids.manager,
		"receptionist@gym.test": &ids.receptionist,
		"trainer@gym.test":      &ids.trainer,
	}
	for email, target := range values {
		id, err := objectIDByFilter(ctx, collection, bson.M{"normalized_email": email})
		if err != nil {
			return fmt.Errorf("load employee %s id: %w", email, err)
		}
		*target = id
	}
	return nil
}

func captureBranchIDs(ctx context.Context, collection *mongo.Collection, ids *demoIDs) error {
	values := map[string]*primitive.ObjectID{
		"IFG-CEN": &ids.branchCentral,
		"IFG-WST": &ids.branchWest,
		"IFG-EST": &ids.branchEast,
	}
	for branchCode, target := range values {
		id, err := objectIDByFilter(ctx, collection, bson.M{"branch_code": branchCode})
		if err != nil {
			return fmt.Errorf("load branch %s id: %w", branchCode, err)
		}
		*target = id
	}
	return nil
}

func captureCourseIDs(ctx context.Context, collection *mongo.Collection, ids *demoIDs) error {
	values := map[string]*primitive.ObjectID{
		"Foundation Strength": &ids.courseBasic,
		"Hybrid Conditioning": &ids.courseIntermediate,
		"Performance Barbell": &ids.courseAdvanced,
	}
	for title, target := range values {
		id, err := objectIDByFilter(ctx, collection, bson.M{"title": title})
		if err != nil {
			return fmt.Errorf("load course %s id: %w", title, err)
		}
		*target = id
	}
	return nil
}

func captureMemberIDs(ctx context.Context, collection *mongo.Collection, ids *demoIDs) error {
	values := map[string]*primitive.ObjectID{
		"CCID-DEMO-001": &ids.memberActive,
		"CCID-DEMO-002": &ids.memberSuspended,
		"CCID-DEMO-003": &ids.memberPending,
		"CCID-DEMO-004": &ids.memberExpired,
		"CCID-DEMO-005": &ids.memberRefunded,
		"CCID-DEMO-006": &ids.memberSessionReady,
	}
	for ccid, target := range values {
		id, err := objectIDByFilter(ctx, collection, bson.M{"ccid": ccid})
		if err != nil {
			return fmt.Errorf("load member %s id: %w", ccid, err)
		}
		*target = id
	}
	return nil
}

func objectIDByFilter(ctx context.Context, collection *mongo.Collection, filter bson.M) (primitive.ObjectID, error) {
	var result struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	if err := collection.FindOne(ctx, filter).Decode(&result); err != nil {
		return primitive.NilObjectID, err
	}
	if result.ID.IsZero() {
		return primitive.NilObjectID, fmt.Errorf("matched %s document has zero _id", collection.Name())
	}
	return result.ID, nil
}

func ptrTime(value time.Time) *time.Time {
	return &value
}

func oid(hex string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return id
}

func databaseNameFromEnv() string {
	dbName := strings.TrimSpace(os.Getenv("DB_NAME"))
	if dbName == "" {
		return "gym_management"
	}
	return dbName
}

type demoIDs struct {
	admin        primitive.ObjectID
	manager      primitive.ObjectID
	receptionist primitive.ObjectID
	trainer      primitive.ObjectID

	branchCentral primitive.ObjectID
	branchWest    primitive.ObjectID
	branchEast    primitive.ObjectID

	courseBasic        primitive.ObjectID
	courseIntermediate primitive.ObjectID
	courseAdvanced     primitive.ObjectID

	memberActive       primitive.ObjectID
	memberSuspended    primitive.ObjectID
	memberPending      primitive.ObjectID
	memberExpired      primitive.ObjectID
	memberRefunded     primitive.ObjectID
	memberSessionReady primitive.ObjectID

	subActive       primitive.ObjectID
	subSuspended    primitive.ObjectID
	subPending      primitive.ObjectID
	subExpired      primitive.ObjectID
	subRefunded     primitive.ObjectID
	subSessionReady primitive.ObjectID

	attendanceToday   primitive.ObjectID
	attendanceReport  primitive.ObjectID
	attendanceMakeup  primitive.ObjectID
	attendanceSession primitive.ObjectID

	sessionMorning primitive.ObjectID
	sessionEvening primitive.ObjectID
	sessionFuture  primitive.ObjectID

	refund primitive.ObjectID
}
