package testutil

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/app"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/SonBestCodeVien5/gym-management-system/pkg/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultTestMongoURI = "mongodb://admin:password123@localhost:27017/?authSource=admin&directConnection=true"

type TestApp struct {
	Router       *gin.Engine
	DB           *mongo.Database
	Client       *mongo.Client
	AdminEmail   string
	AdminPass    string
	AdminToken   string
	AdminRefresh string
}

func NewTestApp(t *testing.T) *TestApp {
	t.Helper()

	gin.SetMode(gin.TestMode)
	client, db := newTestDatabase(t)
	cfg := app.Config{
		Auth: service.AuthConfig{
			AccessSecret:  "test-access-secret",
			RefreshSecret: "test-refresh-secret",
			AccessTTL:     15 * time.Minute,
			RefreshTTL:    24 * time.Hour,
		},
		BootstrapAdmin: service.BootstrapAdminConfig{
			EmployeeID: "ADMIN001",
			FullName:   "Integration Admin",
			Email:      "admin.integration@gym.test",
			Password:   "admin-password-123",
		},
		CORSOrigins: []string{"http://localhost:5173"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := database.EnsureIndexes(ctx, db); err != nil {
		t.Fatalf("EnsureIndexes() error = %v", err)
	}

	router, err := app.NewRouter(ctx, db, cfg)
	if err != nil {
		t.Fatalf("app.NewRouter() error = %v", err)
	}

	testApp := &TestApp{
		Router:     router,
		DB:         db,
		Client:     client,
		AdminEmail: cfg.BootstrapAdmin.Email,
		AdminPass:  cfg.BootstrapAdmin.Password,
	}
	testApp.AdminToken, testApp.AdminRefresh = testApp.Login(t, testApp.AdminEmail, testApp.AdminPass)
	return testApp
}

func newTestDatabase(t *testing.T) (*mongo.Client, *mongo.Database) {
	t.Helper()

	uri := os.Getenv("GYM_TEST_MONGODB_URI")
	if strings.TrimSpace(uri) == "" {
		uri = defaultTestMongoURI
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(1500*time.Millisecond))
	if err != nil {
		t.Skipf("MongoDB integration tests skipped: connect failed: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		t.Skipf("MongoDB integration tests skipped: ping failed: %v", err)
	}

	dbName := "gym_test_" + primitive.NewObjectID().Hex()
	db := client.Database(dbName)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if strings.HasPrefix(db.Name(), "gym_test_") {
			_ = db.Drop(cleanupCtx)
		}
		_ = client.Disconnect(cleanupCtx)
	})

	return client, db
}
