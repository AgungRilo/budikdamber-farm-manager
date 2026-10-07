package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	AppEnv      string
	DatabaseURL string
}

func Load() Config {
	// Di lokal baca .env; di Railway env sudah di-set, jadi error diabaikan.
	_ = godotenv.Load()

	cfg := Config{
		Port:        getEnv("PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "development"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL wajib diisi")
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
