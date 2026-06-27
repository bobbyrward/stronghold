package www

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/eventlog"
	"github.com/bobbyrward/stronghold/internal/hardcover"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/www/api"
)

func Run() error {
	ctx := context.Background()
	parentContext, interruptStop := signal.NotifyContext(ctx, os.Interrupt)
	defer interruptStop()

	slog.InfoContext(ctx, "Starting API server")

	// Connect to database
	db, err := models.ConnectDB()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", slog.Any("err", err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run auto migration
	err = models.AutoMigrate(db)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to auto-migrate database", slog.Any("err", err))
		return fmt.Errorf("failed to automigrate database: %w", err)
	}

	// Clean up old event logs
	eventlog.Cleanup(ctx, db, 90)

	echoServer := echo.New()

	// Add slog middleware for request logging
	echoServer.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		HandleError: true,
		LogValuesFunc: func(c *echo.Context, values middleware.RequestLoggerValues) error {
			if values.Error == nil {
				slog.InfoContext(c.Request().Context(), "HTTP request",
					slog.String("method", values.Method),
					slog.String("uri", values.URI),
					slog.Int("status", values.Status),
					slog.Duration("latency", values.Latency))
			} else {
				slog.ErrorContext(c.Request().Context(), "HTTP request error",
					slog.String("method", values.Method),
					slog.String("uri", values.URI),
					slog.Int("status", values.Status),
					slog.Duration("latency", values.Latency),
					slog.Any("error", values.Error))
			}
			return nil
		},
	}))

	// Add CORS middleware for frontend development
	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	echoServer.Validator = NewValidator()

	// Create Hardcover client
	hc := hardcover.NewClient(config.Config.Hardcover.ApiToken)

	// Register all API routes first (so they take precedence)
	api.RegisterRoutes(echoServer.Group("/api"), db, hc)

	// Serve Vue SPA static files from web/dist
	echoServer.Static("/assets", "web/dist/assets")
	echoServer.File("/stronghold.svg", "web/dist/stronghold.svg")

	// SPA fallback - serve index.html for all non-API routes
	echoServer.GET("/*", func(c *echo.Context) error {
		return c.File("web/dist/index.html")
	})

	// Start blocks until parentContext is cancelled (interrupt), then gracefully
	// shuts down within GracefulTimeout. Replaces v4's e.Start + e.Shutdown,
	// which were removed in v5 in favor of StartConfig.
	sc := echo.StartConfig{
		Address:         ":8000",
		HideBanner:      true,
		GracefulTimeout: 10 * time.Second,
	}
	slog.InfoContext(ctx, "Starting HTTP server on :8000")
	if err := sc.Start(parentContext, echoServer); err != nil && err != http.ErrServerClosed {
		slog.ErrorContext(ctx, "Server error", slog.Any("err", err))
		return err
	}

	slog.InfoContext(ctx, "Server shutdown complete")
	return nil
}
