package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/SonBestCodeVien5/gym-management-system/internal/models"
	"github.com/SonBestCodeVien5/gym-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidDashboardInput = errors.New("invalid dashboard input")
	ErrInvalidDashboardID    = errors.New("invalid dashboard id")
	ErrInvalidDashboardDate  = errors.New("invalid dashboard date")
)

const (
	defaultDashboardRecentLimit = 5
	maxDashboardRecentLimit     = 20
	dashboardBucketDay          = "day"
)

type DashboardRangeFilter struct {
	BranchID string
	From     *time.Time
	To       *time.Time
}

type DashboardRevenueFilter struct {
	BranchID string
	From     *time.Time
	To       *time.Time
	Bucket   string
}

type DashboardRecentMembersFilter struct {
	Limit int
}

type DashboardTodaySessionsFilter struct {
	BranchID string
	Date     *time.Time
}

type DashboardService interface {
	Summary(ctx context.Context, filter DashboardRangeFilter) (*models.DashboardSummary, error)
	Revenue(ctx context.Context, filter DashboardRevenueFilter) (*models.DashboardRevenueResponse, error)
	PlanDistribution(ctx context.Context, filter DashboardRangeFilter) (*models.DashboardPlanDistributionResponse, error)
	RecentMembers(ctx context.Context, filter DashboardRecentMembersFilter) (*models.DashboardRecentMembersResponse, error)
	TodaySessions(ctx context.Context, filter DashboardTodaySessionsFilter) (*models.DashboardTodaySessionsResponse, error)
}

type dashboardServiceImpl struct {
	repo DashboardRepository
	now  func() time.Time
}

type DashboardRepository interface {
	CountActiveMembers(ctx context.Context) (int64, error)
	CountRegisteredMembersCreated(ctx context.Context, from time.Time, to time.Time) (int64, error)
	NetRevenueTotal(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error)
	RevenueBuckets(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) ([]models.DashboardRevenueBucket, error)
	CountCheckins(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error)
	CountSessions(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) (int64, error)
	PlanDistribution(ctx context.Context, from *time.Time, to *time.Time, branchID *primitive.ObjectID) ([]models.DashboardPlanDistributionItem, error)
	RecentMembers(ctx context.Context, limit int) ([]models.DashboardRecentMember, error)
	TodaySessions(ctx context.Context, from time.Time, to time.Time, branchID *primitive.ObjectID) ([]models.Session, error)
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardServiceImpl{
		repo: repo,
		now:  time.Now,
	}
}

func (s *dashboardServiceImpl) Summary(ctx context.Context, filter DashboardRangeFilter) (*models.DashboardSummary, error) {
	branchID, err := parseOptionalObjectID(filter.BranchID)
	if err != nil {
		return nil, ErrInvalidDashboardID
	}
	from, to, err := s.resolveSummaryRange(filter.From, filter.To)
	if err != nil {
		return nil, err
	}
	previousFrom, previousTo := previousRange(from, to)
	todayStart := startOfDayUTC(to)
	previousDayStart := todayStart.AddDate(0, 0, -1)
	weekStart := startOfWeekUTC(to)
	previousWeekStart := weekStart.AddDate(0, 0, -7)

	activeMembers, err := s.repo.CountActiveMembers(ctx)
	if err != nil {
		return nil, err
	}
	currentMembers, err := s.repo.CountRegisteredMembersCreated(ctx, from, to)
	if err != nil {
		return nil, err
	}
	previousMembers, err := s.repo.CountRegisteredMembersCreated(ctx, previousFrom, previousTo)
	if err != nil {
		return nil, err
	}
	currentRevenue, err := s.repo.NetRevenueTotal(ctx, from, to, branchID)
	if err != nil {
		return nil, err
	}
	previousRevenue, err := s.repo.NetRevenueTotal(ctx, previousFrom, previousTo, branchID)
	if err != nil {
		return nil, err
	}
	todayCheckins, err := s.repo.CountCheckins(ctx, todayStart, todayStart.AddDate(0, 0, 1), branchID)
	if err != nil {
		return nil, err
	}
	previousDayCheckins, err := s.repo.CountCheckins(ctx, previousDayStart, todayStart, branchID)
	if err != nil {
		return nil, err
	}
	classesThisWeek, err := s.repo.CountSessions(ctx, weekStart, weekStart.AddDate(0, 0, 7), branchID)
	if err != nil {
		return nil, err
	}
	previousWeekClasses, err := s.repo.CountSessions(ctx, previousWeekStart, weekStart, branchID)
	if err != nil {
		return nil, err
	}

	return &models.DashboardSummary{
		ActiveMembers:        activeMembers,
		ActiveMembersDelta:   currentMembers - previousMembers,
		MonthlyRevenue:       currentRevenue,
		MonthlyRevenueDelta:  currentRevenue - previousRevenue,
		TodayCheckins:        todayCheckins,
		TodayCheckinsDelta:   todayCheckins - previousDayCheckins,
		ClassesThisWeek:      classesThisWeek,
		ClassesThisWeekDelta: classesThisWeek - previousWeekClasses,
		Range: models.DashboardRange{
			From: from,
			To:   to,
		},
	}, nil
}

