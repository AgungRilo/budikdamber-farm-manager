# ============================================
# Budikdamber Farm Manager - Makefile
# ============================================

# Load environment variables from .env
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

APP_NAME   := budikdamber-farm-manager
MAIN_FILE  := cmd/main.go
BIN_DIR    := bin
BIN_FILE   := $(BIN_DIR)/$(APP_NAME)
MIGRATION_PATH := migrations

# Database URL untuk migrate CLI
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: help run build test tidy migrate-up migrate-down migrate-create migrate-version migrate-up-prod migrate-force migrate-drop docker-up docker-down docker-logs clean

## help: Menampilkan daftar perintah yang tersedia
help:
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'

## run: Menjalankan aplikasi
run:
	go run $(MAIN_FILE)

## build: Build binary aplikasi
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_FILE) $(MAIN_FILE)

## test: Menjalankan semua test
test:
	go test ./... -v

## tidy: Membersihkan dependencies
tidy:
	go mod tidy

## docker-up: Menjalankan PostgreSQL via docker-compose
docker-up:
	docker compose up -d

## docker-down: Menghentikan docker-compose
docker-down:
	docker compose down

## docker-logs: Melihat log docker
docker-logs:
	docker compose logs -f db

## migrate-up: Menjalankan migration ke atas
migrate-up:
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

## migrate-down: Rollback 1 migration
migrate-down:
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down 1

## migrate-create: Membuat file migration baru (make migrate-create name=create_users_table)
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

## migrate-version: Melihat versi migration saat ini
migrate-version:
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" version

## migrate-up-prod: Menjalankan migration ke database Railway
migrate-up-prod:
	migrate -path $(MIGRATION_PATH) -database "$(RAILWAY_DATABASE_URL)" up

## migrate-force: Force versi migration (make migrate-force v=1)
migrate-force:
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" force $(v)

## migrate-drop: HAPUS SEMUA tabel di database (hati-hati!)
migrate-drop:
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" drop

## clean: Menghapus hasil build
clean:
	rm -rf $(BIN_DIR)
