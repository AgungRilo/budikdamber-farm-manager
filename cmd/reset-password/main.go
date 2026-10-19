// Command reset-password mengganti password satu user langsung di database.
// Dipakai developer untuk memulihkan akses (misalnya password owner tidak cocok).
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load() // tidak menimpa env yang sudah diset dari Makefile

	dsn := os.Getenv("DATABASE_URL")
	email := os.Getenv("RESET_EMAIL")
	password := os.Getenv("RESET_PASSWORD")
	if dsn == "" || email == "" || len(password) < 8 || len(password) > 72 {
		log.Fatal("DATABASE_URL, RESET_EMAIL, dan RESET_PASSWORD (8-72 karakter) wajib diisi")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	defer conn.Close(ctx)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	tag, err := conn.Exec(ctx,
		`UPDATE users SET password_hash = $1 WHERE LOWER(email) = LOWER($2)`, string(hash), email)
	if err != nil {
		log.Fatalf("gagal update password: %v", err)
	}
	if tag.RowsAffected() == 1 {
		log.Printf("password untuk %s berhasil diganti", email)
		return
	}

	// Email tidak ditemukan: tampilkan daftar email yang ada (tanpa password) untuk membantu.
	log.Printf("email %s tidak ditemukan. User yang ada:", email)
	rows, err := conn.Query(ctx, `SELECT email, role::text, is_active FROM users ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var e, role string
		var active bool
		if err := rows.Scan(&e, &role, &active); err != nil {
			log.Fatal(err)
		}
		log.Printf("  - %s (%s, aktif=%v)", e, role, active)
	}
	os.Exit(1)
}
