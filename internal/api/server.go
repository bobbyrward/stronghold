package api

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

	echoServer.GET("/", func(c echo.Context) error {
		slog.InfoContext(c.Request().Context(), "Health check endpoint called")
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Register all API routes
	RegisterRoutes(echoServer, db)

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

/*
func Run(configFilePath string) error {
	err := loadConfig(configFilePath)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to load config"))
	}

	parentContext, interruptStop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer interruptStop()

	daemon.createEchoServer()
	cancelDaemons, err := daemon.startDaemons(parentContext)
	if err != nil {
		daemon.shutdownEchoServer()
		return errors.Join(err, fmt.Errorf("failed to start daemons"))
	}

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-parentContext.Done()

	cancelDaemons()

	daemon.shutdownEchoServer()

	return nil
}

func  startDaemons(parentContext context.Context) (context.CancelFunc, error) {
	daemonContext, cancelDaemons := context.WithCancel(parentContext)

	slog.Info("Starting daemons...")
	d.echoServer.Logger.Info("Starting daemons...")

	err := d.startBookImporterDaemon(daemonContext)
	if err != nil {
		cancelDaemons()
		return cancelDaemons, err
	}
func Run(configFilePath string) error {
	err := loadConfig(configFilePath)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to load config"))
	}

	parentContext, interruptStop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer interruptStop()

	daemon.createEchoServer()
	cancelDaemons, err := daemon.startDaemons(parentContext)
	if err != nil {
		daemon.shutdownEchoServer()
		return errors.Join(err, fmt.Errorf("failed to start daemons"))
	}

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-parentContext.Done()

	cancelDaemons()

	daemon.shutdownEchoServer()

	return nil
}

func  startDaemons(parentContext context.Context) (context.CancelFunc, error) {
	daemonContext, cancelDaemons := context.WithCancel(parentContext)

	slog.Info("Starting daemons...")
	d.echoServer.Logger.Info("Starting daemons...")

	err := d.startBookImporterDaemon(daemonContext)
	if err != nil {
		cancelDaemons()
		return cancelDaemons, err
	}

	return cancelDaemons, nil
}
*/
