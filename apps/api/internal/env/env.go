package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL  string
	Port         string
	LogLevel     string
	MaxPasteSize int
}

func Load() (Config, error) {
	env := Config{
		DatabaseURL: valueOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/capsule?sslmode=disable"),
		Port:        valueOrDefault("PORT", "4000"),
		LogLevel:    valueOrDefault("LOG_LEVEL", "info"),
	}

	port, err := strconv.Atoi(env.Port)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("PORT must be a valid TCP port")
	}
	if err := validateLogLevel(env.LogLevel); err != nil {
		return Config{}, err
	}

	maxPasteStr := valueOrDefault("MAX_PASTE_SIZE", "1048576")
	maxPaste, err := strconv.Atoi(maxPasteStr)
	if err != nil || maxPaste < 1 {
		return Config{}, fmt.Errorf("MAX_PASTE_SIZE must be a positive integer")
	}
	env.MaxPasteSize = maxPaste

	return env, nil
}

func valueOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func validateLogLevel(level string) error {
	switch strings.ToLower(level) {
	case "debug", "info", "warn", "error":
		return nil
	default:
		return fmt.Errorf("LOG_LEVEL must be one of debug, info, warn, error")
	}
}
