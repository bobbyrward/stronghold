package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type EventLogResponse struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	Category   string    `json:"category"`
	EventType  string    `json:"event_type"`
	Source     string    `json:"source"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Summary    string    `json:"summary"`
	Details    string    `json:"details"`
}

type PaginatedEventLogResponse struct {
	Items   []EventLogResponse `json:"items"`
	Total   int64              `json:"total"`
	Page    int                `json:"page"`
	PerPage int                `json:"per_page"`
	Facets  EventLogFacets     `json:"facets"`
}

type EventLogFacets struct {
	Categories  []string `json:"categories"`
	Sources     []string `json:"sources"`
	EventTypes  []string `json:"event_types"`
	EntityTypes []string `json:"entity_types"`
}

func ListEventLogs(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		slog.InfoContext(ctx, "Listing event logs")

		// Parse pagination params
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}
		perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
		if perPage < 1 {
			perPage = 50
		}
		if perPage > 200 {
			perPage = 200
		}

		// Parse filter params
		category := c.QueryParam("category")
		source := c.QueryParam("source")
		eventType := c.QueryParam("event_type")
		entityType := c.QueryParam("entity_type")
		entityID := c.QueryParam("entity_id")
		q := c.QueryParam("q")
		fromStr := c.QueryParam("from")
		toStr := c.QueryParam("to")

		var fromTime, toTime time.Time
		var hasFrom, hasTo bool
		if fromStr != "" {
			if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
				fromTime = t
				hasFrom = true
			}
		}
		if toStr != "" {
			if t, err := time.Parse(time.RFC3339, toStr); err == nil {
				toTime = t
				hasTo = true
			}
		}

		// Build base query with filters
		query := db.Model(&models.EventLog{})
		if category != "" {
			query = query.Where("category = ?", category)
		}
		if source != "" {
			query = query.Where("source = ?", source)
		}
		if eventType != "" {
			query = query.Where("event_type = ?", eventType)
		}
		if entityType != "" {
			query = query.Where("entity_type = ?", entityType)
		}
		if entityID != "" {
			query = query.Where("entity_id = ?", entityID)
		}
		if q != "" {
			query = query.Where("summary ILIKE ?", "%"+EscapeLikePattern(q)+"%")
		}
		if hasFrom {
			query = query.Where("created_at >= ?", fromTime)
		}
		if hasTo {
			query = query.Where("created_at <= ?", toTime)
		}

		// Count total matching rows
		var total int64
		if err := query.Count(&total).Error; err != nil {
			slog.ErrorContext(ctx, "Failed to count event logs", slog.Any("error", err))
			return InternalError(c, ctx, "Failed to count event logs", err)
		}

		// Fetch paginated results
		var logs []models.EventLog
		offset := (page - 1) * perPage
		if err := query.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&logs).Error; err != nil {
			slog.ErrorContext(ctx, "Failed to list event logs", slog.Any("error", err))
			return InternalError(c, ctx, "Failed to list event logs", err)
		}

		// Build facets (filtered by date range only, not by other filters)
		facetQuery := db.Model(&models.EventLog{})
		if hasFrom {
			facetQuery = facetQuery.Where("created_at >= ?", fromTime)
		}
		if hasTo {
			facetQuery = facetQuery.Where("created_at <= ?", toTime)
		}

		facets := EventLogFacets{}
		facetQuery.Distinct("category").Pluck("category", &facets.Categories)
		facetQuery.Distinct("source").Pluck("source", &facets.Sources)
		facetQuery.Distinct("event_type").Pluck("event_type", &facets.EventTypes)
		facetQuery.Distinct("entity_type").Pluck("entity_type", &facets.EntityTypes)

		// Ensure non-nil slices for JSON
		if facets.Categories == nil {
			facets.Categories = []string{}
		}
		if facets.Sources == nil {
			facets.Sources = []string{}
		}
		if facets.EventTypes == nil {
			facets.EventTypes = []string{}
		}
		if facets.EntityTypes == nil {
			facets.EntityTypes = []string{}
		}

		// Build response
		items := make([]EventLogResponse, len(logs))
		for i, l := range logs {
			items[i] = EventLogResponse{
				ID:         l.ID,
				CreatedAt:  l.CreatedAt,
				Category:   l.Category,
				EventType:  l.EventType,
				Source:     l.Source,
				EntityType: l.EntityType,
				EntityID:   l.EntityID,
				Summary:    l.Summary,
				Details:    l.Details,
			}
		}

		slog.InfoContext(ctx, "Successfully listed event logs",
			slog.Int64("total", total),
			slog.Int("page", page),
			slog.Int("per_page", perPage),
			slog.Int("returned", len(items)),
		)

		return c.JSON(http.StatusOK, PaginatedEventLogResponse{
			Items:   items,
			Total:   total,
			Page:    page,
			PerPage: perPage,
			Facets:  facets,
		})
	}
}

func GetEventLog(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid ID")
		}

		slog.InfoContext(ctx, "Getting event log", slog.Uint64("id", uint64(id)))

		var log models.EventLog
		if err := GetByID(db, ctx, &log, id, "EventLog"); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, "EventLog", id)
			}
			return InternalError(c, ctx, "Failed to get event log", err)
		}

		slog.InfoContext(ctx, "Successfully retrieved event log", slog.Uint64("id", uint64(id)))

		return c.JSON(http.StatusOK, EventLogResponse{
			ID:         log.ID,
			CreatedAt:  log.CreatedAt,
			Category:   log.Category,
			EventType:  log.EventType,
			Source:     log.Source,
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			Summary:    log.Summary,
			Details:    log.Details,
		})
	}
}
