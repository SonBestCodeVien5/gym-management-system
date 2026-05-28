package handlers

import (
	"errors"
	"net/http"

	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmployeeHandler struct {
	employeeService service.EmployeeService
}

func NewEmployeeHandler(employeeService service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
}

type employeeCreateRequest struct {
	EmployeeID string   `json:"employee_id"`
	FullName   string   `json:"full_name"`
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	Role       []string `json:"role"`
	Level      string   `json:"level"`
	Phone      string   `json:"phone"`
	BranchID   []string `json:"branch_id"`
	Status     string   `json:"status"`
}

type employeeUpdateRequest struct {
	EmployeeID *string   `json:"employee_id"`
	FullName   *string   `json:"full_name"`
	Email      *string   `json:"email"`
	Role       *[]string `json:"role"`
	Level      *string   `json:"level"`
	Phone      *string   `json:"phone"`
	BranchID   *[]string `json:"branch_id"`
	Status     *string   `json:"status"`
}

type employeePasswordRequest struct {
	Password string `json:"password"`
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var req employeeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondInvalidRequestBody(c)
		return
	}

	branchIDs, err := parseEmployeeBranchIDs(req.BranchID)
	if err != nil {
		RespondInvalidID(c, "invalid branch id")
		return
	}

	employee, err := h.employeeService.CreateEmployee(c.Request.Context(), service.EmployeeCreateInput{
		EmployeeID: req.EmployeeID,
		FullName:   req.FullName,
		Email:      req.Email,
		Password:   req.Password,
		Role:       req.Role,
		Level:      req.Level,
		Phone:      req.Phone,
		BranchID:   branchIDs,
		Status:     req.Status,
	})
	if err != nil {
		writeEmployeeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "employee created successfully", "data": employee})
}

func (h *EmployeeHandler) List(c *gin.Context) {
	var branchID primitive.ObjectID
	if rawBranchID := c.Query("branch_id"); rawBranchID != "" {
		parsed, err := primitive.ObjectIDFromHex(rawBranchID)
		if err != nil {
			RespondInvalidID(c, "invalid branch id")
			return
		}
		branchID = parsed
	}

	employees, err := h.employeeService.ListEmployees(c.Request.Context(), service.EmployeeListInput{
		Role:     c.Query("role"),
		Status:   c.Query("status"),
		BranchID: branchID,
	})
	if err != nil {
		writeEmployeeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employees fetched successfully", "data": employees})
}

func (h *EmployeeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		RespondInvalidID(c, "invalid employee id")
		return
	}

	employee, err := h.employeeService.GetEmployeeByID(c.Request.Context(), id)
	if err != nil {
		writeEmployeeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee fetched successfully", "data": employee})
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		RespondInvalidID(c, "invalid employee id")
		return
	}

	var req employeeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondInvalidRequestBody(c)
		return
	}

	var branchIDs *[]primitive.ObjectID
	if req.BranchID != nil {
		parsed, err := parseEmployeeBranchIDs(*req.BranchID)
		if err != nil {
			RespondInvalidID(c, "invalid branch id")
			return
		}
		branchIDs = &parsed
	}

	actorID, ok := c.Get(AuthEmployeeIDKey)
	if !ok {
		RespondUnauthorized(c, "missing access token")
		return
	}
	actorIDString, ok := actorID.(string)
	if !ok {
		RespondUnauthorized(c, "invalid access token")
		return
	}

	employee, err := h.employeeService.UpdateEmployee(c.Request.Context(), actorIDString, id, service.EmployeeUpdateInput{
		EmployeeID: req.EmployeeID,
		FullName:   req.FullName,
		Email:      req.Email,
		Role:       req.Role,
		Level:      req.Level,
		Phone:      req.Phone,
		BranchID:   branchIDs,
		Status:     req.Status,
	})
	if err != nil {
		writeEmployeeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee updated successfully", "data": employee})
}

func (h *EmployeeHandler) UpdatePassword(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		RespondInvalidID(c, "invalid employee id")
		return
	}

	var req employeePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondInvalidRequestBody(c)
		return
	}

	if err := h.employeeService.UpdatePassword(c.Request.Context(), id, req.Password); err != nil {
		writeEmployeeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee password updated successfully"})
}

func parseEmployeeBranchIDs(rawBranchIDs []string) ([]primitive.ObjectID, error) {
	branchIDs := make([]primitive.ObjectID, 0, len(rawBranchIDs))
	for _, rawBranchID := range rawBranchIDs {
		branchID, err := primitive.ObjectIDFromHex(rawBranchID)
		if err != nil {
			return nil, err
		}
		branchIDs = append(branchIDs, branchID)
	}
	return branchIDs, nil
}

func writeEmployeeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidEmployeeInput):
		RespondInvalidInput(c, err.Error())
	case errors.Is(err, service.ErrEmployeeNotFound):
		RespondNotFound(c, "employee not found")
	case errors.Is(err, service.ErrEmployeeConflict):
		RespondConflict(c, err.Error())
	default:
		RespondInternal(c)
	}
}
