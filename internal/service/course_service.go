package service

import (
	"context"
	"errors"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrCourseNotFound     = errors.New("course not found")
	ErrInvalidCourseInput = errors.New("invalid course input")
)

// CourseService defines business operations for course management.
type CourseService interface {
	CreateCourse(ctx context.Context, course *models.Course) error
	GetCourseByID(ctx context.Context, id string) (*models.Course, error)
	ListCourses(ctx context.Context) ([]models.Course, error)
	UpdateCourse(ctx context.Context, id string, course *models.Course) error
	DeleteCourse(ctx context.Context, id string) error
}

type courseServiceImpl struct {
	repo repository.CourseRepository
}

// NewCourseService builds the course service with repository dependency.
func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseServiceImpl{repo: repo}
}

// CreateCourse validates input and creates a course.
func (s *courseServiceImpl) CreateCourse(ctx context.Context, course *models.Course) error {
	if course == nil || course.Title == "" || course.Level == "" || course.BasePrice <= 0 || course.SessionCount <= 0 {
		return ErrInvalidCourseInput
	}

	course.ID = primitive.NewObjectID()
	return s.repo.Create(ctx, course)
}

// GetCourseByID returns a course by ID.
func (s *courseServiceImpl) GetCourseByID(ctx context.Context, id string) (*models.Course, error) {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, err
	}
	return course, nil
}

// ListCourses returns all courses.
func (s *courseServiceImpl) ListCourses(ctx context.Context) ([]models.Course, error) {
	return s.repo.List(ctx)
}

// UpdateCourse validates input and updates the given course by ID.
func (s *courseServiceImpl) UpdateCourse(ctx context.Context, id string, course *models.Course) error {
	if course == nil || course.Title == "" || course.Level == "" || course.BasePrice <= 0 || course.SessionCount <= 0 {
		return ErrInvalidCourseInput
	}

	if err := s.repo.UpdateByID(ctx, id, course); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrCourseNotFound
		}
		return err
	}

	return nil
}

// DeleteCourse removes a course by ID.
func (s *courseServiceImpl) DeleteCourse(ctx context.Context, id string) error {
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrCourseNotFound
		}
		return err
	}
	return nil
}
