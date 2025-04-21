package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Run() error {
	parentContext, interruptStop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer interruptStop()

	echoServer := echo.New()
	echoServer.HideBanner = true
	echoServer.Logger.SetLevel(log.INFO)

	echoServer.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	go func() {
		if err := echoServer.Start(":8000"); err != nil && err != http.ErrServerClosed {
			echoServer.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-parentContext.Done()

	shutdownContext, timeoutFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeoutFunc()

	if err := echoServer.Shutdown(shutdownContext); err != nil {
		echoServer.Logger.Fatal(err)
	}

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
