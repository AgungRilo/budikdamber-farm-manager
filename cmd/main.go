package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/config"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/database"
	"github.com/AgungRilo/budikdamber-farm-manager/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("gagal konek database: %v", err)
	}
	defer db.Close()

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router.New(db, cfg.AppEnv),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server jalan di http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown: tunggu Ctrl+C / sinyal stop dari Railway
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("mematikan server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
