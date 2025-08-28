package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func Load[T any](ctx context.Context) (*T, error) {
	var settings T

	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
		slog.Info("no .env file found, skipping load")
	}

	if err := envconfig.Process(ctx, &settings); err != nil {
		return nil, fmt.Errorf("unable to process environment variables: %w", err)
	}
	return &settings, nil
}
