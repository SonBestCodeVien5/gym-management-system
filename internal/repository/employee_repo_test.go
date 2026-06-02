package repository

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const testMongoURI = "mongodb://admin:password123@localhost:27017/?authSource=admin&directConnection=true"

func TestEmployeeRepositoryBootstrapAdminUpsertsExistingRecord(t *testing.T) {
	ctx := context.Background()
	client, db := openTestEmployeeRepoDatabase(t)
	defer func() { _ = client.Disconnect(context.Background()) }()

	if err := database.EnsureIndexes(ctx, db); err != nil {
		t.Fatalf("EnsureIndexes() error = %v", err)
	}

	repo := &employeeRepoImpl{collection: db.Collection("employees")}

	oldPasswordHash, err := bcrypt.GenerateFromPassword([]byte("old-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	existing := &models.Employee{
		ID:              primitive.NewObjectID(),
		EmployeeID:      "OLD001",
		FullName:        "Old Admin",
		Email:           "admin@gym.test",
		NormalizedEmail: "admin@gym.test",
		PasswordHash:    string(oldPasswordHash),
		Status:          "active",
		Role:            []string{"admin"},
		BranchID:        []primitive.ObjectID{},
		CreatedAt:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if _, err := db.Collection("employees").InsertOne(ctx, existing); err != nil {
		t.Fatalf("InsertOne() error = %v", err)
	}

	bootstrapPasswordHash, err := bcrypt.GenerateFromPassword([]byte("new-bootstrap-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	bootstrapAdmin := &models.Employee{
		EmployeeID:      "ADMIN001",
		FullName:        "Gym Admin",
		Email:           "admin@gym.test",
		NormalizedEmail: "admin@gym.test",
		PasswordHash:    string(bootstrapPasswordHash),
		Status:          "active",
		Role:            []string{"admin"},
		Phone:           "",
		Level:           "",
		BranchID:        []primitive.ObjectID{},
	}

	if err := repo.BootstrapAdmin(ctx, bootstrapAdmin); err != nil {
		t.Fatalf("BootstrapAdmin() error = %v", err)
	}

	var saved models.Employee
	if err := db.Collection("employees").FindOne(ctx, bson.M{"normalized_email": "admin@gym.test"}).Decode(&saved); err != nil {
		t.Fatalf("FindOne() error = %v", err)
	}

	if saved.EmployeeID != "ADMIN001" {
		t.Fatalf("EmployeeID = %q, want ADMIN001", saved.EmployeeID)
	}
	if saved.FullName != "Gym Admin" {
		t.Fatalf("FullName = %q, want Gym Admin", saved.FullName)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(saved.PasswordHash), []byte("new-bootstrap-password")); err != nil {
		t.Fatalf("saved password hash does not match bootstrap password: %v", err)
	}
	if saved.ID != existing.ID {
		t.Fatalf("ID = %s, want existing ID %s", saved.ID.Hex(), existing.ID.Hex())
	}
	if !saved.CreatedAt.Equal(existing.CreatedAt) {
		t.Fatalf("CreatedAt = %v, want %v", saved.CreatedAt, existing.CreatedAt)
	}
	if saved.UpdatedAt.IsZero() {
		t.Fatal("UpdatedAt should be set")
	}
}

func openTestEmployeeRepoDatabase(t *testing.T) (*mongo.Client, *mongo.Database) {
	t.Helper()

	uri := strings.TrimSpace(os.Getenv("GYM_TEST_MONGODB_URI"))
	if uri == "" {
		uri = testMongoURI
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetServerSelectionTimeout(1500*time.Millisecond))
	if err != nil {
		t.Skipf("MongoDB integration tests skipped: connect failed: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		t.Skipf("MongoDB integration tests skipped: ping failed: %v", err)
	}

	db := client.Database("gym_repo_test_" + primitive.NewObjectID().Hex())
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = db.Drop(cleanupCtx)
	})

	return client, db
}
