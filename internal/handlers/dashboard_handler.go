package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) Summary(c *gin.Context) {
	filter, ok := h.parseRangeFilter(c)
	if !ok {
		return
	}

	summary, err := h.dashboardService.Summary(c.Request.Context(), filter)
	if err != nil {
		h.handleDashboardError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dashboard summary fetched successfully", "data": summary})
}

func (h *DashboardHandler) Revenue(c *gin.Context) {
	filter, ok := h.parseRevenueFilter(c)
	if !ok {
		return
	}

	revenue, err := h.dashboardService.Revenue(c.Request.Context(), filter)
	if err != nil {
		h.handleDashboardError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dashboard revenue fetched successfully", "data": revenue})
}

func (h *DashboardHandler) PlanDistribution(c *gin.Context) {
	filter, ok := h.parseRangeFilter(c)
	if !ok {
		return
	}

	distribution, err := h.dashboardService.PlanDistribution(c.Request.Context(), filter)
	if err != nil {
		h.handleDashboardError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dashboard plan distribution fetched successfully", "data": distribution})
}

func (h *DashboardHandler) RecentMembers(c *gin.Context) {
	limit, ok := parseOptionalIntQuery(c, "limit")
	if !ok {
		return
	}

	members, err := h.dashboardService.RecentMembers(c.Request.Context(), service.DashboardRecentMembersFilter{Limit: limit})
	if err != nil {
		h.handleDashboardError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dashboard recent members fetched successfully", "data": members})
}

func (h *DashboardHandler) TodaySessions(c *gin.Context) {
	date, ok := parseOptionalTimeQuery(c, "date")
	if !ok {
		return
	}

	sessions, err := h.dashboardService.TodaySessions(c.Request.Context(), service.DashboardTodaySessionsFilter{
		BranchID: c.Query("branch_id"),
		Date:     date,
	})
	if err != nil {
		h.handleDashboardError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "dashboard today sessions fetched successfully", "data": sessions})
}

func (h *DashboardHandler) parseRangeFilter(c *gin.Context) (service.DashboardRangeFilter, bool) {
	from, ok := parseOptionalTimeQuery(c, "from")
	if !ok {
		return service.DashboardRangeFilter{}, false
	}
	to, ok := parseOptionalTimeQuery(c, "to")
	if !ok {
		return service.DashboardRangeFilter{}, false
	}
	return service.DashboardRangeFilter{
		BranchID: c.Query("branch_id"),
		From:     from,
		To:       to,
	}, true
}

func (h *DashboardHandler) parseRevenueFilter(c *gin.Context) (service.DashboardRevenueFilter, bool) {
	rangeFilter, ok := h.parseRangeFilter(c)
	if !ok {
		return service.DashboardRevenueFilter{}, false
	}
	return service.DashboardRevenueFilter{
		BranchID: rangeFilter.BranchID,
		From:     rangeFilter.From,
		To:       rangeFilter.To,
		Bucket:   c.Query("bucket"),
	}, true
}

func (h *DashboardHandler) handleDashboardError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidDashboardID):
		RespondInvalidID(c, "invalid dashboard branch id")
	case errors.Is(err, service.ErrInvalidDashboardDate):
		RespondInvalidDate(c, "invalid dashboard date range")
	case errors.Is(err, service.ErrInvalidDashboardInput):
		RespondInvalidInput(c, err.Error())
	default:
		RespondInternal(c)
	}
}

func parseOptionalTimeQuery(c *gin.Context, key string) (*time.Time, bool) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return nil, true
	}

	parsed, err := parseTimeValue(value)
	if err != nil {
		RespondInvalidDate(c, "invalid "+key+" format")
		return nil, false
	}
	return &parsed, true
}

func parseOptionalIntQuery(c *gin.Context, key string) (int, bool) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return 0, true
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		RespondInvalidInput(c, "invalid "+key)
		return 0, false
	}
	return parsed, true
}
