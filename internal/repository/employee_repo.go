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
	List(ctx context.Context, filter EmployeeListFilter) ([]models.Employee, error)
	UpdateByID(ctx context.Context, id string, update EmployeeUpdate) (*models.Employee, error)
	UpdatePasswordByID(ctx context.Context, id string, passwordHash string, updatedAt time.Time) error
	BootstrapAdmin(ctx context.Context, employee *models.Employee) error
}

type EmployeeListFilter struct {
	Role     string
	Status   string
	BranchID primitive.ObjectID
}

type EmployeeUpdate struct {
	EmployeeID      *string
	FullName        *string
	Email           *string
	NormalizedEmail *string
	Status          *string
	Role            *[]string
	Level           *string
	Phone           *string
	BranchID        *[]primitive.ObjectID
	UpdatedAt       time.Time
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
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicate
	}
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

func (r *employeeRepoImpl) List(ctx context.Context, filter EmployeeListFilter) ([]models.Employee, error) {
	query := bson.M{}
	if filter.Role != "" {
		query["role"] = filter.Role
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if !filter.BranchID.IsZero() {
		query["branch_id"] = filter.BranchID
	}

	cursor, err := r.collection.Find(ctx, query, options.Find().SetSort(bson.D{
		{Key: "created_at", Value: -1},
		{Key: "_id", Value: -1},
	}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var employees []models.Employee
	if err := cursor.All(ctx, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *employeeRepoImpl) UpdateByID(ctx context.Context, id string, update EmployeeUpdate) (*models.Employee, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	set := bson.M{}
	if update.EmployeeID != nil {
		set["employee_id"] = *update.EmployeeID
	}
	if update.FullName != nil {
		set["full_name"] = *update.FullName
	}
	if update.Email != nil {
		set["email"] = *update.Email
	}
	if update.NormalizedEmail != nil {
		set["normalized_email"] = *update.NormalizedEmail
	}
	if update.Status != nil {
		set["status"] = *update.Status
	}
	if update.Role != nil {
		set["role"] = *update.Role
	}
	if update.Level != nil {
		set["level"] = *update.Level
	}
	if update.Phone != nil {
		set["phone"] = *update.Phone
	}
	if update.BranchID != nil {
		set["branch_id"] = *update.BranchID
	}
	set["updated_at"] = update.UpdatedAt

	var employee models.Employee
	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": set},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&employee)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepoImpl) UpdatePasswordByID(ctx context.Context, id string, passwordHash string, updatedAt time.Time) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{
		"password_hash": passwordHash,
		"updated_at":    updatedAt,
	}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
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
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicate
	}
	return err
}
