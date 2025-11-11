package postalcode

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the configuration for postalcode package
type Config struct {
	Database DatabaseConfig
	Import   ImportConfig
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Charset  string
}

// ImportConfig holds data import settings
type ImportConfig struct {
	BatchSize      int
	DefaultTimeout int // seconds
}

// LoadConfig loads configuration from environment variables
// It attempts to load .env file from multiple locations:
// 1. Current directory
// 2. Parent directory (useful when running from cmd/ subdirectories)
// 3. Two levels up (useful when running from tmp/ directory with air)
func LoadConfig() (*Config, error) {
	// Try to load .env file from multiple locations
	envPaths := []string{
		".env",                                   // Current directory
		"../.env",                                // Parent directory
		"../../.env",                             // Two levels up
		filepath.Join(os.Getenv("HOME"), ".env"), // Home directory
	}

	envLoaded := false
	for _, envPath := range envPaths {
		if err := godotenv.Load(envPath); err == nil {
			envLoaded = true
			break
		}
	}

	// Don't fail if .env not found - environment variables might be set externally
	_ = envLoaded

	dbPort, err := strconv.Atoi(getEnv("POSTALCODE_DB_PORT", "3306"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	batchSize, err := strconv.Atoi(getEnv("POSTALCODE_IMPORT_BATCH_SIZE", "1000"))
	if err != nil {
		return nil, fmt.Errorf("invalid IMPORT_BATCH_SIZE: %w", err)
	}

	timeout, err := strconv.Atoi(getEnv("POSTALCODE_IMPORT_TIMEOUT", "300"))
	if err != nil {
		return nil, fmt.Errorf("invalid IMPORT_TIMEOUT: %w", err)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("POSTALCODE_DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("POSTALCODE_DB_USER", "root"),
			Password: getEnv("POSTALCODE_DB_PASSWORD", ""),
			Name:     getEnv("POSTALCODE_DB_NAME", "postalcode"),
			Charset:  getEnv("POSTALCODE_DB_CHARSET", "utf8mb4"),
		},
		Import: ImportConfig{
			BatchSize:      batchSize,
			DefaultTimeout: timeout,
		},
	}, nil
}

// GetDSN returns MySQL DSN string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.Charset,
	)
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
