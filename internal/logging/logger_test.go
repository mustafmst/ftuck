package logging

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.Level != LevelInfo {
		t.Errorf("Expected default level to be LevelInfo, got %v", config.Level)
	}
	if config.Format != "text" {
		t.Errorf("Expected default format to be 'text', got %v", config.Format)
	}
}

func TestConfigFromEnv(t *testing.T) {
	tests := []struct {
		name           string
		envLogLevel    string
		envLogFormat   string
		expectedLevel  LogLevel
		expectedFormat string
	}{
		{
			name:           "default values",
			envLogLevel:    "",
			envLogFormat:   "",
			expectedLevel:  LevelInfo,
			expectedFormat: "text",
		},
		{
			name:           "debug level",
			envLogLevel:    "debug",
			envLogFormat:   "",
			expectedLevel:  LevelDebug,
			expectedFormat: "text",
		},
		{
			name:           "error level with json format",
			envLogLevel:    "error",
			envLogFormat:   "json",
			expectedLevel:  LevelError,
			expectedFormat: "json",
		},
		{
			name:           "case insensitive",
			envLogLevel:    "DEBUG",
			envLogFormat:   "JSON",
			expectedLevel:  LevelDebug,
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.envLogLevel != "" {
				os.Setenv("FTUCK_LOG_LEVEL", tt.envLogLevel)
			} else {
				os.Unsetenv("FTUCK_LOG_LEVEL")
			}
			
			if tt.envLogFormat != "" {
				os.Setenv("FTUCK_LOG_FORMAT", tt.envLogFormat)
			} else {
				os.Unsetenv("FTUCK_LOG_FORMAT")
			}

			config := ConfigFromEnv()
			
			if config.Level != tt.expectedLevel {
				t.Errorf("Expected level %v, got %v", tt.expectedLevel, config.Level)
			}
			
			if config.Format != tt.expectedFormat {
				t.Errorf("Expected format %v, got %v", tt.expectedFormat, config.Format)
			}

			// Clean up
			os.Unsetenv("FTUCK_LOG_LEVEL")
			os.Unsetenv("FTUCK_LOG_FORMAT")
		})
	}
}

func TestInitLogger(t *testing.T) {
	config := &Config{
		Level:  LevelDebug,
		Format: "text",
	}
	
	logger := InitLogger(config)
	if logger == nil {
		t.Error("Expected logger to be initialized, got nil")
	}
}