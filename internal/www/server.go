package www

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
		return errors.Join(err, fmt.Errorf("failed to connect to database"))
	}

	// Run auto migration
	err = models.AutoMigrate(db)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to auto-migrate database", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to automigrate database"))
	}

	echoServer := echo.New()
	echoServer.HideBanner = true

	// Add slog middleware for request logging
	echoServer.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogMethod:   true,
		LogLatency:  true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
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
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	/*
		echoServer.GET("/", func(c echo.Context) error {
			slog.InfoContext(c.Request().Context(), "Health check endpoint called")
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})
	*/

	echoServer.Validator = NewValidator()

	// Register all API routes first (so they take precedence)
	api.RegisterRoutes(echoServer.Group("/api"), db)

	// Serve Vue SPA static files from web/dist
	echoServer.Static("/assets", "web/dist/assets")
	echoServer.File("/stronghold.svg", "web/dist/stronghold.svg")

	// SPA fallback - serve index.html for all non-API routes
	echoServer.GET("/*", func(c echo.Context) error {
		return c.File("web/dist/index.html")
	})

	go func() {
		slog.InfoContext(ctx, "Starting HTTP server on :8000")
		if err := echoServer.Start(":8000"); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "Failed to start HTTP server", slog.Any("err", err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-parentContext.Done()
	slog.InfoContext(ctx, "Received interrupt signal, shutting down server")

	shutdownContext, timeoutFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeoutFunc()

	if err := echoServer.Shutdown(shutdownContext); err != nil {
		slog.ErrorContext(ctx, "Failed to gracefully shutdown server", slog.Any("err", err))
		return err
	}

	slog.InfoContext(ctx, "Server shutdown complete")
	return nil
}
