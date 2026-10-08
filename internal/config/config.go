package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	AppEnv       string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	ttl, err := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	if err != nil {
		log.Fatalf("JWT_EXPIRES_IN tidak valid: %v", err)
	}

	cfg := Config{
		Port:         getEnv("PORT", "8080"),
		AppEnv:       getEnv("APP_ENV", "development"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		JWTExpiresIn: ttl,
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL wajib diisi")
	}
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET wajib diisi, minimal 32 karakter")
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
