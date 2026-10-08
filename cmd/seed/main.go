package main

import (
	"context"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/config"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/database"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	name := os.Getenv("SEED_OWNER_NAME")
	email := os.Getenv("SEED_OWNER_EMAIL")
	password := os.Getenv("SEED_OWNER_PASSWORD")
	if name == "" || email == "" || len(password) < 8 {
		log.Fatal("SEED_OWNER_NAME, SEED_OWNER_EMAIL wajib diisi dan SEED_OWNER_PASSWORD minimal 8 karakter")
	}

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("gagal hash password: %v", err)
	}

	tag, err := db.Exec(ctx, `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, 'owner')
		ON CONFLICT (LOWER(email)) DO NOTHING`,
		name, email, string(hash),
	)
	if err != nil {
		log.Fatalf("gagal membuat owner: %v", err)
	}

	if tag.RowsAffected() == 0 {
		log.Printf("owner %s sudah ada, dilewati", email)
	} else {
		log.Printf("owner %s berhasil dibuat", email)
	}
}
