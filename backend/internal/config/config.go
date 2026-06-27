package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	JWTSecret         string
	DataEncryptionKey string
	AccessTokenTTL    time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/equipment_booking?sslmode=disable"),
		JWTSecret:         getEnv("JWT_SECRET", "change-me-in-production"),
		DataEncryptionKey: getEnv("DATA_ENCRYPTION_KEY", getEnv("JWT_SECRET", "change-me-in-production")),
		AccessTokenTTL:    24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
