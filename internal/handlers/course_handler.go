package handlers

import (
	"errors"
	"net/http"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CourseHandler exposes HTTP endpoints for course management.
type CourseHandler struct {
	courseService service.CourseService
}

// NewCourseHandler wires course service into handler.
func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// createCourseRequest is JSON body for creating/updating a course.
type createCourseRequest struct {
	Title        string `json:"title"`
	Level        string `json:"level"`
	BasePrice    int64  `json:"base_price"`
	SessionCount int    `json:"session_count"`
	Description  string `json:"description"`
}

// Create handles POST /courses.
func (h *CourseHandler) Create(c *gin.Context) {
	// 1) Parse JSON body.
	var req createCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	// 2) Build domain model for service validation.
	course := &models.Course{
		Title:        req.Title,
		Level:        req.Level,
		BasePrice:    req.BasePrice,
		SessionCount: req.SessionCount,
		Description:  req.Description,
	}

	// 3) Delegate to service and map errors.
	if err := h.courseService.CreateCourse(c.Request.Context(), course); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCourseInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	// 4) Success response.
	c.JSON(http.StatusCreated, gin.H{"message": "course created successfully", "data": course})
}

// GetByID handles GET /courses/:id.
func (h *CourseHandler) GetByID(c *gin.Context) {
	// 1) Validate path ID.
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid course id"})
		return
	}

	// 2) Delegate to service.
	course, err := h.courseService.GetCourseByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "course not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	// 3) Success response.
	c.JSON(http.StatusOK, gin.H{"message": "course fetched successfully", "data": course})
}

// List handles GET /courses.
func (h *CourseHandler) List(c *gin.Context) {
	courses, err := h.courseService.ListCourses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "courses fetched successfully", "data": courses})
}

// Update handles PATCH /courses/:id.
func (h *CourseHandler) Update(c *gin.Context) {
	// 1) Validate path ID.
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid course id"})
		return
	}

	// 2) Parse JSON body.
	var req createCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	// 3) Build model and delegate to service.
	course := &models.Course{
		Title:        req.Title,
		Level:        req.Level,
		BasePrice:    req.BasePrice,
		SessionCount: req.SessionCount,
		Description:  req.Description,
	}
	if err := h.courseService.UpdateCourse(c.Request.Context(), id, course); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCourseInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case errors.Is(err, service.ErrCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "course not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "course updated successfully"})
}

// Delete handles DELETE /courses/:id.
func (h *CourseHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid course id"})
		return
	}

	if err := h.courseService.DeleteCourse(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrCourseNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "course not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "course deleted successfully"})
}