func (s *dashboardServiceImpl) Revenue(ctx context.Context, filter DashboardRevenueFilter) (*models.DashboardRevenueResponse, error) {
	branchID, err := parseOptionalObjectID(filter.BranchID)
	if err != nil {
		return nil, ErrInvalidDashboardID
	}
	bucket := strings.TrimSpace(filter.Bucket)
	if bucket == "" {
		bucket = dashboardBucketDay
	}
	if bucket != dashboardBucketDay {
		return nil, ErrInvalidDashboardInput
	}

	from, to, err := s.resolveRevenueRange(filter.From, filter.To)
	if err != nil {
		return nil, err
	}

	buckets, err := s.repo.RevenueBuckets(ctx, from, to, branchID)
	if err != nil {
		return nil, err
	}

	return &models.DashboardRevenueResponse{
		Bucket: bucket,
		From:   from,
		To:     to,
		Items:  fillDailyRevenueBuckets(from, to, buckets),
	}, nil
}

func (s *dashboardServiceImpl) PlanDistribution(ctx context.Context, filter DashboardRangeFilter) (*models.DashboardPlanDistributionResponse, error) {
	branchID, err := parseOptionalObjectID(filter.BranchID)
	if err != nil {
		return nil, ErrInvalidDashboardID
	}

	var from *time.Time
	var to *time.Time
	if filter.From != nil || filter.To != nil {
		resolvedFrom, resolvedTo, err := s.resolveSummaryRange(filter.From, filter.To)
		if err != nil {
			return nil, err
		}
		from = &resolvedFrom
		to = &resolvedTo
	}

	items, err := s.repo.PlanDistribution(ctx, from, to, branchID)
	if err != nil {
		return nil, err
	}
	return &models.DashboardPlanDistributionResponse{Items: items}, nil
}

func (s *dashboardServiceImpl) RecentMembers(ctx context.Context, filter DashboardRecentMembersFilter) (*models.DashboardRecentMembersResponse, error) {
	limit := filter.Limit
	if limit == 0 {
		limit = defaultDashboardRecentLimit
	}
	if limit < 1 || limit > maxDashboardRecentLimit {
		return nil, ErrInvalidDashboardInput
	}

	items, err := s.repo.RecentMembers(ctx, limit)
	if err != nil {
		return nil, err
	}
	return &models.DashboardRecentMembersResponse{Items: items}, nil
}

func (s *dashboardServiceImpl) TodaySessions(ctx context.Context, filter DashboardTodaySessionsFilter) (*models.DashboardTodaySessionsResponse, error) {
	branchID, err := parseOptionalObjectID(filter.BranchID)
	if err != nil {
		return nil, ErrInvalidDashboardID
	}

	date := s.now().UTC()
	if filter.Date != nil {
		date = filter.Date.UTC()
	}
	start := startOfDayUTC(date)

	items, err := s.repo.TodaySessions(ctx, start, start.AddDate(0, 0, 1), branchID)
	if err != nil {
		return nil, err
	}
	return &models.DashboardTodaySessionsResponse{Items: items}, nil
}

func (s *dashboardServiceImpl) resolveSummaryRange(fromInput *time.Time, toInput *time.Time) (time.Time, time.Time, error) {
	now := s.now().UTC()
	to := now
	if toInput != nil {
		to = toInput.UTC()
	}

	from := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC)
	if fromInput != nil {
		from = fromInput.UTC()
	}

	if !from.Before(to) {
		return time.Time{}, time.Time{}, ErrInvalidDashboardDate
	}
	return from, to, nil
}

func (s *dashboardServiceImpl) resolveRevenueRange(fromInput *time.Time, toInput *time.Time) (time.Time, time.Time, error) {
	now := s.now().UTC()
	to := now
	if toInput != nil {
		to = toInput.UTC()
	}

	from := startOfDayUTC(to).AddDate(0, 0, -6)
	if fromInput != nil {
		from = fromInput.UTC()
	}

	if !from.Before(to) {
		return time.Time{}, time.Time{}, ErrInvalidDashboardDate
	}
	return from, to, nil
}

func parseOptionalObjectID(value string) (*primitive.ObjectID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	objectID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return nil, err
	}
	return &objectID, nil
}

func previousRange(from time.Time, to time.Time) (time.Time, time.Time) {
	duration := to.Sub(from)
	return from.Add(-duration), from
}

func startOfDayUTC(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

func startOfWeekUTC(value time.Time) time.Time {
	dayStart := startOfDayUTC(value)
	offset := int(dayStart.Weekday() - time.Monday)
	if offset < 0 {
		offset += 7
	}
	return dayStart.AddDate(0, 0, -offset)
}

func fillDailyRevenueBuckets(from time.Time, to time.Time, buckets []models.DashboardRevenueBucket) []models.DashboardRevenueBucket {
	byLabel := make(map[string]models.DashboardRevenueBucket, len(buckets))
	for _, bucket := range buckets {
		byLabel[bucket.Label] = bucket
	}

	start := startOfDayUTC(from)
	end := startOfDayUTC(to)
	if !end.Before(to) {
		end = end.AddDate(0, 0, 1)
	}

	items := []models.DashboardRevenueBucket{}
	for cursor := start; cursor.Before(end); cursor = cursor.AddDate(0, 0, 1) {
		label := cursor.Format("2006-01-02")
		bucket := byLabel[label]
		bucket.Label = label
		bucket.NetAmount = bucket.GrossAmount - bucket.RefundAmount
		items = append(items, bucket)
	}
	return items
}
