package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/test-go/testify/require"
)

type TestConfig struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	Port        int    `env:"PORT,default=8080"`
}

func TestLoad_Success(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("PORT", "9090")

	cfg, err := Load[TestConfig](context.Background())
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	require.Equal(t, 9090, cfg.Port)
}

func TestLoad_DefaultValue(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")

	cfg, err := Load[TestConfig](context.Background())
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	require.Equal(t, 8080, cfg.Port)
}

func TestLoad_MissingRequired(t *testing.T) {
	_, err := Load[TestConfig](context.Background())
	require.Error(t, err)
}

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
