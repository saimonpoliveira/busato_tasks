package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	GinMode            string
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHours int
	UploadDir          string
	MaxUploadSizeMB    int64
	CORSOrigins        []string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExp, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %w", err)
	}

	maxUpload, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "10"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_UPLOAD_SIZE_MB: %w", err)
	}

	corsOrigins := resolveCORSOrigins()

	return &Config{
		Port:               getEnv("PORT", "8080"),
		GinMode:            getEnv("GIN_MODE", "release"),
		DatabaseURL:        normalizeDatabaseURL(getEnv("DATABASE_URL", "")),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpirationHours: jwtExp,
		UploadDir:          getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSizeMB:    maxUpload,
		CORSOrigins:        corsOrigins,
	}, nil
}

func resolveCORSOrigins() []string {
	if explicit := getEnv("CORS_ORIGINS", ""); explicit != "" {
		return splitAndTrim(explicit)
	}

	// Railway injects public domain for the frontend service
	if railwayDomain := getEnv("RAILWAY_PUBLIC_DOMAIN", ""); railwayDomain != "" {
		return []string{fmt.Sprintf("https://%s", railwayDomain)}
	}

	return []string{"http://localhost:5173"}
}

func normalizeDatabaseURL(url string) string {
	if url == "" {
		return ""
	}

	url = strings.Replace(url, "postgresql://", "postgres://", 1)

	if strings.Contains(url, "sslmode=") {
		return url
	}

	separator := "?"
	if strings.Contains(url, "?") {
		separator = "&"
	}

	if strings.Contains(url, "railway.internal") {
		return url + separator + "sslmode=disable"
	}

	if strings.Contains(url, "railway") || strings.Contains(url, "neon.tech") {
		return url + separator + "sslmode=require"
	}

	return url
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
