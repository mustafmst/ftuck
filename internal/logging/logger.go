package logging

import (
	"log/slog"
	"os"
)

// LogLevel represents the different log levels available
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Config holds the logging configuration
type Config struct {
	Level  LogLevel
	Format string // "json" or "text" 
}

// DefaultConfig returns the default logging configuration
func DefaultConfig() *Config {
	return &Config{
		Level:  LevelInfo,
		Format: "text",
	}
}

// InitLogger initializes and configures the structured logger
func InitLogger(config *Config) *slog.Logger {
	var level slog.Level
	switch config.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if config.Format == "json" {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	
	return logger
}

// MustInitLogger initializes the logger and panics if it fails
func MustInitLogger(config *Config) *slog.Logger {
	logger := InitLogger(config)
	return logger
}