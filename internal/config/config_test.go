package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Success(t *testing.T) {
	// Set required environment variables
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/testdb")
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("STORAGE_ENDPOINT", "localhost:9000")
	os.Setenv("STORAGE_BUCKET", "test-bucket")
	os.Setenv("STORAGE_ACCESS_KEY", "test-access")
	os.Setenv("STORAGE_SECRET_KEY", "test-secret")
	defer cleanupEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Check defaults
	if cfg.Server.Port != "8080" {
		t.Errorf("Server.Port = %v, want 8080", cfg.Server.Port)
	}
	if cfg.Server.Env != "development" {
		t.Errorf("Server.Env = %v, want development", cfg.Server.Env)
	}
}

func TestLoad_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		expectedError string
	}{
		{
			name: "missing DATABASE_URL",
			envVars: map[string]string{
				"JWT_SECRET":         "secret",
				"STORAGE_ENDPOINT":   "localhost:9000",
				"STORAGE_BUCKET":     "bucket",
				"STORAGE_ACCESS_KEY": "access",
				"STORAGE_SECRET_KEY": "secret",
			},
			expectedError: "DATABASE_URL is required",
		},
		{
			name: "missing JWT_SECRET",
			envVars: map[string]string{
				"DATABASE_URL":       "postgres://localhost:5432/db",
				"STORAGE_ENDPOINT":   "localhost:9000",
				"STORAGE_BUCKET":     "bucket",
				"STORAGE_ACCESS_KEY": "access",
				"STORAGE_SECRET_KEY": "secret",
			},
			expectedError: "JWT_SECRET is required",
		},
		{
			name: "missing STORAGE_ENDPOINT",
			envVars: map[string]string{
				"DATABASE_URL":       "postgres://localhost:5432/db",
				"JWT_SECRET":         "secret",
				"STORAGE_BUCKET":     "bucket",
				"STORAGE_ACCESS_KEY": "access",
				"STORAGE_SECRET_KEY": "secret",
			},
			expectedError: "STORAGE_ENDPOINT is required",
		},
		{
			name: "missing STORAGE_BUCKET",
			envVars: map[string]string{
				"DATABASE_URL":       "postgres://localhost:5432/db",
				"JWT_SECRET":         "secret",
				"STORAGE_ENDPOINT":   "localhost:9000",
				"STORAGE_ACCESS_KEY": "access",
				"STORAGE_SECRET_KEY": "secret",
			},
			expectedError: "STORAGE_BUCKET is required",
		},
		{
			name: "missing STORAGE_ACCESS_KEY",
			envVars: map[string]string{
				"DATABASE_URL":       "postgres://localhost:5432/db",
				"JWT_SECRET":         "secret",
				"STORAGE_ENDPOINT":   "localhost:9000",
				"STORAGE_BUCKET":     "bucket",
				"STORAGE_SECRET_KEY": "secret",
			},
			expectedError: "STORAGE_ACCESS_KEY is required",
		},
		{
			name: "missing STORAGE_SECRET_KEY",
			envVars: map[string]string{
				"DATABASE_URL":       "postgres://localhost:5432/db",
				"JWT_SECRET":         "secret",
				"STORAGE_ENDPOINT":   "localhost:9000",
				"STORAGE_BUCKET":     "bucket",
				"STORAGE_ACCESS_KEY": "access",
			},
			expectedError: "STORAGE_SECRET_KEY is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanupEnv()
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer cleanupEnv()

			_, err := Load()
			if err == nil {
				t.Fatal("Load() error = nil, want error")
			}

			if err.Error() != "config validation failed: "+tt.expectedError {
				t.Errorf("Load() error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestLoad_CustomValues(t *testing.T) {
	cleanupEnv()

	// Set all environment variables with custom values
	os.Setenv("PORT", "9090")
	os.Setenv("ENV", "production")
	os.Setenv("DATABASE_URL", "postgres://db:5432/mydb")
	os.Setenv("JWT_SECRET", "my-secret")
	os.Setenv("JWT_ACCESS_EXPIRY", "30m")
	os.Setenv("JWT_REFRESH_EXPIRY", "720h")
	os.Setenv("STORAGE_ENDPOINT", "s3.amazonaws.com")
	os.Setenv("STORAGE_REGION", "eu-west-1")
	os.Setenv("STORAGE_BUCKET", "my-bucket")
	os.Setenv("STORAGE_ACCESS_KEY", "my-access")
	os.Setenv("STORAGE_SECRET_KEY", "my-secret")
	os.Setenv("STORAGE_USE_PATH_STYLE", "true")
	os.Setenv("STORAGE_INSECURE", "true")
	os.Setenv("RATE_LIMIT", "200")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://example.com")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_FORMAT", "text")
	defer cleanupEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Verify all custom values
	if cfg.Server.Port != "9090" {
		t.Errorf("Server.Port = %v, want 9090", cfg.Server.Port)
	}
	if cfg.Server.Env != "production" {
		t.Errorf("Server.Env = %v, want production", cfg.Server.Env)
	}
	if cfg.Database.URL != "postgres://db:5432/mydb" {
		t.Errorf("Database.URL = %v, want postgres://db:5432/mydb", cfg.Database.URL)
	}
	if cfg.JWT.Secret != "my-secret" {
		t.Errorf("JWT.Secret = %v, want my-secret", cfg.JWT.Secret)
	}
	if cfg.JWT.AccessExpiry != 30*time.Minute {
		t.Errorf("JWT.AccessExpiry = %v, want 30m", cfg.JWT.AccessExpiry)
	}
	if cfg.JWT.RefreshExpiry != 720*time.Hour {
		t.Errorf("JWT.RefreshExpiry = %v, want 720h", cfg.JWT.RefreshExpiry)
	}
	if cfg.Storage.Endpoint != "s3.amazonaws.com" {
		t.Errorf("Storage.Endpoint = %v, want s3.amazonaws.com", cfg.Storage.Endpoint)
	}
	if cfg.Storage.Region != "eu-west-1" {
		t.Errorf("Storage.Region = %v, want eu-west-1", cfg.Storage.Region)
	}
	if !cfg.Storage.UsePathStyle {
		t.Error("Storage.UsePathStyle = false, want true")
	}
	if !cfg.Storage.Insecure {
		t.Error("Storage.Insecure = false, want true")
	}
	if cfg.RateLimit != 200 {
		t.Errorf("RateLimit = %v, want 200", cfg.RateLimit)
	}
	if len(cfg.CORS.AllowedOrigins) != 2 {
		t.Errorf("CORS.AllowedOrigins length = %v, want 2", len(cfg.CORS.AllowedOrigins))
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %v, want debug", cfg.Log.Level)
	}
	if cfg.Log.Format != "text" {
		t.Errorf("Log.Format = %v, want text", cfg.Log.Format)
	}
}

func TestIsDevelopment(t *testing.T) {
	tests := []struct {
		env  string
		want bool
	}{
		{"development", true},
		{"dev", true},
		{"production", false},
		{"prod", false},
		{"staging", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					Env: tt.env,
				},
			}
			if got := cfg.IsDevelopment(); got != tt.want {
				t.Errorf("IsDevelopment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		env  string
		want bool
	}{
		{"production", true},
		{"prod", true},
		{"development", false},
		{"dev", false},
		{"staging", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					Env: tt.env,
				},
			}
			if got := cfg.IsProduction(); got != tt.want {
				t.Errorf("IsProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"123", 123},
		{"0", 0},
		{"-456", -456},
		{"invalid", 0},
		{"", 0},
		{"12.34", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := parseInt(tt.input); got != tt.want {
				t.Errorf("parseInt(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"1", true},
		{"false", false},
		{"False", false},
		{"FALSE", false},
		{"0", false},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := parseBool(tt.input); got != tt.want {
				t.Errorf("parseBool(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input           string
		defaultDuration time.Duration
		want            time.Duration
	}{
		{"15m", time.Hour, 15 * time.Minute},
		{"2h", time.Hour, 2 * time.Hour},
		{"30s", time.Minute, 30 * time.Second},
		{"invalid", time.Hour, time.Hour},
		{"", 5 * time.Minute, 5 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := parseDuration(tt.input, tt.defaultDuration); got != tt.want {
				t.Errorf("parseDuration(%q, %v) = %v, want %v", tt.input, tt.defaultDuration, got, tt.want)
			}
		})
	}
}

func TestParseSlice(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{"single", []string{"single"}},
		{"one,two", []string{"one", "two"}},
		{"", []string{}},
		{"a,,b", []string{"a", "b"}},
		{"http://localhost:3000,https://example.com", []string{"http://localhost:3000", "https://example.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseSlice(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("parseSlice(%q) length = %v, want %v", tt.input, len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseSlice(%q)[%d] = %v, want %v", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	defaultValue := "default"

	// Test with environment variable set
	os.Setenv(key, "custom-value")
	defer os.Unsetenv(key)

	if got := getEnv(key, defaultValue); got != "custom-value" {
		t.Errorf("getEnv() = %v, want custom-value", got)
	}

	// Test with environment variable unset
	os.Unsetenv(key)
	if got := getEnv(key, defaultValue); got != defaultValue {
		t.Errorf("getEnv() = %v, want %v", got, defaultValue)
	}
}

func TestValidate(t *testing.T) {
	baseConfig := &Config{
		Database: DatabaseConfig{URL: "postgres://localhost:5432/db"},
		JWT:      JWTConfig{Secret: "secret"},
		Storage: StorageConfig{
			Endpoint:  "localhost:9000",
			Bucket:    "bucket",
			AccessKey: "access",
			SecretKey: "secret",
		},
	}

	tests := []struct {
		name    string
		modify  func(*Config)
		wantErr string
	}{
		{
			name:    "valid config",
			modify:  func(_ *Config) {},
			wantErr: "",
		},
		{
			name: "missing database URL",
			modify: func(c *Config) {
				c.Database.URL = ""
			},
			wantErr: "DATABASE_URL is required",
		},
		{
			name: "missing JWT secret",
			modify: func(c *Config) {
				c.JWT.Secret = ""
			},
			wantErr: "JWT_SECRET is required",
		},
		{
			name: "missing storage endpoint",
			modify: func(c *Config) {
				c.Storage.Endpoint = ""
			},
			wantErr: "STORAGE_ENDPOINT is required",
		},
		{
			name: "missing storage bucket",
			modify: func(c *Config) {
				c.Storage.Bucket = ""
			},
			wantErr: "STORAGE_BUCKET is required",
		},
		{
			name: "missing storage access key",
			modify: func(c *Config) {
				c.Storage.AccessKey = ""
			},
			wantErr: "STORAGE_ACCESS_KEY is required",
		},
		{
			name: "missing storage secret key",
			modify: func(c *Config) {
				c.Storage.SecretKey = ""
			},
			wantErr: "STORAGE_SECRET_KEY is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of base config
			cfg := &Config{
				Database: baseConfig.Database,
				JWT:      baseConfig.JWT,
				Storage:  baseConfig.Storage,
			}
			tt.modify(cfg)

			err := cfg.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			} else {
				if err == nil {
					t.Errorf("Validate() error = nil, want %v", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

// Helper function to cleanup environment variables
func cleanupEnv() {
	envVars := []string{
		"PORT", "ENV", "DATABASE_URL", "JWT_SECRET",
		"JWT_ACCESS_EXPIRY", "JWT_REFRESH_EXPIRY",
		"STORAGE_ENDPOINT", "STORAGE_REGION", "STORAGE_BUCKET",
		"STORAGE_ACCESS_KEY", "STORAGE_SECRET_KEY",
		"STORAGE_USE_PATH_STYLE", "STORAGE_INSECURE",
		"RATE_LIMIT", "CORS_ALLOWED_ORIGINS",
		"LOG_LEVEL", "LOG_FORMAT",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}
}
