package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BranchHandler exposes HTTP endpoints for branch management.
type BranchHandler struct {
	branchService service.BranchService
}

// NewBranchHandler wires branch service into handler.
func NewBranchHandler(branchService service.BranchService) *BranchHandler {
	return &BranchHandler{branchService: branchService}
}

// branchRequest is JSON body for creating/updating a branch.
type branchRequest struct {
	BranchCode string   `json:"branch_code"`
	Name       string   `json:"name"`
	Address    string   `json:"address"`
	Province   string   `json:"province"`
	Location   geoInput `json:"location"`
	ManagerID  string   `json:"manager_id"`
}

type geoInput struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// Create handles POST /branches.
func (h *BranchHandler) Create(c *gin.Context) {
	// 1) Parse JSON body.
	var req branchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	// 2) Parse optional manager_id.
	managerID := primitive.NilObjectID
	if req.ManagerID != "" {
		parsed, err := primitive.ObjectIDFromHex(req.ManagerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid manager id"})
			return
		}
		managerID = parsed
	}

	// 3) Build domain model and delegate to service.
	branch := &models.Branch{
		BranchCode: req.BranchCode,
		Name:       req.Name,
		Address:    req.Address,
		Province:   req.Province,
		Location: models.GeoLocation{
			Type:        req.Location.Type,
			Coordinates: req.Location.Coordinates,
		},
		ManagerID: managerID,
	}

	if err := h.branchService.CreateBranch(c.Request.Context(), branch); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidBranchInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "branch created successfully", "data": branch})
}

// GetByID handles GET /branches/:id.
func (h *BranchHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch id"})
		return
	}

	branch, err := h.branchService.GetBranchByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBranchNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "branch not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "branch fetched successfully", "data": branch})
}

// List handles GET /branches.
func (h *BranchHandler) List(c *gin.Context) {
	branches, err := h.branchService.ListBranches(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "branches fetched successfully", "data": branches})
}

// Nearby handles GET /branches/nearby.
func (h *BranchHandler) Nearby(c *gin.Context) {
	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch input"})
		return
	}

	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch input"})
		return
	}

	var maxDistance int64
	if c.Query("max_distance") != "" {
		maxDistance, err = strconv.ParseInt(c.Query("max_distance"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch input"})
			return
		}
		if maxDistance <= 0 {
			maxDistance = -1
		}
	}

	var limit int64
	if c.Query("limit") != "" {
		limit, err = strconv.ParseInt(c.Query("limit"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch input"})
			return
		}
	}

	branches, err := h.branchService.NearbyBranches(c.Request.Context(), lng, lat, maxDistance, limit)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidBranchInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "nearby branches fetched successfully", "data": branches})
}

// Update handles PATCH /branches/:id.
func (h *BranchHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch id"})
		return
	}

	var req branchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	managerID := primitive.NilObjectID
	if req.ManagerID != "" {
		parsed, err := primitive.ObjectIDFromHex(req.ManagerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid manager id"})
			return
		}
		managerID = parsed
	}

	branch := &models.Branch{
		BranchCode: req.BranchCode,
		Name:       req.Name,
		Address:    req.Address,
		Province:   req.Province,
		Location: models.GeoLocation{
			Type:        req.Location.Type,
			Coordinates: req.Location.Coordinates,
		},
		ManagerID: managerID,
	}

	if err := h.branchService.UpdateBranch(c.Request.Context(), id, branch); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidBranchInput):
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		case errors.Is(err, service.ErrBranchNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "branch not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "branch updated successfully"})
}

// Delete handles DELETE /branches/:id.
func (h *BranchHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid branch id"})
		return
	}

	if err := h.branchService.DeleteBranch(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrBranchNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "branch not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "branch deleted successfully"})
}
