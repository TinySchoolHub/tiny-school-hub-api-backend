package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Storage  StorageConfig
	RateLimit int
	CORS     CORSConfig
	Log      LogConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URL string
}

// JWTConfig holds JWT token configuration
type JWTConfig struct {
	Secret         string
	AccessExpiry   time.Duration
	RefreshExpiry  time.Duration
}

// StorageConfig holds S3-compatible storage configuration
type StorageConfig struct {
	Endpoint       string
	Region         string
	Bucket         string
	AccessKey      string
	SecretKey      string
	UsePathStyle   bool
	Insecure       bool
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			Secret:         getEnv("JWT_SECRET", ""),
			AccessExpiry:   parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"), 15*time.Minute),
			RefreshExpiry:  parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"), 168*time.Hour),
		},
		Storage: StorageConfig{
			Endpoint:     getEnv("STORAGE_ENDPOINT", ""),
			Region:       getEnv("STORAGE_REGION", "us-east-1"),
			Bucket:       getEnv("STORAGE_BUCKET", ""),
			AccessKey:    getEnv("STORAGE_ACCESS_KEY", ""),
			SecretKey:    getEnv("STORAGE_SECRET_KEY", ""),
			UsePathStyle: parseBool(getEnv("STORAGE_USE_PATH_STYLE", "false")),
			Insecure:     parseBool(getEnv("STORAGE_INSECURE", "false")),
		},
		RateLimit: parseInt(getEnv("RATE_LIMIT", "100")),
		CORS: CORSConfig{
			AllowedOrigins: parseSlice(getEnv("CORS_ALLOWED_ORIGINS", "*")),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if c.Storage.Endpoint == "" {
		return fmt.Errorf("STORAGE_ENDPOINT is required")
	}

	if c.Storage.Bucket == "" {
		return fmt.Errorf("STORAGE_BUCKET is required")
	}

	if c.Storage.AccessKey == "" {
		return fmt.Errorf("STORAGE_ACCESS_KEY is required")
	}

	if c.Storage.SecretKey == "" {
		return fmt.Errorf("STORAGE_SECRET_KEY is required")
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development" || c.Server.Env == "dev"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production" || c.Server.Env == "prod"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}

func parseBool(value string) bool {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return b
}

func parseDuration(value string, defaultDuration time.Duration) time.Duration {
	d, err := time.ParseDuration(value)
	if err != nil {
		return defaultDuration
	}
	return d
}

func parseSlice(value string) []string {
	if value == "" {
		return []string{}
	}
	// Simple split by comma
	result := []string{}
	current := ""
	for _, char := range value {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
