package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/AgungRilo/budikdamber-farm-manager/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	status, dbStatus, code := "ok", "up", http.StatusOK
	if err := h.db.Ping(ctx); err != nil {
		status, dbStatus, code = "degraded", "down", http.StatusServiceUnavailable
	}

	c.JSON(code, response.Body{
		Success: status == "ok",
		Message: "health check",
		Data: gin.H{
			"status":   status,
			"database": dbStatus,
			"time":     time.Now().UTC(),
		},
	})
}
