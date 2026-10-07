package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/handler"
)

func New(db *pgxpool.Pool, appEnv string) *gin.Engine {
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	health := handler.NewHealthHandler(db)
	r.GET("/health", health.Check)

	// Route API berikutnya masuk ke sini:
	// api := r.Group("/api/v1")

	return r
}
