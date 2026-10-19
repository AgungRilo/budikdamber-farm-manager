// Command seed-demo mengisi database dengan data contoh untuk demo/presentasi.
package main

import (
	"context"
	_ "embed"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

//go:embed demo.sql
var demoSQL string // isi demo.sql ikut ter-compile ke binary

func main() {
	_ = godotenv.Load() // tidak menimpa env yang sudah diset (penting untuk seed-demo-prod)

	dsn := os.Getenv("DATABASE_URL")
	workerPW := os.Getenv("SEED_WORKER_PASSWORD")
	if dsn == "" || len(workerPW) < 8 {
		log.Fatal("DATABASE_URL dan SEED_WORKER_PASSWORD (minimal 8 karakter) wajib diisi")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("DATABASE_URL tidak valid: %v", err)
	}
	cfg.OnNotice = func(_ *pgconn.PgConn, n *pgconn.Notice) { log.Println(n.Message) } // tampilkan RAISE NOTICE

	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	defer conn.Close(ctx)

	hash, err := bcrypt.GenerateFromPassword([]byte(workerPW), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := conn.Exec(ctx, `
		INSERT INTO users (name, email, password_hash, role, is_active)
		VALUES ('Pekerja Demo', 'pekerja@budikdamber.local', $1, 'worker', TRUE)
		ON CONFLICT (LOWER(email)) DO NOTHING`, string(hash)); err != nil {
		log.Fatalf("gagal membuat akun worker demo: %v", err)
	}

	// Tanpa argumen, pgx mengirim query lewat simple protocol,
	// sehingga blok DO $$ ... $$ berjalan utuh dalam satu transaksi.
	if _, err := conn.Exec(ctx, demoSQL); err != nil {
		log.Fatalf("gagal membuat data demo: %v", err)
	}
	log.Println("seed demo selesai")
}
