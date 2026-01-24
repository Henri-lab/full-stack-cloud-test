package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	Environment string
	CORSOrigin  string
}

func Load() *Config {
	env := getEnv("ENVIRONMENT", "development")
	jwtSecret := os.Getenv("JWT_SECRET")

	// In production, JWT_SECRET must be set
	if env == "production" && jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required in production")
	}

	// In development, generate a random secret if not set (but warn)
	if jwtSecret == "" {
		jwtSecret = generateRandomSecret()
		log.Println("WARNING: JWT_SECRET not set, using random secret. Sessions will not persist across restarts.")
	}

	// Validate JWT secret length
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/fullstack?sslmode=disable"),
		JWTSecret:   jwtSecret,
		Port:        getEnv("PORT", "8080"),
		Environment: env,
		CORSOrigin:  getEnv("CORS_ORIGIN", "*"),
	}
}

func generateRandomSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Failed to generate random secret")
	}
	return hex.EncodeToString(bytes)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
