package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/test-go/testify/require"
)

// TestConfig represents a test configuration struct.
type TestConfig struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	Port        int    `env:"PORT,default=8080"`
}

// TestLoad_Success checks that Load[T] can correctly load the configuration
// from explicitly set environment variables.
func TestLoad_Success(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("PORT", "9090")

	cfg, err := Load[TestConfig](context.Background())
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	require.Equal(t, 9090, cfg.Port)
}

// TestLoad_DefaultValue checks that default values are applied correctly
// when the environment variable is not defined.
func TestLoad_DefaultValue(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")

	cfg, err := Load[TestConfig](context.Background())
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	require.Equal(t, 8080, cfg.Port)
}

// TestLoad_MissingRequired checks that an error is returned
// when a required environment variable is missing.
func TestLoad_MissingRequired(t *testing.T) {
	_, err := Load[TestConfig](context.Background())
	require.Error(t, err)
}

// TestLoad_FromEnvFile checks that Load[T] can load configuration
// from a .env file in the current working directory.
func TestLoad_FromEnvFile(t *testing.T) {
	// Creamos archivo temporal .env
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	content := []byte("DATABASE_URL=postgres://fileuser:filepass@localhost:5432/filedb\nPORT=7070\n")
	require.NoError(t, os.WriteFile(envFile, content, 0644))

	// Cambiamos cwd temporalmente para que godotenv encuentre el archivo
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(oldWd)

	cfg, err := Load[TestConfig](context.Background())
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "postgres://fileuser:filepass@localhost:5432/filedb", cfg.DatabaseURL)
	require.Equal(t, 7070, cfg.Port)
}
