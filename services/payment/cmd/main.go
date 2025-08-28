package main

import (
	"context"
	"log/slog"
	"os/signal"
	"payment-system/pkg/config"
	"payment-system/pkg/db"
	"payment-system/pkg/logger"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	paymentCfg "github.com/walker-16/payment-system/services/payment/internal/config"
	"github.com/walker-16/payment-system/services/payment/internal/handler"
	"github.com/walker-16/payment-system/services/payment/internal/order"
	"github.com/walker-16/payment-system/services/payment/internal/repository"
)

// defaultShutdownTimeout
const defaultShutdownTimeout = 10 * time.Second

func main() {
	// set up context that is cancelled on SIGN/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// initialize logger.
	logger := logger.NewSlogLogger(logger.LoggerConfig{
		Format:    logger.FormatJSON,
		Level:     slog.LevelDebug,
		AddSource: false,
	})

	// load configuration.
	cfg, err := config.Load[paymentCfg.PaymentConfiguration](ctx)
	if err != nil {
		logger.Fatal("failed to load configuration", "error", err)
	}

	// initialize db client.
	dbConfig := db.Config{
		DSN:             cfg.DB.DNS,
		MaxConns:        cfg.DB.MaxConns,
		MinConns:        cfg.DB.MinConns,
		MaxConnIdleTime: cfg.DB.MaxConnIdleTime,
		MaxConnLifetime: cfg.DB.MaxConnLifetime,
		AppName:         paymentCfg.AppName,
	}
	db, err := db.New(ctx, dbConfig)
	if err != nil {
		logger.Fatal("failed to create db client", "error", err)
	}
	defer db.Close()

	// create and run server.
	app := newServer(db, logger)
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("payment server started", "port", cfg.Port)
		serverErr <- app.Listen(":" + cfg.Port)
	}()

	// wait for shutdown signal or server error.
	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			logger.Error("payment server stopped unexpectedly", "error", err)
		}
	}

	// graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("failed to shutdown payment server gracefully", "error", err)
	}

	logger.Info("payment server exited succesfully")
}

func newServer(db db.DB, logger logger.Logger) *fiber.App {
	// create a new Fiber app.
	app := fiber.New()

	// TODO: configurate middelware.
	app.Use(recover.New())
	app.Use(fiberLogger.New())

	// Register routes.
	registerRoutes(app, db, logger)
	return app
}

func registerRoutes(app *fiber.App, db db.DB, logger logger.Logger) {
	v1 := app.Group("/v1")
	orderService := newOrderService()
	repository := repository.NewPaymentRepository(db)
	h := handler.NewPaymentHandler(orderService, repository, logger)
	v1.Post("/payments", h.CreatePayment)
}

// NOTE: A mock implementation of the Order Service is used here, as the actual
// service is out of scope for this exercise. However, the interface and structure
// are defined. In a production system, this would be replaced with a proper
// service client that handles retries and robust error handling when communicating
// with the Order Service.
func newOrderService() order.Service {
	return order.NewMockOrderService(order.MockSuccess)
}
