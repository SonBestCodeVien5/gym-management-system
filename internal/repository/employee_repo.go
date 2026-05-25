package repository

import (
	"context"
	"errors"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmployeeRepository interface {
	Create(ctx context.Context, employee *models.Employee) error
	GetByID(ctx context.Context, id string) (*models.Employee, error)
	GetByNormalizedEmail(ctx context.Context, email string) (*models.Employee, error)
	BootstrapAdmin(ctx context.Context, employee *models.Employee) error
}

type employeeRepoImpl struct {
	collection *mongo.Collection
}

func NewEmployeeRepository(db *mongo.Database) (EmployeeRepository, error) {
	collection := db.Collection("employees")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "normalized_email", Value: 1}},
			Options: options.Index().SetName("normalized_email_unique").SetUnique(true).SetSparse(true),
		},
		{
			Keys:    bson.D{{Key: "employee_id", Value: 1}},
			Options: options.Index().SetName("employee_id_unique").SetUnique(true).SetSparse(true),
		},
	}
	if _, err := collection.Indexes().CreateMany(ctx, indexes); err != nil {
		return nil, err
	}

	return &employeeRepoImpl{collection: collection}, nil
}

func (r *employeeRepoImpl) Create(ctx context.Context, employee *models.Employee) error {
	_, err := r.collection.InsertOne(ctx, employee)
	return err
}

func (r *employeeRepoImpl) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var employee models.Employee
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&employee)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepoImpl) GetByNormalizedEmail(ctx context.Context, email string) (*models.Employee, error) {
	var employee models.Employee
	err := r.collection.FindOne(ctx, bson.M{"normalized_email": email}).Decode(&employee)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepoImpl) BootstrapAdmin(ctx context.Context, employee *models.Employee) error {
	if employee == nil {
		return nil
	}

	now := time.Now()
	if employee.ID.IsZero() {
		employee.ID = primitive.NewObjectID()
	}
	employee.CreatedAt = now
	employee.UpdatedAt = now

	filter := bson.M{"normalized_email": employee.NormalizedEmail}
	update := bson.M{"$setOnInsert": employee}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}
